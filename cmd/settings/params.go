package settings

import (
	"fmt"
	"os"
	"strconv"
)

func GetParamsFromArgs(s Settings) Settings {
	if len(os.Args) > 1 {
		for i, _ := range os.Args {
			if os.Args[i] == "timer" {
				i++
				s.Timer, _ = strconv.Atoi(os.Args[i])
				fmt.Println("use timer: ", os.Args[i])
			}
			if os.Args[i] == "timeout" {
				i++
				s.Timeout, _ = strconv.Atoi(os.Args[i])
				fmt.Println("use timeout: ", os.Args[i])
			}
			if os.Args[i] == "file" {
				i++
				s.File = os.Args[i]
				fmt.Println("use file: ", os.Args[i])
			}
			if os.Args[i] == "useragent" {
				i++
				s.Useragent = os.Args[i]
				fmt.Println("use useragent: ", os.Args[i])
			}
		}
	}
	return s
}
