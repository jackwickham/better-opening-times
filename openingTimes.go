package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

type GetOpeningTimesRequest struct {
	venue    string
	activity string
}

type dateInfo struct {
	Date   string `json:"raw"`
	Pretty string `json:"full_date_pretty"`
}

type datesResponse struct {
	Data []dateInfo
}

type openingTime struct {
	Duration       string `json:"duration"`
	StartTimestamp int64  `json:"timestamp"`
}

type openingTimesResponse struct {
	Data interface{} `json:"data"` // data can either be a map or an array
}

var durationRegex = regexp.MustCompile("(\\d+)min")

func loadDates(client *http.Client, request GetOpeningTimesRequest) (*datesResponse, error) {
	url := fmt.Sprintf("https://better-admin.org.uk/api/activities/venue/%s/activity-category/%s/dates", request.venue, request.activity)
	resp, err := client.Get(url)
	if err != nil {
		log.Println("Failed to load dates:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var dates datesResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&dates); err != nil {
		log.Println("Failed to decode dates response:", err)
		return nil, err
	}
	return &dates, nil
}

func loadOpeningTimes(client *http.Client, date dateInfo, request GetOpeningTimesRequest) (*RangeSet, error) {
	url := fmt.Sprintf("https://better-admin.org.uk/api/activities/venue/%s/activity/%s/times?date=%s", request.venue, request.activity, date.Date)
	resp, err := client.Get(url)
	if err != nil {
		log.Println("Failed to load times:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var timesResponse openingTimesResponse
	if err := json.NewDecoder(resp.Body).Decode(&timesResponse); err != nil {
		log.Println("Failed to decode times response:", err)
		return nil, err
	}

	times := make([]openingTime, 0)
	switch v := timesResponse.Data.(type) {
	case []interface{}:
		for _, item := range v {
			serialized, _ := json.Marshal(item)
			var deserialized openingTime
			json.Unmarshal(serialized, &deserialized)
			times = append(times, deserialized)
		}
	case map[string]interface{}:
		for _, item := range v {
			serialized, _ := json.Marshal(item)
			var deserialized openingTime
			json.Unmarshal(serialized, &deserialized)
			times = append(times, deserialized)
		}
	default:
		log.Printf("Unexpected return type for times: %T\n", v)
		return nil, errors.New("Unepected times format")
	}

	slices.SortFunc(times, func(a, b openingTime) int {
		return cmp.Compare(a.StartTimestamp, b.StartTimestamp)
	})
	rangeSet := NewRangeSet()
	for _, slot := range times {
		matches := durationRegex.FindStringSubmatch(slot.Duration)
		if matches == nil {
			log.Printf("Failed to parse %s as duration\n", slot.Duration)
			return nil, errors.New("Failed to parse duration")
		}
		durationMins, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			log.Printf("Failed to parse duration int\n")
			return nil, err
		}
		durationSec := durationMins * 60
		rangeSet.Add(Range{
			StartTimestamp: slot.StartTimestamp,
			EndTimestamp:   slot.StartTimestamp + durationSec,
		})
	}

	return &rangeSet, nil
}

type outputTimeInfo struct {
	Start string
	End   string
}

type outputDateInfo struct {
	Date  string
	Times []outputTimeInfo
}

type openingTimesOutput struct {
	Dates    []outputDateInfo
	Venue    *VenueDetails
	Activity Activity
}

func OpeningTimesHandler(request GetOpeningTimesRequest, w http.ResponseWriter, r *http.Request) {
	client := MakeHttpClient()

	venueChan := make(chan *VenueDetails, 1)
	go LoadVenueDetailsToChan(client, request.venue, venueChan)

	activityChan := make(chan MaybeActivity, 1)
	go LoadActivityDetailsToChan(client, request.venue, request.activity, activityChan)

	dates, err := loadDates(client, request)
	if err != nil {
		http.Error(w, "Failed to load available dates", http.StatusInternalServerError)
		return
	}

	if len(dates.Data) > 21 {
		dates.Data = dates.Data[:21]
	}
	results := make([]outputDateInfo, len(dates.Data))

	location, _ := time.LoadLocation("Europe/London")

	eg := new(errgroup.Group)
	for ti, tdate := range dates.Data {
		i, date := ti, tdate
		eg.Go(func() error {
			timeRanges, err := loadOpeningTimes(client, date, request)
			if err != nil {
				return err
			}

			times := make([]outputTimeInfo, len(timeRanges.Ranges))
			for j, timeRange := range timeRanges.Ranges {
				times[j] = outputTimeInfo{
					Start: time.Unix(timeRange.StartTimestamp, 0).In(location).Format("15:04"),
					End:   time.Unix(timeRange.EndTimestamp, 0).In(location).Format("15:04"),
				}
			}

			results[i] = outputDateInfo{
				Date:  date.Pretty,
				Times: times,
			}

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		http.Error(w, "Failed to load times", http.StatusInternalServerError)
		return
	}

	venueDetails, ok := <-venueChan
	if !ok {
		http.Error(w, "Failed to load venue details", http.StatusInternalServerError)
		return
	}

	activityDetails := <-activityChan
	if activityDetails.err != nil {
		http.Error(w, "Failed to load activity details", http.StatusInternalServerError)
		return
	}

	outputData := openingTimesOutput{
		Dates:    results,
		Venue:    venueDetails,
		Activity: activityDetails.activity,
	}

	err = Templates.ExecuteTemplate(w, "opening-times.html", outputData)
	if err != nil {
		log.Println("Failed to execute template: ", err)
		http.Error(w, "Failed to render output", http.StatusInternalServerError)
	}
}
