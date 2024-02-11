package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Activity struct {
	Name     string
	Slug     string
	Children []activitySummary
}

type MaybeActivity struct {
	activity Activity
	err      error
}

func LoadActivityDetailsToChan(client *http.Client, venueSlug string, activitySlug string, results chan<- MaybeActivity) {
	resp, err := client.Get(fmt.Sprintf("https://better-admin.org.uk/api/activities/venue/%s/categories/%s", venueSlug, activitySlug))
	if err != nil {
		log.Println("Failed to load activity: ", err)
		results <- MaybeActivity{err: err}
	}
	defer resp.Body.Close()

	var response struct{ Data Activity }
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Failed to parse activity: ", err)
		results <- MaybeActivity{err: err}
	}

	a := response.Data
	if len(response.Data.Children) == 0 {
		// If no actual children, pretend it's a child of itself
		a.Children = []activitySummary{{
			Name: a.Name,
			Slug: a.Slug,
		}}
	}
	results <- MaybeActivity{activity: a}
}
