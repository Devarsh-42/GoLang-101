package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// ============================================================
// main.go - Entry point with CLI commands and HTTP server
// Using WeatherAPI.com (free plan)
//
// CLI Usage:
//   go run . serve                              → Start the HTTP server
//   go run . fetch --key=xxx --q=London         → Fetch weather & save to weather.json
//   go run . fetch --key=xxx --q=28.6,77.2      → Fetch by lat,lon
//   go run . show                               → Show saved weather data
//   go run . csv                                → Export saved data to weather.csv
//
// HTTP Endpoints (when server is running):
//   GET    /                        → Welcome message
//   GET    /weather?key=..&q=..     → Fetch weather from API
//   GET    /weather/local           → Read saved weather.json data
//   POST   /weather/local           → Create new weather data in weather.json
//   PUT    /weather/local           → Replace entire weather data in weather.json
//   PATCH  /weather/local           → Update specific fields in weather.json
//   DELETE /weather/local           → Delete weather.json file
// ============================================================

func main() {
	// If no command is given, show usage help
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// The first argument is the command (serve, fetch, show, csv)
	command := os.Args[1]

	switch command {
	case "serve":
		startServer()
	case "fetch":
		handleFetchCLI()
	case "show":
		handleShowCLI()
	case "csv":
		handleCSVCLI()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

// printUsage shows all available CLI commands
func printUsage() {
	fmt.Println("=== Day-11 Weather CLI (WeatherAPI.com) ===")
	fmt.Println("Usage: go run . <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  serve                             Start the HTTP server on :8080")
	fmt.Println("  fetch --key=API_KEY --q=QUERY     Fetch weather and save to weather.json")
	fmt.Println("        [--lang=en]                 QUERY can be: city name, lat,lon, zip, IP")
	fmt.Println("  show                              Display saved weather data from weather.json")
	fmt.Println("  csv                               Export saved weather data to weather.csv")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run . fetch --key=abc123 --q=London")
	fmt.Println("  go run . fetch --key=abc123 --q=28.6139,77.2090")
	fmt.Println("  go run . fetch --key=abc123 --q=10001")
}

// ============================================================
// CLI Command Handlers
// ============================================================

// handleFetchCLI parses CLI flags, fetches weather, and saves to JSON
func handleFetchCLI() {
	// Parse flags from command-line arguments (e.g. --key=abc --q=London)
	flags := parseFlags(os.Args[2:])

	apiKey := flags["key"]
	query := flags["q"]   // Location: city name, lat,lon, zip, IP
	lang := flags["lang"] // Optional: language code (en, hi, etc.)

	// Validate required parameters
	if apiKey == "" || query == "" {
		fmt.Println("Error: --key and --q are required")
		fmt.Println("Example: go run . fetch --key=YOUR_API_KEY --q=London")
		return
	}

	fmt.Printf("Fetching weather for q=%s ...\n", query)

	// Call the WeatherAPI.com API
	weather, err := fetchWeather(apiKey, query, lang)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Save the response to weather.json using Marshal
	err = saveToJSON(weather)
	if err != nil {
		fmt.Printf("Error saving data: %v\n", err)
		return
	}

	// Print a summary to the terminal
	fmt.Println("Weather data saved to weather.json!")
	printWeatherSummary(weather)
}

// handleShowCLI loads and displays saved weather data from weather.json
func handleShowCLI() {
	// Load data from weather.json using Unmarshal
	weather, err := loadFromJSON()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Hint: Run 'go run . fetch --key=... --q=London' first")
		return
	}

	printWeatherSummary(weather)
}

// handleCSVCLI exports saved weather data to a CSV file
func handleCSVCLI() {
	err := exportToCSV()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Weather data exported to weather.csv!")
}

// ============================================================
// HTTP Server & Handlers
// ============================================================

// startServer starts the HTTP server with middleware and routes
func startServer() {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/{$}", handleHello)           // GET / → Welcome message
	mux.HandleFunc("/goodbye", handleGoodBye)     // GET /goodbye → Goodbye message
	mux.HandleFunc("/weather", handleWeatherAPI)  // GET /weather?key=..&q=.. → Fetch weather
	mux.HandleFunc("/weather/local", handleLocal) // GET/POST/PUT/PATCH/DELETE → CRUD on weather.json

	// Wrap the mux with middleware (logging → JSON content type → handler)
	handler := loggingMiddleware(jsonContentTypeMiddleware(mux))

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /                       → Welcome")
	fmt.Println("  GET    /weather?key=..&q=..    → Fetch weather from WeatherAPI.com")
	fmt.Println("  GET    /weather/local           → Read saved weather.json")
	fmt.Println("  POST   /weather/local           → Create new weather data")
	fmt.Println("  PUT    /weather/local           → Replace all weather data")
	fmt.Println("  PATCH  /weather/local           → Update specific fields")
	fmt.Println("  DELETE /weather/local           → Delete weather.json")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("Server failed", "error", err)
	}
}

