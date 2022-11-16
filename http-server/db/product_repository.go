package dao

import (
	"context"
	"database/sql"
	"entrytask/http-server/shared/constants"
	"entrytask/http-server/shared/dto"
	"entrytask/http-server/shared/model"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(DB *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: DB,
	}
}

func (u *ProductRepository) GetListProducts(request model.GetListProductsRequest) (products []dto.Product, total int, err error) {
	_, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*300))
	defer cancel()

	queryStatement := GetQueryStatement(request)
	fmt.Println(queryStatement)
	rows, _ := u.db.Query(queryStatement)

	product := dto.Product{}
	for rows.Next() {
		rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Description, &product.Price, &product.Currency, &product.Images)
		product.Comments, _, _ = u.GetComment(model.GetCommentRequest{
			ProductId: product.Id,
			PageIndex: constants.DEFAULT_PAGE_INDEX,
			PageSize:  constants.DEFAULT_PAGE_SIZE,
		}, true)
		products = append(products, product)
	}

	total = u.GetTotalRowsProduct(request)

	return products, total, nil
}

func (u *ProductRepository) GetProduct(request model.GetProductRequest) (product dto.Product, err error) {
	_, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*300))
	defer cancel()

	queryStatement := fmt.Sprintf("SELECT id, category_id, name, description, price, currency, images "+
		"FROM product_tab  "+
		"WHERE id = %d", request.ProductId)
	rows, _ := u.db.Query(queryStatement)

	for rows.Next() {
		rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Description, &product.Price, &product.Currency, &product.Images)
	}

	return product, nil
}

func (u *ProductRepository) PostComment(request model.PostCommentRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*300))
	defer cancel()
	_, err := u.db.ExecContext(ctx,
		"INSERT INTO comment_tab (product_id, user_id, account, content, parent_comment_id, inserted_at) values(?, ? ,?, ?, ?, ?)",
		request.ProductId, request.UserId, request.Account, request.Content, request.ParentCommentId, time.Now().UTC())
	if err != nil {
		return errors.New("Failed")
	}

	return nil
}

func (u *ProductRepository) GetComment(request model.GetCommentRequest, getAll bool) (comments []dto.Comment, total int, err error) {
	_, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*300))
	defer cancel()

	queryStatement := fmt.Sprintf("SELECT id, product_id, user_id, account, content, parent_comment_id, inserted_at "+
		"FROM comment_tab  "+
		"WHERE product_id = %d "+
		"ORDER BY inserted_at ", request.ProductId)

	if !getAll {
		queryStatement += fmt.Sprintf("LIMIT %d, %d", request.PageIndex*request.PageSize, request.PageSize)
	}

	rows, _ := u.db.Query(queryStatement)

	comment := dto.Comment{}
	for rows.Next() {
		rows.Scan(&comment.Id, &comment.ProductId, &comment.UserId, &comment.Account, &comment.Content, &comment.ParentCommentId, &comment.InsertedAt)
		comments = append(comments, comment)
	}

	total = u.GetTotalRowsComment(request)

	return comments, total, nil
}

func (u *ProductRepository) GetTotalRowsComment(request model.GetCommentRequest) int {
	queryStatement := fmt.Sprintf("SELECT COUNT(*) "+
		"FROM comment_tab  "+
		"WHERE product_id = %d ", request.ProductId)

	rows, err := u.db.Query(queryStatement)
	if err != nil {
		return 0
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	return count
}

func (u *ProductRepository) GetTotalRowsProduct(request model.GetListProductsRequest) int {
	queryString := "SELECT COUNT(*) " +
		"FROM product_tab  " +
		`WHERE name like '%` + request.Name + `%' `

	if request.CategoryId > 0 {
		queryString += fmt.Sprintf("and category_id = %d ", request.CategoryId)
	}

	rows, err := u.db.Query(queryString)
	if err != nil {
		return 0
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	return count
}

func GetQueryStatement(request model.GetListProductsRequest) string {
	skip, take := request.PageIndex*request.PageSize, request.PageSize

	queryString := "SELECT id, category_id, name, description, price, currency, images " +
		"FROM product_tab  " +
		`WHERE name like '%` + request.Name + `%' `

	if request.CategoryId > 0 {
		queryString += fmt.Sprintf("and category_id = %d ", request.CategoryId)
	}

	queryString += fmt.Sprintf("ORDER BY id ")
	queryString += fmt.Sprintf("LIMIT %d, %d", skip, take)
	return queryString
}
