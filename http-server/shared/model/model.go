package model

import (
	"entrytask/http-server/shared/dto"
	"time"
)

type Request struct {
	Cmd  int         `json:"cmd_type"`
	Data interface{} `json:"data"`
}

type Response struct {
	Status  int         `json:"error_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	PassWord string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	PassWord string `json:"password"`
}

type GetListProductsRequest struct {
	CategoryId int    `json:"category_id"`
	Name       string `json:"name"`
	PageIndex  int    `json:"page_index"`
	PageSize   int    `json:"page_size"`
}

type GetProductRequest struct {
	ProductId int
}

type PostCommentRequest struct {
	ProductId       int
	UserId          int
	Account         string
	ParentCommentId int
	Content         string
	InsertedAt      time.Time
}

type GetCommentRequest struct {
	ProductId int
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}

type GetListProductsResponse struct {
	Products []dto.Product `json:"products"`
	Total    int           `json:"total"`
}

type GetListCommentResponse struct {
	Comments []dto.Comment `json:"products"`
	Total    int           `json:"total"`
}
