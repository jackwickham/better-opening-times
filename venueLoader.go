package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type VenueDetails struct {
	Name string
	Slug string
}

type venueResponse struct {
	Data VenueDetails
}

func LoadVenueDetails(client *http.Client, venueSlug string) (*VenueDetails, error) {
	resp, err := client.Get(fmt.Sprintf("https://better-admin.org.uk/api/activities/venues/%s", venueSlug))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response venueResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}

func LoadVenueDetailsToChan(client *http.Client, venueSlug string, c chan<- *VenueDetails) {
	details, err := LoadVenueDetails(client, venueSlug)
	if err != nil {
		log.Println("Failed to load venue details: ", err)
		close(c)
		return
	}
	c <- details
}
