package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

	// Struct representig a ref and its
	// associated pipelines
	Ref struct {
		Value     string
		Pipelines []Pipeline
	}
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

func (pj *Project) FetchRefPipelines(api *API) error {
	if pj.Ref.Value == "" {
		return fmt.Errorf("FetchRefPipelines error: ref must be set before fetching its pipelines")
	}

	statusCode, err := api.Call("GET",
		fmt.Sprintf("/projects/%d/pipelines", int(pj.ID)),
		url.Values{"ref": []string{pj.Ref.Value}}, nil, &pj.Ref.Pipelines)
	if err != nil {
		return fmt.Errorf("FetchRefPipelines error: %v", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("FetchRefPipelines: unexpected status code %d, expected 200", statusCode)
	}
	return nil
}
