package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ISO8601 = "2006-01-02T15:04:05.99999"
)

// getJSONFromURL fetches a JSON object from a given URL and returns it as an
// interface.
func getJSONFromURL(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	// Confirm that the response body is JSON before decoding. We need to add the text/plain
	// content type to the list since we're fetching from a GitHub raw URL.
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "text/plain; charset=utf-8" {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return data, nil
}
