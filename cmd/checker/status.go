package checker

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"zeithub.com/site-status-controller/pkg/config/protobuf"
)

type Settings struct {
	Timer     int
	Timeout   int
	File      string
	Useragent string
}

type Response struct {
	Uuid string
	Url      string
	Status   int
	Hash     string
	TitleNow string
}

func TestSite(UuidNow string, urlNow string) *RmqProto.BotJobResponse {
	res, err := http.Get(urlNow)
	status:=true

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
	ContentNow := doc.Find("body").Text()
	hasher := sha1.New()
	hasher.Write([]byte(ContentNow))
	HashNow := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	HashNow = string(HashNow)
	StatusNow := res.StatusCode

	if StatusNow == 200{
		status = true
	}else{
		status = false
	}

	response:=&RmqProto.BotJobResponse{
		Uuid: UuidNow,
		PageUrl: urlNow,
		Status: status,

	}

	log.Printf("Received a message: %s ,%s ", TitleNow , UuidNow)


	return response
}
