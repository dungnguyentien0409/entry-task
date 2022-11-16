package dto

import "time"

type User struct {
	Id       int64
	Account  string
	Nickname string
	Password string
	Salt     string
}

type Product struct {
	Id          int
	CategoryId  int
	Name        string
	Description string
	Price       float32
	Currency    string
	Images      string
}

type Comment struct {
	Id              int
	ProductId       int
	UserId          int
	Account         string
	Content         string
	ParentCommentId int
	InsertedAt      time.Time
}

type DBConfig struct {
	ConnectionString string
	DriverName       string
	MaxIdleConns     int
	MaxOpenConns     int
	ConnMaxLifeTime  time.Duration
}
