package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sganon/code-bot/slack"
)

func jsonResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{ "text": "an unexpected error occured encoding response" }`)
		return
	}
	w.WriteHeader(status)
	fmt.Fprintln(w, string(b))
}

func slackErrorResponse(w http.ResponseWriter, msg string) {
	attch := slack.ErrorMessage(msg)
	jsonResponse(w, http.StatusOK, slackResponse{
		ResponseType: inChannelType,
		Attachments:  []slack.Attachment{attch},
	})
}
