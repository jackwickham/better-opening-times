package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var Templates = template.Must(template.ParseGlob("templates/*.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) == 0 {
		CenterSelectionHandler(w, r)
		return
	}

	urlComponents := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(urlComponents) == 1 {
		ActivitySelectionHandler(urlComponents[0], w, r)
	} else if len(urlComponents) == 2 {
		request := GetOpeningTimesRequest{
			venue:    urlComponents[0],
			activity: urlComponents[1],
		}
		OpeningTimesHandler(request, w, r)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	basePath := "/better-opening-times/"

	mux.Handle(basePath, http.StripPrefix(basePath, http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8072", mux))
}
