package handler

import (
	"crypto/rsa"
	"encoding/json"
	"entrytask/http-server/product_service"
	"entrytask/http-server/shared/constants"
	"entrytask/http-server/shared/jwt"
	"entrytask/http-server/shared/model"
	"entrytask/http-server/user_service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	clientService  user_service.IUserService
	productService product_service.IProductService
	publicKey      *rsa.PublicKey
}

func NewHandler(service user_service.IUserService, productService product_service.IProductService, publicKey *rsa.PublicKey) Handler {
	return Handler{
		clientService:  service,
		productService: productService,
		publicKey:      publicKey,
	}
}

func (handler *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get("Authorization")
		jwtToken = strings.Split(jwtToken, "Bearer ")[1]

		_, err := jwtHelper.ValidateToken(jwtToken, handler.publicKey)
		if err != nil {
			sendAuthErrorResp(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (handler *Handler) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Set CORS headers for the main request.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	}
}

func (handler *Handler) HandleLogin(writer http.ResponseWriter, request *http.Request) {
	var loginRequest model.LoginRequest
	json.NewDecoder(request.Body).Decode(&loginRequest)
	response := handler.clientService.Login(loginRequest)

	sendResponse(writer, response)
}

func (handler *Handler) HandleRegister(writer http.ResponseWriter, request *http.Request) {
	var registerRequest model.RegisterRequest
	json.NewDecoder(request.Body).Decode(&registerRequest)
	response := handler.clientService.Register(registerRequest)

	sendResponse(writer, response)
}

func (handler *Handler) HandleGetListProducts(writer http.ResponseWriter, request *http.Request) {
	getListProductsRequest, err := DecodeGetListProductRequest(request)

	if err != nil {
		log.Println(err.Error())
	}

	response := handler.productService.GetListProducts(getListProductsRequest)

	sendResponse(writer, response)
}

func (handler *Handler) HandlePing(writer http.ResponseWriter, request *http.Request) {
	var response = model.Response{}
	response.Data = 1

	sendResponse(writer, response)
}

func (handler *Handler) HandleGetProduct(writer http.ResponseWriter, request *http.Request) {
	var getProductRequest model.GetProductRequest
	var err error
	var productId = strings.TrimPrefix(request.URL.Path, "/products/")

	getProductRequest.ProductId, err = strconv.Atoi(productId)

	if err != nil {
		log.Println(err)
	}

	response := handler.productService.GetProduct(getProductRequest)

	sendResponse(writer, response)
}

func (handler *Handler) HandlePostComment(writer http.ResponseWriter, request *http.Request) {
	var postCommentRequest model.PostCommentRequest
	json.NewDecoder(request.Body).Decode(&postCommentRequest)
	response := handler.productService.PostComment(postCommentRequest)

	sendResponse(writer, response)
}

func (handler *Handler) HandleGetComment(writer http.ResponseWriter, request *http.Request) {
	var productId = strings.TrimPrefix(request.URL.Path, "/comments/")

	getCommentRequest, err := DecodeGetCommentRequest(request)
	if err != nil {
		log.Println(err)
	}

	getCommentRequest.ProductId, err = strconv.Atoi(productId)
	if err != nil {
		log.Println(err)
	}

	response := handler.productService.GetComment(getCommentRequest)

	sendResponse(writer, response)
}

func sendResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func sendAuthErrorResp(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Not Authorized"))
}

func DecodeGetListProductRequest(request *http.Request) (getListProductsRequest model.GetListProductsRequest, err error) {
	values := request.URL.Query()
	getListProductsRequest.Name = values.Get("name")
	getListProductsRequest.CategoryId, err = strconv.Atoi(values.Get("category_id"))
	getListProductsRequest.PageIndex, err = strconv.Atoi(values.Get("page_index"))

	if err != nil {
		getListProductsRequest.PageIndex = constants.DEFAULT_PAGE_INDEX
	}

	getListProductsRequest.PageSize, err = strconv.Atoi(values.Get("page_size"))

	if err != nil || getListProductsRequest.PageSize > constants.DEFAULT_PAGE_SIZE {
		getListProductsRequest.PageSize = constants.DEFAULT_PAGE_SIZE
	}

	return
}

func DecodeGetCommentRequest(request *http.Request) (getListCommentRequest model.GetCommentRequest, err error) {
	values := request.URL.Query()
	getListCommentRequest.PageIndex, err = strconv.Atoi(values.Get("page_index"))

	if err != nil {
		getListCommentRequest.PageIndex = constants.DEFAULT_PAGE_INDEX
	}

	getListCommentRequest.PageSize, err = strconv.Atoi(values.Get("page_size"))

	if err != nil || getListCommentRequest.PageIndex > constants.DEFAULT_PAGE_SIZE {
		getListCommentRequest.PageSize = constants.DEFAULT_PAGE_SIZE
	}

	return
}
