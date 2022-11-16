package dao

import (
	"entrytask/http-server/shared/dto"
	"entrytask/http-server/shared/model"
)

type IProductRepository interface {
	GetListProducts(request model.GetListProductsRequest) (products []dto.Product, total int, err error)
	GetProduct(request model.GetProductRequest) (product dto.Product, err error)
	PostComment(request model.PostCommentRequest) error
	GetComment(request model.GetCommentRequest, getAll bool) (comments []dto.Comment, total int, err error)
}
