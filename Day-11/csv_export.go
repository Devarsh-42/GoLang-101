package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ============================================================
// csv_export.go - Export weather data to a CSV file
// ============================================================

const csvFile = "weather.csv" // Output CSV file name

// exportToCSV reads the saved weather.json and writes a CSV file from it
func exportToCSV() error {

	// First, load the weather data from our local JSON file
	data, err := loadFromJSON()
	if err != nil {
		return fmt.Errorf("cannot export CSV - no saved data: %w", err)
	}

	// Create (or overwrite) the CSV file
	file, err := os.Create(csvFile)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Make sure all buffered data is written

	// Write the CSV header row
	header := []string{
		"City", "Region", "Country", "Latitude", "Longitude",
		"Localtime", "TempC", "TempF", "FeelsLikeC",
		"Humidity", "PressureMb", "WindKph", "WindDir",
		"Cloud", "VisKm", "UV", "Condition",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write one data row with the current weather
	row := []string{
		data.Location.Name,
		data.Location.Region,
		data.Location.Country,
		fmt.Sprintf("%.4f", data.Location.Lat),
		fmt.Sprintf("%.4f", data.Location.Lon),
		data.Location.Localtime,
		fmt.Sprintf("%.1f", data.Current.TempC),
		fmt.Sprintf("%.1f", data.Current.TempF),
		fmt.Sprintf("%.1f", data.Current.FeelsLikeC),
		strconv.Itoa(data.Current.Humidity),
		fmt.Sprintf("%.0f", data.Current.PressureMb),
		fmt.Sprintf("%.1f", data.Current.WindKph),
		data.Current.WindDir,
		strconv.Itoa(data.Current.Cloud),
		fmt.Sprintf("%.1f", data.Current.VisKm),
		fmt.Sprintf("%.1f", data.Current.UV),
		data.Current.Condition.Text,
	}
	if err := writer.Write(row); err != nil {
		return fmt.Errorf("failed to write CSV row: %w", err)
	}

	return nil
}
