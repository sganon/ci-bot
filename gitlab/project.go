package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
}

// GetProjectByName will return specific project from GET /projects via its name
// name could be of form group/project_name if forks exist
func (api *API) GetProjectByName(name string) (proj Project, err error) {
	splitedName := strings.Split(name, "/")
	if len(splitedName) > 2 {
		return proj, fmt.Errorf("project name should be of form [GROUP]/PROJECT")
	}

	endpoint, err := url.Parse(api.baseURL + "/projects")
	if err != nil {
		return proj, fmt.Errorf("error parsing GET /projects url: %v", err)
	}
	// we search for name in every case
	query := url.Values{}
	if len(splitedName) == 2 {
		query.Add("search", splitedName[1])
	} else {
		query.Add("search", name)
	}
	endpoint.RawQuery = query.Encode()

	log.Debugf("Requesting: %s", endpoint.String())
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return proj, fmt.Errorf("error creating GET /projects request: %v", err)
	}
	req.Header.Add("Private-Token", api.token)
	res, err := api.client.Do(req)
	if err != nil {
		return proj, fmt.Errorf("error sending GET /projects request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return proj, fmt.Errorf("unexpected status code: %d, expected 200", res.StatusCode)
	}
	defer res.Body.Close()
	var projects []Project
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&projects); err != nil {
		return proj, fmt.Errorf("error decoding response: %v", err)
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
