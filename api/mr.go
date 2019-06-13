package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sganon/code-bot/gitlab"
	"github.com/sganon/code-bot/slack"
)

func (api API) handleMR(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{ "error": "unable to parse form" }`)
		return
	}

	username := r.Form.Get("text")
	if username == "" {
		slackErrorResponse(w, "Usage error: you need to provide an username")
		return
	}

	user, err := gitlab.GetUserByName(api.glAPI, username)
	if err != nil {
		slackErrorResponse(w, "Cannot find user: "+username)
		return
	}

	mrs, err := user.GetAssignedMRs(api.glAPI)
	if err != nil {
		slackErrorResponse(w, "An error occured fetching your MRs")
		return
	}
	var attchs []slack.Attachment
	timeFormat := "Mon _2 Jan 2006"
	for _, mr := range mrs {
		proj, err := gitlab.GetProjectByID(api.glAPI, mr.ProjectID)
		if err != nil {
			slackErrorResponse(w, "An unexpected error occured fetching project details")
			return
		}
		cAt, _ := time.Parse(time.RFC3339, mr.CreatedAt)
		uAt, _ := time.Parse(time.RFC3339, mr.UpdatedAt)
		attch := slack.Attachment{
			Title:      mr.Title,
			AuthorName: mr.Author.Name,
			AuthorIcon: mr.Author.AvatarURL,
			Color:      "#008bd2",
			Fields: []slack.Field{
				slack.Field{Title: "Project", Value: proj.NameWithNamespace},
				slack.Field{Title: "Source Branch", Value: mr.SourceBranch, Short: true},
				slack.Field{Title: "Target Branch", Value: mr.TargetBranch, Short: true},
				slack.Field{Title: "Opened At", Value: cAt.Format(timeFormat), Short: true},
				slack.Field{Title: "Updated At", Value: uAt.Format(timeFormat), Short: true},
			},
		}
		attchs = append(attchs, attch)
	}
	jsonResponse(w, http.StatusOK, slackResponse{
		ResponseType: inChannelType,
		Attachments:  attchs,
	})
}
