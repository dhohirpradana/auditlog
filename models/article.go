package models

// Article represent the article structs
type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
