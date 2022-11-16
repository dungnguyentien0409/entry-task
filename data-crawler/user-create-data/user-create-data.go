package main

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MaxOpenConnsLifetime = 4
	MaxUser              = 10000000
)

func main() {
	db, err := SetUpConnection()
	if err != nil {
		log.Println("Could not connect to database")
		return
	}

	for i := 0; i < MaxUser; i += 10 {
		var usernames []string
		for j := 0; j < 10; j++ {
			if i != 0 || j != 0 {
				username := "user" + strconv.Itoa(i+j)
				usernames = append(usernames, username)
			}
		}
		password := []byte("123")
		md5 := md5.New()
		md5.Write(password)

		//202cb962ac59075b964b07152d234b70
		encodePassword := hex.EncodeToString(md5.Sum(nil))

		rand.Seed(time.Now().UnixNano())
		go InsertUser(usernames, encodePassword, db)
		time.Sleep(20 * time.Millisecond)
	}
}

func SetUpConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/entry_task")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(MaxOpenConnsLifetime * time.Minute)

	return db, err
}

func GenerateSalt() string {
	salt := ""

	for i := 0; i < 64; i++ {
		random_int := rand.Intn(74) + 48
		character := rune(random_int)
		salt += string(character)
	}

	return salt
}

func GenerateHashPassword(password string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(salt))
	hash.Write([]byte(password))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func InsertUser(accounts []string, password string, db *sql.DB) (int8, string, int32, string) {
	salt := GenerateSalt()
	hashPassword := GenerateHashPassword(password, salt)
	insertQuery := "INSERT INTO user_tab (account, nickname, password, salt) VALUES"
	vals := []interface{}{}

	for _, account := range accounts {
		insertQuery += "(?,?,?,?),"
		vals = append(vals, account, account, hashPassword, salt)
	}
	insertQuery = insertQuery[:len(insertQuery)-1]

	//prepare the statement
	statement, _ := db.Prepare(insertQuery)
	//format all vals at once
	_, err := statement.Exec(vals...)

	if err != nil {
		log.Println(err)
		return 1, "Fail to create new user", -1, ""
	}

	return 0, "Create new user successfully", -1, accounts[0]
}
