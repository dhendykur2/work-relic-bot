package models

type Task struct {
	Task        string `json:"task"`
	DueDate     string `json:"due_date"`
	Description string `json:"description"`
	DoneAt      string `json:"done_at,omitempty"`
}
