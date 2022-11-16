package product_service

import (
	"entrytask/http-server/shared/model"
)

type IProductService interface {
	GetListProducts(getListProductsRequest interface{}) model.Response
	GetProduct(getProductRequest interface{}) model.Response
	PostComment(postComment interface{}) model.Response
	GetComment(getCommentRequest interface{}) model.Response
}
