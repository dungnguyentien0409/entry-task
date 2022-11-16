package main

import (
	"fmt"
	_ "net/http/pprof"
	"time"
)

func main() {
	fmt.Println(time.Now())
}
