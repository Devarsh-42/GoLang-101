package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// ============================================================
// Tests for HTTP handlers and storage functions (WeatherAPI.com)
// ============================================================

// TestHandleHello tests the welcome endpoint
func TestHandleHello(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	handleHello(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK but got %v", resp.Status)
	}
	expected := `{"message": "Welcome to Day-11 Weather API!"}`
	if string(body) != expected {
		t.Errorf("Expected body '%s' but got '%s'", expected, string(body))
	}
}

// TestHandleGoodBye tests the goodbye endpoint
func TestHandleGoodBye(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/goodbye", nil)
	w := httptest.NewRecorder()

	handleGoodBye(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK but got %v", resp.Status)
	}
	expected := `{"message": "Goodbye, World!"}`
	if string(body) != expected {
		t.Errorf("Expected body '%s' but got '%s'", expected, string(body))
	}
}

// TestHandleWeatherAPI_MissingParams tests that missing query params return 400
func TestHandleWeatherAPI_MissingParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/weather", nil) // No query params
	w := httptest.NewRecorder()

	handleWeatherAPI(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

// TestSaveAndLoadJSON tests the marshal/unmarshal round trip
func TestSaveAndLoadJSON(t *testing.T) {
	// Create sample weather data (WeatherAPI.com structure)
	sample := &WeatherResponse{
		Location: Location{
			Name:    "New Delhi",
			Region:  "Delhi",
			Country: "India",
			Lat:     28.6139,
			Lon:     77.2090,
			TzID:    "Asia/Kolkata",
		},
		Current: CurrentWeather{
			TempC:      30.5,
			TempF:      86.9,
			Humidity:   60,
			FeelsLikeC: 33.0,
			Condition:  Condition{Text: "Partly cloudy", Code: 1003},
		},
	}

	// Save to JSON (uses marshal)
	err := saveToJSON(sample)
	if err != nil {
		t.Fatalf("saveToJSON failed: %v", err)
	}

	// Load back from JSON (uses unmarshal)
	loaded, err := loadFromJSON()
	if err != nil {
		t.Fatalf("loadFromJSON failed: %v", err)
	}

	// Verify the data matches
	if loaded.Location.Lat != sample.Location.Lat || loaded.Location.Lon != sample.Location.Lon {
		t.Errorf("Lat/Lon mismatch: got %.4f,%.4f want %.4f,%.4f",
			loaded.Location.Lat, loaded.Location.Lon, sample.Location.Lat, sample.Location.Lon)
	}
	if loaded.Current.TempC != sample.Current.TempC {
		t.Errorf("TempC mismatch: got %.2f want %.2f", loaded.Current.TempC, sample.Current.TempC)
	}
	if loaded.Location.Name != "New Delhi" {
		t.Errorf("City name mismatch: got %s want New Delhi", loaded.Location.Name)
	}

	// Clean up the test file
	os.Remove(weatherFile)
}

// TestParseFlags tests the CLI flag parser
func TestParseFlags(t *testing.T) {
	args := []string{"--key=abc123", "--q=London", "--lang=en"}
	flags := parseFlags(args)

	if flags["key"] != "abc123" {
		t.Errorf("Expected key=abc123, got %s", flags["key"])
	}
	if flags["q"] != "London" {
		t.Errorf("Expected q=London, got %s", flags["q"])
	}
	if flags["lang"] != "en" {
		t.Errorf("Expected lang=en, got %s", flags["lang"])
	}
}

// TestHandleSavedData_NoFile tests that /weather/local GET returns 404 when no file exists
func TestHandleSavedData_NoFile(t *testing.T) {
	// Make sure no weather.json exists
	os.Remove(weatherFile)

	req := httptest.NewRequest("GET", "/weather/local", nil)
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// TestLocalPost tests POST /weather/local creates weather.json
func TestLocalPost(t *testing.T) {
	os.Remove(weatherFile) // Make sure file doesn't exist

	body := `{"location":{"name":"Mumbai","country":"India","lat":19.07,"lon":72.87},"current":{"temp_c":32.0}}`
	req := httptest.NewRequest("POST", "/weather/local", strings.NewReader(body))
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d", w.Code)
	}

	// Verify file was created
	loaded, err := loadFromJSON()
	if err != nil {
		t.Fatalf("Failed to load created file: %v", err)
	}
	if loaded.Location.Name != "Mumbai" {
		t.Errorf("Expected city Mumbai but got %s", loaded.Location.Name)
	}

	os.Remove(weatherFile)
}

// TestLocalPost_Conflict tests POST returns 409 if weather.json already exists
func TestLocalPost_Conflict(t *testing.T) {
	// Create the file first
	sample := &WeatherResponse{Location: Location{Name: "Delhi"}}
	saveToJSON(sample)

	body := `{"location":{"name":"London"},"current":{"temp_c":15.0}}`
	req := httptest.NewRequest("POST", "/weather/local", strings.NewReader(body))
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 but got %d", w.Code)
	}

	os.Remove(weatherFile)
}

