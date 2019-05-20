package gitlab

type PiepelineID int

type Pipeline struct {
	ID     PiepelineID `json:"id"`
	Status string      `json:"status"`
	Ref    string      `json:"ref"`
	SHA    string      `json:"sha"`
	WebURL string      `json:"web_url"`
}
