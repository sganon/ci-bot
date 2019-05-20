package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type PiepelineID int

type Pipeline struct {
	ID     PiepelineID `json:"id"`
	Status string      `json:"status"`
	Ref    string      `json:"ref"`
	SHA    string      `json:"sha"`
	WebURL string      `json:"web_url"`
}

func (api *API) GetRefLastPipeline(proj Project, ref string) (pp Pipeline, err error) {
	// TODO: make an api method for http call
	// will reduce code duplication and could make this method part of its model
	endpoint, err := url.Parse(api.baseURL +
		"/projects/" + strconv.Itoa(int(proj.ID)) + "/pipelines")
	if err != nil {
		return pp, fmt.Errorf("error parsing GET /projects/%d/pipelines/ url: %v", proj.ID, err)
	}
	// we search for name in every case
	query := url.Values{}
	query.Add("ref", ref)
	endpoint.RawQuery = query.Encode()

	log.Debugf("Requesting: %s", endpoint.String())
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return pp, fmt.Errorf("error creating GET /projects/%d/pipelines/ request: %v", proj.ID, err)
	}
	req.Header.Add("Private-Token", api.token)
	res, err := api.client.Do(req)
	if err != nil {
		return pp, fmt.Errorf("error sending GET /projects/%d/pipelines/ request: %v", proj.ID, err)
	}
	if res.StatusCode != http.StatusOK {
		return pp, fmt.Errorf("unexpected status code: %d, expected 200", res.StatusCode)
	}
	defer res.Body.Close()
	var pipelines []Pipeline
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pipelines); err != nil {
		return pp, fmt.Errorf("error decoding response body")
	}
	return pipelines[0], nil
}
