package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func BenchmarkPrint(b *testing.B) {
	// To run: go test -bench=. -benchtime=300x
	fmt.Println()
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxConnsPerHost:     200,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     time.Duration(10) * time.Second,
		},
	}

	b.ResetTimer()
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			account, password := "user"+strconv.Itoa(rand.Intn(201)), []byte("123")
			md5 := md5.New()
			md5.Write(password)
			encodePassword := hex.EncodeToString(md5.Sum(nil))

			requestBody := fmt.Sprintf(`{"account":"%s","password": "%s"}`, account, encodePassword)
			var jsonStr = []byte(requestBody)
			request, _ := http.NewRequest("POST", "http://localhost:49/login", bytes.NewBuffer(jsonStr))
			request.Header.Add("Content-type", "application/json")
			_, err := client.Do(request)
			if err != nil {
				log.Println(err)
				panic(err)
			}
		}
	})
}
