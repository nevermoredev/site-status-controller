package checker

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)
func CheckTitle(urlNow string) string {

	res, err := http.Get(urlNow)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc.Find("title").Text()
}