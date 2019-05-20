package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// API is the Gitlab's API wrapper
type API struct {
	addr    string
	token   string
	baseURL string
	client  *http.Client
}

func NewAPI(addr string, token string) *API {
	client := http.Client{
		Timeout: time.Second * 60,
	}
	return &API{
		addr:    addr,
		baseURL: addr + "/api/v4",
		token:   token,
		client:  &client,
	}
}

// Call handles http call and json unmarshaling
// it will return response code and error if any
// if error occurs before response status code will be 0
func (api *API) Call(method string, path string, query url.Values, body io.Reader, dest interface{}) (int, error) {
	logFields := log.Fields{
		"method": method,
		"path":   path,
		"query":  query.Encode(),
	}
	endpoint, err := url.Parse(api.baseURL + path)
	if err != nil {
		log.WithFields(logFields).Error(err)
		return 0, fmt.Errorf("api call error: creating url: %v", err)
	}
	endpoint.RawQuery = query.Encode()
	logFields["url"] = endpoint.String()
	log.WithFields(logFields).Debugln("Call Gitlab API")

	req, err := http.NewRequest(method, endpoint.String(), body)
	if err != nil {
		log.WithFields(logFields).Error(err)
		return 0, fmt.Errorf("api call error: creating request: %v", err)
	}
	req.Header.Add("Private-Token", api.token)
	req.Header.Add("Accept", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		log.WithFields(logFields).Error(err)
		return res.StatusCode, fmt.Errorf("api call error: sending request: %v", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(dest); err != nil {
		return res.StatusCode, fmt.Errorf("api call error: decoding response: %v", err)
	}

	return res.StatusCode, nil
}
