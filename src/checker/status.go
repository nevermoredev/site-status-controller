package checker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)


type Settings struct {
	Timer int
	Timeout   int
	File      string
	Useragent string
}

func TestSite(u string, s Settings) bool {
	if u == "" {
		return false
	}
	client := http.Client{
		Timeout: time.Duration(s.Timeout) * time.Second,
	}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Add("User-Agent", s.Useragent)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
		return false
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		if resp.StatusCode == 200 {
			fmt.Println(u + ": " + "ok")
		} else {
			fmt.Println("Site ", u, " returned error code: ", resp.StatusCode)
		}
	}
	defer resp.Body.Close()
	return true
}
