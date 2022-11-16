package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)
import (
	_ "github.com/go-sql-driver/mysql"
)

const preUrl = "https://shopee.sg/api/v4/search/search_items"

var MaxProduct = 5000

func main() {
	InsertProduct()
}

func InsertProduct() {
	database, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/entry_task")

	if err != nil {
		panic(err.Error())
	}

	pageSize, pageIndex := 10, 0

	for pageSize*pageIndex < MaxProduct {
		data := GetCrawlData()

		sqlStr := `INSERT INTO product_tab(category_id,name,description,price,currency,images) VALUES`
		vals := []interface{}{}

		for _, item := range data.Items {
			itemBasic := item.ItemBasic

			for i := 0; i < 1000; i++ {
				sqlStr += "(?,?,?,?,?,?),"
				itemName := itemBasic.Name + strconv.Itoa(i)
				categoryId := GetCategoryId()
				description := GetDescription(categoryId, itemName, itemBasic.Price, itemBasic.Currency)
				vals = append(vals, categoryId, itemName, description, itemBasic.Price,
					itemBasic.Currency, strings.Join(itemBasic.Images, ","))
			}
		}

		sqlStr = sqlStr[0 : len(sqlStr)-1]
		stmt, _ := database.Prepare(sqlStr)

		_, err := stmt.Exec(vals...)
		if err != nil {
			panic(err.Error())
		}

		pageIndex++
	}

	database.Close()
}

func GetDescription(categoryId int, name string, price int64, currency string) string {
	return strconv.Itoa(categoryId) + " " + name + " " + strconv.FormatInt(price, 10) + " " + currency
}

func GetCategoryId() int {
	min, max := 1, 3
	return rand.Intn(max-min+1) + min
}

func GetCrawlData() CrawlData {
	var body = ReadBody("Samsung", 0, MaxProduct)
	var list CrawlData

	json.Unmarshal([]byte(body), &list)

	return list
}

func ReadBody(keyword string, pageId int, pageSize int) string {
	params := map[string]string{
		"by":        "relevancy",
		"keyword":   keyword,
		"limit":     strconv.Itoa(pageSize),          // take
		"newest":    strconv.Itoa(pageId * pageSize), //skip
		"order":     "desc",
		"page_type": "search",
		"scenario":  "PAGE_GLOBAL_SEARCH",
		"version":   "2",
	}
	url := GenerateParams(preUrl, params)
	response, error := http.Get(url)

	if error != nil {
		panic(error.Error())
	}

	body, _ := ioutil.ReadAll(response.Body)

	return string(body)
}

func GenerateParams(url string, params map[string]string) string {
	if len(params) == 0 {
		return url
	}

	url += "?"
	for key := range params {
		url += key + "=" + params[key] + "&"
	}
	url = url[:len(url)-1]

	return url
}

type CrawlData struct {
	NoMore bool
	Items  []Item
}

type Item struct {
	ItemBasic ItemBasic `json:"item_basic"`
}

type ItemBasic struct {
	ItemId      int
	ShopId      int
	Name        string
	Description string
	Price       int64
	Currency    string
	Images      []string
}
