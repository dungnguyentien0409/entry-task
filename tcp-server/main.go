package main

import (
	"database/sql"
	"entrytask/protocol"
	"entrytask/tcp-server/config"
	"entrytask/tcp-server/db"
	"entrytask/tcp-server/handler"
	"entrytask/tcp-server/services"
	"entrytask/tcp-server/shared/constants"
	"entrytask/tcp-server/shared/dto"
	myCache "entrytask/tcp-server/shared/redis"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func init() {
	config.Configure()
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	signBytes, _ := ioutil.ReadFile(constants.PATH + "/tcp-server/jwtRS256.key")

	dbConfig := dto.DBConfig{
		ConnectionString: os.Getenv("DB_CONNECTION"),
		DriverName:       "mysql",
		MaxIdleConns:     6,
		MaxOpenConns:     6,
		ConnMaxLifeTime:  time.Minute,
	}
	UserDB := StartDB(dbConfig)
	defer CloseDB(UserDB)

	userRepository := dao.NewUserRepository(UserDB)
	userCache := myCache.NewCache(os.Getenv("REDIS_URL"))
	userService := services.NewTCPService(userRepository, userCache, signBytes)
	userHandler := handler.NewHandler(userService)

	server := protocol.Server{
		Addr: os.Getenv("TCP_PORT"),
		Handler: func(request []byte) (response []byte, err error) {
			return userHandler.OnHandle(request)
		},
	}

	if server.Start() != nil {
		log.Printf("Error while starting server %+v", err)
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
