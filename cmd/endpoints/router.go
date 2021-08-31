package endpoints

import (
	"fmt"
	"time"
	"zeithub.com/site-status-controller/src/checker"
	"zeithub.com/site-status-controller/src/resource"
	"zeithub.com/site-status-controller/src/settings"
)

func GetTitle(id string, url string) (string, string, string) {

	data := checker.CheckTitle(url)

	return data, id, url
}

func GetStatus(data map[string]string) map[string]uint16 {

	var sitesMap map[string]uint16

	fmt.Print(data)
	sitesMap["site"] = 200
	sitesMap["site2"] = 200

	sett := settings.InitSettings()
	s := settings.GetParamsFromArgs(sett)
	sites := resource.GetSitesFromFile(s.File)
	ch := make(chan byte, 1)
	for _, site := range sites {
		go func(site string) {
			for {
				checker.TestSite(site, checker.Settings(sett))
				time.Sleep(time.Duration(s.Timer) * time.Second)
			}
			ch <- 1
		}(site)
	}
	<-ch

	return sitesMap
}
