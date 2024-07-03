package models

import(
	"time"
)

type User struct {
    ID              uint      `json:"id" gorm:"primaryKey"`
    Name            string    `json:"name" validate:"required"`
    Email           string    `json:"email" validate:"required,email"`
    RegistrationDate time.Time `json:"registration_date" gorm:"default:now()"`
    Role            string    `json:"role" validate:"required,oneof=admin manager developer"`
}


