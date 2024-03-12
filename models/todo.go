package models

// Todo represent the todo structs
type Todo struct {
	UserID    int    `json:"userId" validate:"nonzero"`
	ID        int    `json:"id"`
	Title     string `json:"title" validate:"nonnil"`
	Completed bool   `json:"completed" validate:"nonnil"`
}
