package gitlab

// Tag is a gitlab tag
type Tag struct {
	Name    string `json:"name"`
	Message string `json:"Message"`
	Target  string `json:"target"`
	Release struct {
		TagName     string `json:"tag_name"`
		Description string `json:"description"`
	} `json:"release"`

	Pipelines []Pipeline
}
