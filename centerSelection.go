package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type venuesResponse struct {
	Data []VenueDetails
}

func CenterSelectionHandler(w http.ResponseWriter, r *http.Request) {
	client := MakeHttpClient()
	resp, err := client.Get("https://better-admin.org.uk/api/activities/venues")
	if err != nil {
		log.Println("Failed to load venues: ", err)
		http.Error(w, "Failed to load venues", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var response venuesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Failed to parse venue response: ", err)
		http.Error(w, "Failed to parse venues response", http.StatusInternalServerError)
		return
	}

	err = Templates.ExecuteTemplate(w, "center-selection.html", response.Data)
	if err != nil {
		log.Println("Failed to execute template: ", err)
		http.Error(w, "Failed to render output", http.StatusInternalServerError)
	}
}
