package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io/ioutil"
	"net/http"

	"fmt"

	"strings"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"encoding/hex"
	"encoding/json"

	"github.com/sganon/ci-bot/gitlab"
	"github.com/sganon/ci-bot/slack"
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
		fmt.Fprintln(w, `{ "text": "usage: [group]/project tag" }`)
		return
	}
	name := parts[0]
	tag := parts[1]

	pj, err := gitlab.GetProjectByName(api.glAPI, name)
	if err != nil {
		fmt.Fprintln(w, `{ "text": "unable to find matching project" }`)
		return
	}

	pj.Tag.Name = tag
	err = pj.FetchTagPipelines(api.glAPI)
	if err != nil {
		fmt.Fprintln(w, `{ "text": "error fetching pipelines" }`)
		return
	}

	err = pj.FetchTag(api.glAPI)
	if err != nil {
		fmt.Fprintln(w, `{ "text": "error fetching tag" }`)
		return
	}

	b, err := json.Marshal(struct {
		ResponseType string             `json:"response_type"`
		Attachements []slack.Attachment `json:"attachments"`
	}{
		ResponseType: "in_channel",
		Attachements: []slack.Attachment{pj.Attachement()},
	})
	if err != nil {
		fmt.Fprintln(w, `{ "text": "error encoding response" }`)
		return
	}

	fmt.Fprintf(w, string(b))
}
