package main

// ============================================================
// models.go - Structs that map to WeatherAPI.com JSON response
// We use `json:"..."` tags so Go knows how to marshal/unmarshal
//
// API: http://api.weatherapi.com/v1/current.json?key=KEY&q=QUERY
// Docs: https://www.weatherapi.com/docs/
// ============================================================

// WeatherResponse is the top-level response from WeatherAPI.com
type WeatherResponse struct {
	Location Location       `json:"location"` // Location info (city, country, coords)
	Current  CurrentWeather `json:"current"`  // Current weather data
}

// Location holds info about the matched location
type Location struct {
	Name           string  `json:"name"`            // City/town name
	Region         string  `json:"region"`          // Region or state
	Country        string  `json:"country"`         // Country name
	Lat            float64 `json:"lat"`             // Latitude
	Lon            float64 `json:"lon"`             // Longitude
	TzID           string  `json:"tz_id"`           // Timezone ID (e.g. "Asia/Kolkata")
	LocaltimeEpoch int64   `json:"localtime_epoch"` // Local time as unix timestamp
	Localtime      string  `json:"localtime"`       // Local time as string
}

// CurrentWeather holds the current weather conditions
type CurrentWeather struct {
	LastUpdated string    `json:"last_updated"` // Last updated time string
	TempC       float64   `json:"temp_c"`       // Temperature in Celsius
	TempF       float64   `json:"temp_f"`       // Temperature in Fahrenheit
	IsDay       int       `json:"is_day"`       // 1 = day, 0 = night
	Condition   Condition `json:"condition"`    // Weather condition (text + icon)
	WindMph     float64   `json:"wind_mph"`     // Wind speed in mph
	WindKph     float64   `json:"wind_kph"`     // Wind speed in kph
	WindDegree  int       `json:"wind_degree"`  // Wind direction in degrees
	WindDir     string    `json:"wind_dir"`     // Wind direction (e.g. "NW")
	PressureMb  float64   `json:"pressure_mb"`  // Pressure in millibars
	PrecipMm    float64   `json:"precip_mm"`    // Precipitation in mm
	Humidity    int       `json:"humidity"`     // Humidity percentage
	Cloud       int       `json:"cloud"`        // Cloud cover percentage
	FeelsLikeC  float64   `json:"feelslike_c"`  // Feels like in Celsius
	FeelsLikeF  float64   `json:"feelslike_f"`  // Feels like in Fahrenheit
	VisKm       float64   `json:"vis_km"`       // Visibility in km
	UV          float64   `json:"uv"`           // UV index
	GustMph     float64   `json:"gust_mph"`     // Wind gust in mph
	GustKph     float64   `json:"gust_kph"`     // Wind gust in kph
}

// Condition describes the weather (e.g. "Partly cloudy")
type Condition struct {
	Text string `json:"text"` // Description (e.g. "Sunny", "Partly cloudy")
	Icon string `json:"icon"` // Icon URL
	Code int    `json:"code"` // Condition code
}
