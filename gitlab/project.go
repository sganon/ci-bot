package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sganon/code-bot/slack"
)

type ProjectID int

type Project struct {
	ID                ProjectID `json:"id"`
	Description       string    `json:"description"`
	DefaultBranch     string    `json:"default_branch"`
	SSHURL            string    `json:"ssh_url_to_repo"`
	HTTPURL           string    `json:"http_url_to_repo"`
	WebURL            string    `json:"web_url"`
	ReadmeURL         string    `json:"readme_url"`
	TagList           []string  `json:"tag_list"`
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	CreatedAt         time.Time `json:"created_at"`
	LastActivity      time.Time `json:"last_activity_at"`
	ForksCount        int       `json:"forks_count"`
	AvatarURL         string    `json:"avatar_url"`
	StarCount         int       `json:"star_count"`

	Tag Tag
}

// GetProjectByName will return specific project from GET /projects via its name
// name could be of form group/project_name if forks exist
func GetProjectByName(api *API, name string) (proj Project, err error) {
	splitedName := strings.Split(name, "/")
	if len(splitedName) > 2 {
		return proj, fmt.Errorf("project name should be of form [GROUP]/PROJECT")
	}

	// we search for name in every case
	query := url.Values{}
	if len(splitedName) == 2 {
		query.Add("search", splitedName[1])
	} else {
		query.Add("search", name)
	}

	var projects []Project
	statusCode, err := api.Call("GET", "/projects", query, nil, &projects)
	if err != nil {
		return proj, fmt.Errorf("GetProjectByName error: %v", err)
	}
	if statusCode != http.StatusOK {
		return proj, fmt.Errorf("unexpected status code: %d, expected 200", statusCode)
	}

	switch len(projects) {
	// This should happen when no forks and name is exact
	case 0:
		return proj, fmt.Errorf("empty list of projects")
	case 1:
		return projects[1], nil
	default:
		for _, p := range projects {
			if p.PathWithNamespace == name {
				return p, nil
			}
		}
	}

	return proj, fmt.Errorf("unable to find unique matching project")
}

func (pj *Project) FetchTagPipelines(api *API) error {
	if pj.Tag.Name == "" {
		return fmt.Errorf("FetchRefPipelines error: ref must be set before fetching its pipelines")
	}

	statusCode, err := api.Call("GET",
		fmt.Sprintf("/projects/%d/pipelines", int(pj.ID)),
		url.Values{"ref": []string{pj.Tag.Name}}, nil, &pj.Tag.Pipelines)
	if err != nil {
		return fmt.Errorf("FetchTagPipelines error: %v", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("FetchTagPipelines: unexpected status code %d, expected 200", statusCode)
	}
	return nil
}

func (pj *Project) FetchTag(api *API) error {
	if pj.Tag.Name == "" {
		return fmt.Errorf("FetchTag error: ref must be set before fetching its pipelines")
	}
	statusCode, err := api.Call("GET",
		fmt.Sprintf("/projects/%d/repository/tags/%s", int(pj.ID), pj.Tag.Name),
		nil, nil, &pj.Tag)
	if err != nil {
		return fmt.Errorf("FetchTag error: %v", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("FetchTag: unexpected status code %d, expected 200", statusCode)
	}
	return nil
}

func (pj Project) Attachment() slack.Attachment {
	main := fmt.Sprintf("New release of %s: *<%s/tags/%s|%s>*",
		pj.PathWithNamespace, pj.WebURL, pj.Tag.Name, pj.Tag.Name)

	fields := []slack.Field{
		slack.Field{
			Title: "Changelog",
			Value: strings.ReplaceAll(pj.Tag.Release.Description, "*", "•"),
		},
	}
	if len(pj.Tag.Pipelines) > 0 {
		lastPipeline := pj.Tag.Pipelines[len(pj.Tag.Pipelines)-1]
		fields = append(fields, slack.Field{
			Title: "Last pipeline",
			Value: fmt.Sprintf("<%s|#%d>: %s", lastPipeline.WebURL, lastPipeline.ID, lastPipeline.Status),
		})
	}

	return slack.Attachment{
		Fallback: main,
		Color:    "#008bd2",
		Pretext:  main,
		Fields:   fields,
	}
}
