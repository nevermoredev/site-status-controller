package checker

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	RmqProto "github.com/zeithub/site-status-controller/pkg/config/protobuf"
)

type Settings struct {
	Timer     int
	Timeout   int
	File      string
	Useragent string
}

func TestSite(PageId string, urlNow string, oldTitle string, Action uint32) *RmqProto.BotJobResponse {

	res, errHttp := http.Get(urlNow)
	response := &RmqProto.BotJobResponse{}
	status := true

	if errHttp != nil {
		log.Printf("Cant get resource: %v", errHttp)
		return nil
	}

	doc, errDoc := goquery.NewDocumentFromReader(res.Body)
	TitleNow := doc.Find("title").Text()

	if errDoc != nil {
		log.Printf("Cant get resource: %v", errDoc)
		return nil
	}

	// Action = 1 (New) if Action = 2(Old)
	switch Action {

	case 1:
		ContentNow := doc.Find("body").Text()
		hasher := sha1.New()
		hasher.Write([]byte(ContentNow))
		HashNow := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		response = &RmqProto.BotJobResponse{
			PageId:  PageId,
			PageUrl: urlNow,
			Status:  status,
			Title:   TitleNow,
			Hash:    HashNow,
		}

	case 2:
		if res.StatusCode != 200 || res.StatusCode != 500 || res.StatusCode != 501 || res.StatusCode != 504 {
			status = false

			log.Print("status code error: %d %s", res.StatusCode, res.Status)
			response = &RmqProto.BotJobResponse{
				PageId:  PageId,
				PageUrl: urlNow,
				Status:  status,
			}

		} else {
			status = true
			// Similarity test

			swg := metrics.NewSmithWatermanGotoh()
			similarity := strutil.Similarity(oldTitle, TitleNow, swg)
			if similarity < 70 {
				status = false
				log.Print("status code error: %d %s", res.StatusCode, res.Status)
				response = &RmqProto.BotJobResponse{
					PageId:  PageId,
					PageUrl: urlNow,
					Status:  status,
				}
			} else {
				ContentNow := doc.Find("body").Text()
				hasher := sha1.New()
				hasher.Write([]byte(ContentNow))
				HashNow := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

				response = &RmqProto.BotJobResponse{
					PageId:  PageId,
					PageUrl: urlNow,
					Status:  status,
					Title:   TitleNow,
					Hash:    HashNow,
				}
			}

			log.Printf("Received a message: %s, %s ", TitleNow, PageId)
		}

	}

	defer res.Body.Close()
	return response
}
