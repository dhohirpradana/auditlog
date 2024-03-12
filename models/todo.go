package models

// Todo represent the todo structs
type Todo struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId" validate:"nonzero"`
	Title     string `json:"title" validate:"nonnil,nonzero"`
	Completed bool   `json:"completed" validate:"nonnil,nonzero"`
}
