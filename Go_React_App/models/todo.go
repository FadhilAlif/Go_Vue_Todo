package models

type Todo struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}
