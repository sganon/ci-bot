package gitlab

type MR struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"project_id"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
	Author       User   `json:"author"`
}
