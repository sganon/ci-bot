package gitlab

import (
	"net/http"
	"time"
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
