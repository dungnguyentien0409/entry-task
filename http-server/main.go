package main

import (
	"database/sql"
	"entrytask/http-server/config"
	"entrytask/http-server/db"
	"entrytask/http-server/handler"
	"entrytask/http-server/product_service"
	"entrytask/http-server/user_service"
	"entrytask/tcp-server/shared/constants"
	"entrytask/tcp-server/shared/dto"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	config.Configure()
}

func main() {
	userService, err := user_service.NewUserService(200, 300, os.Getenv("TCP_PORT"))

	if err != nil {
		log.Println(err)
	}

	dbConfig := dto.DBConfig{
		ConnectionString: os.Getenv("DB_CONNECTION"),
		DriverName:       "mysql",
		MaxIdleConns:     6,
		MaxOpenConns:     6,
		ConnMaxLifeTime:  time.Minute,
	}
	ProductDB := StartDB(dbConfig)
	defer CloseDB(ProductDB)

	productRepository := dao.NewProductRepository(ProductDB)
	productService := product_service.NewProductService(productRepository)
	publicBytes, _ := ioutil.ReadFile(constants.PATH + "/http-server/jwtRS256.key.pub")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicBytes)

	userHandler := handler.NewHandler(userService, productService, publicKey)

	http.HandleFunc("/register", userHandler.CORSMiddleware(userHandler.HandleRegister))
	http.HandleFunc("/login", userHandler.CORSMiddleware(userHandler.HandleLogin))
	http.HandleFunc("/ping", userHandler.CORSMiddleware(userHandler.HandlePing))

	http.Handle("/products", userHandler.CORSMiddleware(userHandler.AuthMiddleware(userHandler.HandleGetListProducts)))
	http.Handle("/products/", userHandler.CORSMiddleware(userHandler.AuthMiddleware(userHandler.HandleGetProduct)))
	http.Handle("/comments", userHandler.CORSMiddleware(userHandler.AuthMiddleware(userHandler.HandlePostComment)))
	http.Handle("/comments/", userHandler.CORSMiddleware(userHandler.AuthMiddleware(userHandler.HandleGetComment)))

	if err := http.ListenAndServe(os.Getenv("HTTP_PORT"), nil); err != nil {
		log.Printf("HTTP server was shutdown: %v\n", err)
	}
}

func StartDB(cfg dto.DBConfig) *sql.DB {
	var err error
	DB, err := sql.Open(cfg.DriverName, cfg.ConnectionString)
	if err != nil {
		log.Fatalf("could not connect to database: %s\n", err)
	}

	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)

	err = DB.Ping()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("connected to database!")
	return DB
}

func CloseDB(DB *sql.DB) {
	err := DB.Close()
	if err != nil {
		fmt.Printf("could not close connection to db: %s\n", err)
	}
	fmt.Println("close connection is done!")
}
