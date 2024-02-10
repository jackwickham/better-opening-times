package main

import "net/http"

type originInjectingRoundTripper struct {
	delegate http.RoundTripper
}

func (o *originInjectingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rather than using CORS, the API just 404s if you try to use it without the expected origin
	if req.URL.Hostname() == "better-admin.org.uk" {
		req.Header.Add("Origin", "https://bookings.better.org.uk")
	}
	return o.delegate.RoundTrip(req)
}

func MakeHttpClient() *http.Client {
	return &http.Client{
		Transport: &originInjectingRoundTripper{
			delegate: http.DefaultTransport,
		},
	}
}
