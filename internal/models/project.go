package models

import "time"

type Project struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description" validate:"required,max=100"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	ManagerID    int       `json:"manager_id" validate:"required,gt=0"`
}

