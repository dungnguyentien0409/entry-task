package product_service

import (
	dao "entrytask/http-server/db"
	"entrytask/http-server/shared/constants"
	"entrytask/http-server/shared/model"
	_ "github.com/go-sql-driver/mysql"
)

type ProductService struct {
	productRepository dao.IProductRepository
}

func NewProductService(productRepository dao.IProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (p *ProductService) GetListProducts(getListProductsRequest interface{}) model.Response {
	request := getListProductsRequest.(model.GetListProductsRequest)
	response := model.Response{}

	if request.PageSize > constants.DEFAULT_PAGE_SIZE {
		request.PageSize = constants.DEFAULT_PAGE_SIZE
	}

	data, total, err := p.productRepository.GetListProducts(request)
	if err != nil {
		response.Status = constants.GET_LIST_PRODUCTS_FAILED
		response.Message = "get list products failed"
		return response
	}

	response.Status = constants.GET_LIST_PRODUCTS_SUCCESSED
	response.Message = "get list products success"
	response.Data = model.GetListProductsResponse{Products: data, Total: total}
	return response
}

func (p *ProductService) GetProduct(getProductRequest interface{}) model.Response {
	request := getProductRequest.(model.GetProductRequest)
	response := model.Response{}

	res, err := p.productRepository.GetProduct(request)
	if err != nil {
		response.Status = constants.GET_PRODUCT_FAILED
		response.Message = "get product failed"
		return response
	}

	response.Status = constants.GET_PRODUCT_SUCCESSED
	response.Message = "get product successed"
	response.Data = res
	return response
}

func (p *ProductService) PostComment(postComment interface{}) model.Response {
	request := postComment.(model.PostCommentRequest)
	response := model.Response{}

	if request.Content == "" {
		response.Status = constants.COMMENT_EMPTY
		response.Message = "Comment cannot be empty"
		return response
	}

	err := p.productRepository.PostComment(request)
	if err != nil {
		response.Status = constants.POST_COMMENT_FAILED
		response.Message = "post comment failed"
		return response
	}

	response.Status = constants.POST_COMMENT_SUCCESSED
	response.Message = "post comment success"
	response.Data = nil
	return response
}

func (p *ProductService) GetComment(getCommentRequest interface{}) model.Response {
	request := getCommentRequest.(model.GetCommentRequest)
	response := model.Response{}

	if request.PageSize > constants.DEFAULT_PAGE_SIZE {
		request.PageSize = constants.DEFAULT_PAGE_SIZE
	}

	data, total, err := p.productRepository.GetComment(request, false)
	if err != nil {
		response.Status = constants.GET_PRODUCT_FAILED
		response.Message = "get comment failed"
		return response
	}

	response.Status = constants.GET_COMMENT_SUCCESSED
	response.Message = "get comment success"
	response.Data = model.GetListCommentResponse{Comments: data, Total: total}
	return response
}
