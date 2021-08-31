package checker

import (
	"crypto"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"hash"
	"log"
	"net/http"
)

type Settings struct {
	Timer     int
	Timeout   int
	File      string
	Useragent string
}

type Response struct {
	Url      string
	Status   int
	Hash     hash.Hash
	TitleNow string
}

func TestSite(urlNow string) Response {

	res, err := http.Get(urlNow)
	if err != nil {
		log.Printf("%s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("%s", err)
	}

	TitleNow := doc.Find("title").Text()
	//ContentNow := doc.Find("body").Text()
	checksum := crypto.MD5
	HashNow := checksum.New()
	//HashNow, err = HashNow.Write([]byte(ContentNow))
	StatusNow := res.StatusCode

	fmt.Printf("Checksum 1: %s\n", checksum)

	//log.Printf("Received a message: %s", title)

	responseNow := Response{urlNow, StatusNow, HashNow, TitleNow}

	return responseNow
}
