package resource

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetSitesFromFile(namefile string) []string {
	data, err := ioutil.ReadFile(namefile)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(2)
	}
	result := string(data)

	test := strings.Split(result, "\r\n")
	return test
}