// handleHello responds with a welcome message
func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Welcome to Day-11 Weather API!"}`))
}

// handleGoodBye responds with a goodbye message
func handleGoodBye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Goodbye, World!"}`))
}

// handleWeatherAPI fetches weather from WeatherAPI.com via HTTP query params
func handleWeatherAPI(w http.ResponseWriter, r *http.Request) {
	// Read query parameters from the URL
	apiKey := r.URL.Query().Get("key")
	query := r.URL.Query().Get("q")   // Location query
	lang := r.URL.Query().Get("lang") // Optional language

	// Validate required params
	if apiKey == "" || query == "" {
		http.Error(w, `{"error": "key and q query params are required"}`, http.StatusBadRequest)
		return
	}

	// Fetch weather data from the API
	weather, err := fetchWeather(apiKey, query, lang)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Save to weather.json locally
	saveToJSON(weather)

	// Marshal the response back to JSON and send it
	jsonBytes, err := json.MarshalIndent(weather, "", "  ")
	if err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}

// ============================================================
// handleLocal routes /weather/local to the right handler based on HTTP method
// This is how we support GET, POST, PUT, PATCH, DELETE on one path
// ============================================================
func handleLocal(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleLocalGet(w, r) // Read weather.json
	case http.MethodPost:
		handleLocalPost(w, r) // Create new weather data
	case http.MethodPut:
		handleLocalPut(w, r) // Replace entire weather data
	case http.MethodPatch:
		handleLocalPatch(w, r) // Update specific fields
	case http.MethodDelete:
		handleLocalDelete(w, r) // Delete weather.json
	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleLocalGet returns the locally saved weather.json data
// GET /weather/local
func handleLocalGet(w http.ResponseWriter, r *http.Request) {
	weather, err := loadFromJSON()
	if err != nil {
		http.Error(w, `{"error": "no saved data found, fetch weather first"}`, http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(weather, "", "  ")
	if err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

// handleLocalPost creates a new weather.json from the request body
// POST /weather/local  (send full WeatherResponse JSON in body)
// Returns 409 if weather.json already exists — use PUT to overwrite
func handleLocalPost(w http.ResponseWriter, r *http.Request) {
	// Check if weather.json already exists
	if _, err := os.Stat(weatherFile); err == nil {
		http.Error(w, `{"error": "weather.json already exists, use PUT to replace"}`, http.StatusConflict)
		return
	}

	// Decode the JSON body into our struct
	var weather WeatherResponse
	err := json.NewDecoder(r.Body).Decode(&weather)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "invalid JSON body: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Save to weather.json
	err = saveToJSON(&weather)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	w.Write([]byte(`{"message": "weather data created successfully"}`))
}

// handleLocalPut replaces the entire weather.json with new data from request body
// PUT /weather/local  (send full WeatherResponse JSON in body)
func handleLocalPut(w http.ResponseWriter, r *http.Request) {
	// Decode the full JSON body
	var weather WeatherResponse
	err := json.NewDecoder(r.Body).Decode(&weather)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "invalid JSON body: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Overwrite weather.json completely
	err = saveToJSON(&weather)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "weather data replaced successfully"}`))
}

// handleLocalPatch updates only the fields you send in the request body
// PATCH /weather/local  (send partial JSON — only fields you want to change)
//
// Example body to update just the temperature:
//
//	{"current": {"temp_c": 25.0, "temp_f": 77.0}}
func handleLocalPatch(w http.ResponseWriter, r *http.Request) {
	// First, load existing data
	existing, err := loadFromJSON()
	if err != nil {
		http.Error(w, `{"error": "no saved data to update, fetch or POST first"}`, http.StatusNotFound)
		return
	}

	// Marshal existing data to JSON bytes (so we can merge)
	existingBytes, _ := json.Marshal(existing)

	// Unmarshal into a map for flexible merging
	var existingMap map[string]interface{}
	json.Unmarshal(existingBytes, &existingMap)

	// Decode the partial update from request body into a map
	var patchMap map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&patchMap)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "invalid JSON body: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Merge: overwrite existing fields with patch values
	mergeMap(existingMap, patchMap)

	// Convert merged map back to our struct
	mergedBytes, _ := json.Marshal(existingMap)
	var updated WeatherResponse
	err = json.Unmarshal(mergedBytes, &updated)
	if err != nil {
		http.Error(w, `{"error": "failed to apply patch"}`, http.StatusInternalServerError)
		return
	}

	// Save updated data
	err = saveToJSON(&updated)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "weather data updated successfully"}`))
}

// handleLocalDelete deletes the weather.json file
// DELETE /weather/local
func handleLocalDelete(w http.ResponseWriter, r *http.Request) {
	// Check if the file exists first
	if _, err := os.Stat(weatherFile); os.IsNotExist(err) {
		http.Error(w, `{"error": "weather.json does not exist, nothing to delete"}`, http.StatusNotFound)
		return
	}

	// Delete the file
	err := os.Remove(weatherFile)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "failed to delete: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "weather.json deleted successfully"}`))
}

// mergeMap recursively merges src into dst (src values overwrite dst values)
// This lets PATCH update nested fields like {"current": {"temp_c": 25}}
func mergeMap(dst, src map[string]interface{}) {
	for key, srcVal := range src {
		// If both dst and src have a nested map for this key, merge recursively
		if dstVal, ok := dst[key]; ok {
			if dstMap, ok := dstVal.(map[string]interface{}); ok {
				if srcMap, ok := srcVal.(map[string]interface{}); ok {
					mergeMap(dstMap, srcMap)
					continue
				}
			}
		}
		// Otherwise, just overwrite
		dst[key] = srcVal
	}
}

// ============================================================
// Helper Functions
// ============================================================

// parseFlags converts ["--key=abc", "--q=London"] into a map {"key": "abc", "q": "London"}
func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for _, arg := range args {
		// Remove leading dashes
		for len(arg) > 0 && arg[0] == '-' {
			arg = arg[1:]
		}
		// Split on '=' to get key and value
		for i, ch := range arg {
			if ch == '=' {
				flags[arg[:i]] = arg[i+1:]
				break
			}
		}
	}
	return flags
}