// TestLocalPut tests PUT /weather/local replaces weather.json
func TestLocalPut(t *testing.T) {
	// Create initial data
	sample := &WeatherResponse{Location: Location{Name: "Delhi"}, Current: CurrentWeather{TempC: 30}}
	saveToJSON(sample)

	// PUT new data
	body := `{"location":{"name":"Paris","country":"France","lat":48.85,"lon":2.35},"current":{"temp_c":18.0}}`
	req := httptest.NewRequest("PUT", "/weather/local", strings.NewReader(body))
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}

	// Verify it was replaced
	loaded, _ := loadFromJSON()
	if loaded.Location.Name != "Paris" {
		t.Errorf("Expected city Paris but got %s", loaded.Location.Name)
	}

	os.Remove(weatherFile)
}

// TestLocalPatch tests PATCH /weather/local updates specific fields
func TestLocalPatch(t *testing.T) {
	// Create initial data
	sample := &WeatherResponse{
		Location: Location{Name: "Delhi", Country: "India", Lat: 28.61},
		Current:  CurrentWeather{TempC: 30, Humidity: 50},
	}
	saveToJSON(sample)

	// PATCH only the temperature
	body := `{"current":{"temp_c":25.0}}`
	req := httptest.NewRequest("PATCH", "/weather/local", strings.NewReader(body))
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}

	// Verify temp changed but city stayed
	loaded, _ := loadFromJSON()
	if loaded.Current.TempC != 25.0 {
		t.Errorf("Expected temp 25.0 but got %f", loaded.Current.TempC)
	}
	if loaded.Location.Name != "Delhi" {
		t.Errorf("Expected city Delhi but got %s", loaded.Location.Name)
	}

	os.Remove(weatherFile)
}

// TestLocalDelete tests DELETE /weather/local removes weather.json
func TestLocalDelete(t *testing.T) {
	// Create the file
	sample := &WeatherResponse{Location: Location{Name: "Delhi"}}
	saveToJSON(sample)

	req := httptest.NewRequest("DELETE", "/weather/local", nil)
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}

	// Verify file is gone
	if _, err := os.Stat(weatherFile); !os.IsNotExist(err) {
		t.Error("weather.json should have been deleted")
		os.Remove(weatherFile)
	}
}

// TestLocalDelete_NotFound tests DELETE returns 404 if no file exists
func TestLocalDelete_NotFound(t *testing.T) {
	os.Remove(weatherFile)

	req := httptest.NewRequest("DELETE", "/weather/local", nil)
	w := httptest.NewRecorder()

	handleLocal(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// TestLoggingMiddleware tests that the logging middleware passes through requests
func TestLoggingMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := loggingMiddleware(inner)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

// TestJSONContentTypeMiddleware tests that the middleware sets correct Content-Type
func TestJSONContentTypeMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := jsonContentTypeMiddleware(inner)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json but got %s", contentType)
	}
}

// TestExportToCSV tests the CSV export produces a valid file
func TestExportToCSV(t *testing.T) {
	// First save some test data (WeatherAPI.com structure)
	sample := &WeatherResponse{
		Location: Location{
			Name:      "New Delhi",
			Region:    "Delhi",
			Country:   "India",
			Lat:       28.6139,
			Lon:       77.2090,
			Localtime: "2025-02-09 12:00",
		},
		Current: CurrentWeather{
			TempC:      30.5,
			TempF:      86.9,
			Humidity:   60,
			PressureMb: 1012,
			WindKph:    15.0,
			WindDir:    "NW",
			Cloud:      25,
			VisKm:      10.0,
			UV:         5.0,
			Condition:  Condition{Text: "Partly cloudy"},
		},
	}
	saveToJSON(sample)

	// Export to CSV
	err := exportToCSV()
	if err != nil {
		t.Fatalf("exportToCSV failed: %v", err)
	}

	// Verify the CSV file was created
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Error("weather.csv was not created")
	}

	// Clean up
	os.Remove(weatherFile)
	os.Remove(csvFile)
}

// TestMarshalUnmarshal tests direct JSON marshal/unmarshal
func TestMarshalUnmarshal(t *testing.T) {
	original := WeatherResponse{
		Location: Location{
			Name:    "London",
			Country: "United Kingdom",
			Lat:     51.52,
			Lon:     -0.11,
		},
		Current: CurrentWeather{
			TempC:      22.5,
			FeelsLikeC: 21.0,
			Humidity:   55,
			PressureMb: 1013,
			WindKph:    12.5,
		},
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal back
	var decoded WeatherResponse
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify fields match
	if decoded.Location.Lat != original.Location.Lat {
		t.Errorf("Lat mismatch: got %f want %f", decoded.Location.Lat, original.Location.Lat)
	}
	if decoded.Current.TempC != original.Current.TempC {
		t.Errorf("TempC mismatch: got %f want %f", decoded.Current.TempC, original.Current.TempC)
	}
}
