package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/sganon/code-bot/gitlab"
	"github.com/sganon/code-bot/slack"
)

// API expose routes to handle slacke request
type API struct {
	host         string
	port         string
	signinSecret string
	glAPI        *gitlab.API

	handler http.Handler
}

func New(host, port, signinSecret string, glAPI *gitlab.API) *API {
	api := API{
		host:         host,
		port:         port,
		signinSecret: signinSecret,
		glAPI:        glAPI,
	}
	api.routes()
	api.handler = signinHandler{api.handler, api.signinSecret}
	return &api
}

type signinHandler struct {
	h            http.Handler
	signinSecret string
}

func (h signinHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	baseString := "v0:" + r.Header.Get("X-Slack-Request-Timestamp") + ":" + string(bodyBytes)
	hash := hmac.New(sha256.New, []byte(h.signinSecret))
	hash.Write([]byte(baseString))
	if "v0="+hex.EncodeToString(hash.Sum(nil)) != r.Header.Get("X-Slack-Signature") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	h.h.ServeHTTP(w, r)
}

func (api *API) routes() {
	router := httprouter.New()

	router.POST("/ack/", api.handleACK)
	router.POST("/release/", api.handleRelease)
	router.POST("/mr/", api.handleMR)

	api.handler = router
}

func (api API) Serve() {
	addr := fmt.Sprintf("%s:%s", api.host, api.port)
	log.Infof("serving api on: %s", addr)
	http.ListenAndServe(addr, handlers.RecoveryHandler()(api.handler))
}

func (api API) handleACK(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "OK")
}

type slackResponse struct {
	ResponseType string             `json:"response_type"`
	Attachments  []slack.Attachment `json:"attachments"`
}

const inChannelType = "in_channel"

func (api API) handleRelease(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{ "error": "unable to parse form" }`)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	text := r.Form.Get("text")
	parts := strings.Split(text, " ")
	if len(parts) != 2 {
		slackErrorResponse(w, "*usage*: [group]/project tag")
		return
	}
	name := parts[0]
	tag := parts[1]

	pj, err := gitlab.GetProjectByName(api.glAPI, name)
	if err != nil {
		slackErrorResponse(w, "Cannot find project")
		return
	}

	pj.Tag.Name = tag
	err = pj.FetchTagPipelines(api.glAPI)
	if err != nil {
		slackErrorResponse(w, "Error fetching project pipelines")
		return
	}

	err = pj.FetchTag(api.glAPI)
	if err != nil {
		slackErrorResponse(w, "Error fetching project tagss")
		return
	}

	jsonResponse(w, http.StatusOK, slackResponse{
		ResponseType: inChannelType,
		Attachments:  []slack.Attachment{pj.Attachment()},
	})
}
