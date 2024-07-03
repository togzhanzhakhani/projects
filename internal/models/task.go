package models

import "time"

type Task struct {
	ID           int       `json:"id"`
	Title        string    `json:"title" validate:"required"`
	Description  string    `json:"description" validate:"required,max=100"`
	Priority     string    `json:"priority" validate:"oneof=low medium high"`
	Status       string    `json:"status" validate:"oneof=todo in_progress done"`
	AssigneeID   int       `json:"assignee_id" validate:"required,gt=0"`
	ProjectID    int       `json:"project_id" validate:"required,gt=0"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	CompletedAt  time.Time `json:"completed_at" validate:"required,gtfield=CreatedAt"`
}

