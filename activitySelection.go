package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
)

type activitySummary struct {
	Name string
	Slug string
}

type activityOuputData struct {
	VenueDetails *VenueDetails
	Activities   []Activity
}

func ActivitySelectionHandler(center string, w http.ResponseWriter, r *http.Request) {
	client := MakeHttpClient()

	venueChan := make(chan *VenueDetails, 1)
	go LoadVenueDetailsToChan(client, center, venueChan)

	resp, err := client.Get(fmt.Sprintf("https://better-admin.org.uk/api/activities/venue/%s/categories", center))
	if err != nil {
		log.Println("Failed to load activity categories: ", err)
		http.Error(w, "Failed to load activity categories", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var categories struct{ Data []activitySummary }
	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		log.Println("Failed to parse activity categories response: ", err)
		http.Error(w, "Failed to parse activity categories response", http.StatusInternalServerError)
		return
	}

	activityQueue := make(chan activitySummary, len(categories.Data))
	activityResults := make(chan MaybeActivity, len(categories.Data))

	for i := 0; i < 4; i++ {
		go func(client *http.Client, venueSlug string, queue <-chan activitySummary, results chan<- MaybeActivity) {
			for activityCategory := range queue {
				LoadActivityDetailsToChan(client, venueSlug, activityCategory.Slug, results)
			}
		}(client, center, activityQueue, activityResults)
	}

	for _, activity := range categories.Data {
		activityQueue <- activity
	}
	close(activityQueue)

	results := make([]Activity, len(categories.Data))
	for i := 0; i < len(categories.Data); i++ {
		res := <-activityResults
		if res.err != nil {
			http.Error(w, "Failed to load activity details", http.StatusInternalServerError)
			return
		}
		results[i] = res.activity
	}
	slices.SortFunc(results, func(a Activity, b Activity) int {
		return cmp.Compare(a.Name, b.Name)
	})

	venueDetails, ok := <-venueChan
	if !ok {
		http.Error(w, "Failed to load venue details", http.StatusInternalServerError)
		return
	}

	outputData := activityOuputData{
		VenueDetails: venueDetails,
		Activities:   results,
	}

	err = Templates.ExecuteTemplate(w, "activity-selection.html", outputData)
	if err != nil {
		log.Println("Failed to execute template: ", err)
		http.Error(w, "Failed to render output", http.StatusInternalServerError)
	}
}
