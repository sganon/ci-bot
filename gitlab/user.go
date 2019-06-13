package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	State     string `json:"active"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func GetUserByName(api *API, username string) (user User, err error) {
	query := url.Values{}
	query.Add("username", username)

	var users []User
	statusCode, err := api.Call("GET", "/users", query, nil, &users)
	if err != nil {
		return user, fmt.Errorf("GetUserByName error: %v", err)
	}
	if statusCode != http.StatusOK {
		return user, fmt.Errorf("unexpected status code: %d, expected 200", statusCode)
	}
	return users[0], err
}

// GetAssignedMRs returns all opened MR assigned to the user
func (u User) GetAssignedMRs(api *API) (mrs []MR, err error) {
	query := url.Values{}
	query.Add("state", "opened")
	query.Add("scope", "all")
	query.Add("assignee_id", strconv.Itoa(u.ID))
	query.Add("sort", "asc")

	statusCode, err := api.Call("GET", "/merge_requests", query, nil, &mrs)
	if err != nil {
		return mrs, fmt.Errorf("GetAssignedMRs error: %v", err)
	}
	if statusCode != http.StatusOK {
		return mrs, fmt.Errorf("unexpected status code: %d, expected 200", statusCode)
	}

	return mrs, err
}
