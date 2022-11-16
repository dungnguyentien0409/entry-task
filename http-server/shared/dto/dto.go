package dto

import "time"

type User struct {
	Id       int64  `json:"id"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type Product struct {
	Id          int       `json:"id"`
	CategoryId  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Currency    string    `json:"currency"`
	Images      string    `json:"images"`
	Comments    []Comment `json:"comments"`
}

type Comment struct {
	Id              int       `json:"id"`
	ProductId       int       `json:"product_id"`
	UserId          int       `json:"user_id"`
	Account         string    `json:"account"`
	Content         string    `json:"content"`
	ParentCommentId int       `json:"parent_comment_id"`
	InsertedAt      time.Time `json:"inserted_at"`
}
