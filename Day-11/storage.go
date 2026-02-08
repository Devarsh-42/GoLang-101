package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ============================================================
// storage.go - Save and load weather data to/from a JSON file
// Uses json.Marshal (struct -> JSON) and json.Unmarshal (JSON -> struct)
// ============================================================

const weatherFile = "weather.json" // File where we store weather data locally

// saveToJSON takes weather data and writes it to weather.json
// json.MarshalIndent makes the output human-readable with indentation
func saveToJSON(data *WeatherResponse) error {

	// Marshal (encode) the struct into pretty-printed JSON bytes
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal weather data: %w", err)
	}

	// Write the JSON bytes to the file (0644 = owner read/write, others read)
	err = os.WriteFile(weatherFile, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", weatherFile, err)
	}

	return nil
}

// loadFromJSON reads weather.json and returns the parsed weather data
func loadFromJSON() (*WeatherResponse, error) {

	// Read the entire file into a byte slice
	jsonBytes, err := os.ReadFile(weatherFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", weatherFile, err)
	}

	// Unmarshal (decode) the JSON bytes into our Go struct
	var data WeatherResponse
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", weatherFile, err)
	}

	return &data, nil
}
