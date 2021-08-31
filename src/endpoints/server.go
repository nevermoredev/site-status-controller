package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)




func UpServer() *mux.Route {
	r := mux.NewRouter()
	r.HandleFunc("/GetTitle", func(w http.ResponseWriter, r *http.Request){
		idNow := r.URL.Query().Get("id")
		urlNow := r.URL.Query().Get("url")
		var responseLocaleMap [3]string
		if len(idNow) > 0 && len(urlNow) > 0 {
			titleResponse, idResponse, urlResponse := GetTitle(idNow, urlNow)
			responseLocaleMap[0] = titleResponse
			responseLocaleMap[1] = idResponse
			responseLocaleMap[2] = urlResponse
		}
		buf ,err := json.Marshal(responseLocaleMap)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(buf)
	})

	//r.HandleFunc("/GetStatus", func(w http.ResponseWriter, r *http.Request){
	//
	//	var sitesMap map[string]string
	//	var responseLocaleMap map[string]int
	//	sitesMap["site"] = "dsadasdasdas"
	//	sitesMap["site2"] = "dsadsadsadasdsa"
	//
	//	responseLocaleMap = GetStatus(sitesMap)
	//	if len(responseLocaleMap) > 0{
	//		fmt.Printf("%s", responseLocaleMap)
	//	}
	//})
	//fmt.Print("Server up...")
	//log.Fatal(http.ListenAndServe("localhost:26009", r))

	return nil
}