// printWeatherSummary prints weather data in a readable format to the terminal
func printWeatherSummary(w *WeatherResponse) {
	fmt.Println()
	fmt.Println("=== Weather Summary ===")
	fmt.Printf("City       : %s, %s, %s\n", w.Location.Name, w.Location.Region, w.Location.Country)
	fmt.Printf("Location   : %.4f, %.4f\n", w.Location.Lat, w.Location.Lon)
	fmt.Printf("Timezone   : %s\n", w.Location.TzID)
	fmt.Printf("Local Time : %s\n", w.Location.Localtime)
	fmt.Printf("Updated    : %s\n", w.Current.LastUpdated)
	fmt.Printf("Temp       : %.1f°C / %.1f°F\n", w.Current.TempC, w.Current.TempF)
	fmt.Printf("Feels Like : %.1f°C / %.1f°F\n", w.Current.FeelsLikeC, w.Current.FeelsLikeF)
	fmt.Printf("Humidity   : %d%%\n", w.Current.Humidity)
	fmt.Printf("Pressure   : %.0f mb\n", w.Current.PressureMb)
	fmt.Printf("Wind       : %.1f kph %s\n", w.Current.WindKph, w.Current.WindDir)
	fmt.Printf("Cloud      : %d%%\n", w.Current.Cloud)
	fmt.Printf("Visibility : %.1f km\n", w.Current.VisKm)
	fmt.Printf("UV Index   : %.1f\n", w.Current.UV)
	fmt.Printf("Condition  : %s\n", w.Current.Condition.Text)
	fmt.Println()
}
