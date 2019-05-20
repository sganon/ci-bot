package gitlab

// PipelineID is the unique id of a CI pipeline
type PipelineID int

// Pipeline is the representation of a CI pipeline
type Pipeline struct {
	ID     PipelineID `json:"id"`
	Status string     `json:"status"`
	Ref    string     `json:"ref"`
	SHA    string     `json:"sha"`
	WebURL string     `json:"web_url"`
}
