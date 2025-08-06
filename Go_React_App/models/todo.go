package models

import (
	"github.com/google/uuid"
)

type Todo struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Completed bool      `json:"completed"`
	Body      string    `json:"body"`
}
