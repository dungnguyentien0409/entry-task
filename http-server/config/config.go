package config

import (
	"entrytask/tcp-server/shared/constants"
	"io/ioutil"
	"os"
	"strings"
)

func Configure() {
	ReadEnv()
}

func ReadEnv() {
	bytes, _ := ioutil.ReadFile(constants.PATH + "http-server/local.env")

	var values = strings.Split(string(bytes), "\n")
	for _, value := range values {
		var item = strings.Split(value, "=")

		if len(item) >= 2 {
			os.Setenv(item[0], item[1])
		}
	}
}
