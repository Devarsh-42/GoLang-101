package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ============================================================
// weather.go - Fetches weather data from WeatherAPI.com
// Using the FREE Current Weather API
// Docs: https://www.weatherapi.com/docs/
// ============================================================

// Base URL for WeatherAPI.com current weather endpoint
const weatherAPIBaseURL = "http://api.weatherapi.com/v1/current.json"

// fetchWeather calls WeatherAPI.com and returns parsed weather data
// Parameters:
//   - apiKey: your WeatherAPI.com API key (free plan works!)
//   - query: location query — can be lat,lon / city name / zip code / IP
//   - lang: language code for descriptions (e.g. "en", "hi") — optional
func fetchWeather(apiKey, query, lang string) (*WeatherResponse, error) {

	// Build the request URL with query parameters
	// q= accepts: "28.6,77.2" or "London" or "10001" etc.
	url := fmt.Sprintf("%s?key=%s&q=%s", weatherAPIBaseURL, apiKey, query)

	// Add optional "lang" parameter if provided
	if lang != "" {
		url += "&lang=" + lang
	}

	// Make the HTTP GET request to WeatherAPI.com
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call WeatherAPI: %w", err)
	}
	defer resp.Body.Close() // Always close the response body when done

	// Read the entire response body into a byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the API returned a non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WeatherAPI error (status %d): %s", resp.StatusCode, string(body))
	}

	// Unmarshal (decode) the JSON response into our Go struct
	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &weather, nil
}
