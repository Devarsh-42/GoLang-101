# Day 11 — Weather CLI & REST API (Go)

> A complete Go project that integrates with [WeatherAPI.com](https://www.weatherapi.com/) to fetch real-time weather data, serve it over a REST API with full **CRUD** operations, persist it locally as JSON, and export it to CSV — all built with **zero external dependencies**.

---

## Table of Contents

- [Overview](#overview)
- [Topics Covered](#topics-covered)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
- [CLI Commands](#cli-commands)
  - [Start the HTTP Server](#1-start-the-http-server)
  - [Fetch Weather Data](#2-fetch-weather-data)
  - [Show Saved Data](#3-show-saved-data)
  - [Export to CSV](#4-export-to-csv)
- [HTTP API Endpoints](#http-api-endpoints)
  - [Postman Testing Guide](#postman-testing-guide)
  - [Example Request & Response](#example-request--response)
- [CRUD Operations](#crud-operations)
  - [POST — Create](#post--create)
  - [GET — Read](#get--read)
  - [PUT — Replace](#put--replace)
  - [PATCH — Partial Update](#patch--partial-update)
  - [DELETE — Remove](#delete--remove)
- [Running Tests](#running-tests)
- [File-by-File Breakdown](#file-by-file-breakdown)

---

## Overview

This project demonstrates building a **production-style Go application** from scratch using only the standard library. It covers the full lifecycle of working with external APIs — fetching data, processing it, storing it locally, serving it over HTTP, and allowing CRUD manipulation of the stored data.

**What this project does:**

1. **Fetches** real-time weather data from WeatherAPI.com (free tier)
2. **Stores** the response locally in `weather.json` using JSON marshal/unmarshal
3. **Serves** the data over an HTTP server with middleware (logging + JSON headers)
4. **Supports full CRUD** — Create, Read, Update (partial), Replace, and Delete on local data
5. **Exports** weather data to a `weather.csv` file
6. **CLI interface** — All features are accessible via command-line commands
7. **16 unit tests** covering every handler, middleware, storage, and CRUD operation

---

## Topics Covered

| #   | Topic                                                                               | Where                         |
| --- | ----------------------------------------------------------------------------------- | ----------------------------- |
| 1   | **HTTP Server** — `net/http`, `http.NewServeMux`, `ListenAndServe`                  | `main.go`                     |
| 2   | **HTTP Handlers** — `http.HandlerFunc`, request/response pattern                    | `main.go`                     |
| 3   | **HTTP Methods** — GET, POST, PUT, PATCH, DELETE routing                            | `main.go`                     |
| 4   | **Middleware Pattern** — Function wrapping (`func(http.Handler) http.Handler`)      | `middleware.go`               |
| 5   | **Structured Logging** — `log/slog` for request logging                             | `middleware.go`               |
| 6   | **JSON Marshal / Unmarshal** — `encoding/json` (struct ↔ JSON)                      | `storage.go`, `main.go`       |
| 7   | **JSON Tags** — `json:"field_name"` struct tags for API mapping                     | `models.go`                   |
| 8   | **Structs** — Nested structs modelling real API responses                           | `models.go`                   |
| 9   | **HTTP Client** — Making external API calls with `http.Get`                         | `weather.go`                  |
| 10  | **File I/O** — `os.ReadFile`, `os.WriteFile`, `os.Create`, `os.Remove`              | `storage.go`, `csv_export.go` |
| 11  | **CSV Writing** — `encoding/csv` for tabular data export                            | `csv_export.go`               |
| 12  | **Error Handling** — `fmt.Errorf` with `%w` wrapping, status codes                  | Throughout                    |
| 13  | **CLI Argument Parsing** — Custom flag parser (no external libs)                    | `main.go`                     |
| 14  | **Unit Testing** — `testing` package, `httptest.NewRequest`, `httptest.NewRecorder` | `main_test.go`                |
| 15  | **Query Parameters** — `r.URL.Query().Get()` for reading URL params                 | `main.go`                     |
| 16  | **HTTP Status Codes** — 200, 201, 400, 404, 405, 409, 500                           | `main.go`                     |
| 17  | **CRUD REST Pattern** — Full Create/Read/Update/Delete on a single endpoint         | `main.go`                     |
| 18  | **Recursive Map Merging** — Deep merge for PATCH partial updates                    | `main.go`                     |
| 19  | **`defer`** — Deferred file/body closing                                            | `weather.go`, `csv_export.go` |
| 20  | **String Formatting** — `fmt.Sprintf`, `strconv.Itoa`                               | Throughout                    |

---

## Project Structure

```
Day-11/
├── go.mod              # Go module definition (day11, no external deps)
├── README.md           # This file
├── main.go             # Entry point: CLI commands + HTTP server + CRUD handlers
├── models.go           # Go structs mapping to WeatherAPI.com JSON response
├── weather.go          # HTTP client — fetches data from WeatherAPI.com
├── storage.go          # JSON file persistence (marshal/unmarshal to weather.json)
├── csv_export.go       # Export weather data to weather.csv
├── middleware.go        # HTTP middleware (logging + JSON content-type)
├── main_test.go        # 16 unit tests for all features
├── weather.json        # (Generated) Saved weather data
└── weather.csv         # (Generated) Exported CSV data
```

---

## How It Works

```
                        ┌─────────────────────┐
   CLI / HTTP Request   │   WeatherAPI.com     │
         │              │  (Free Current API)  │
         ▼              └──────────┬───────────┘
  ┌─────────────┐                  │ JSON Response
  │   main.go   │◄─────────────────┘
  │  (Router)   │
  └──────┬──────┘
         │
         ├──► weather.go     →  Fetches data from API
         ├──► models.go      →  Parses JSON into Go structs
         ├──► storage.go     →  Saves/loads weather.json (Marshal/Unmarshal)
         ├──► csv_export.go  →  Exports to weather.csv
         └──► middleware.go  →  Wraps handlers with logging & JSON headers
```

**Request Flow (HTTP Server):**

```
Client Request
      │
      ▼
loggingMiddleware          → Logs method, path, duration
      │
      ▼
jsonContentTypeMiddleware  → Sets Content-Type: application/json
      │
      ▼
Route Handler              → Processes the request
      │
      ▼
Response sent back
```

---

## Getting Started

### Prerequisites

- **Go 1.21+** installed ([download](https://go.dev/dl/))
- **WeatherAPI.com API key** (free) — [sign up here](https://www.weatherapi.com/signup.aspx)

### Setup

```bash
# Clone the repo and navigate to Day-11
cd GoLang-101/Day-11

# Verify Go is installed
go version

# Run tests to make sure everything works
go test ./... -v
```

No `go mod tidy` needed — there are **zero external dependencies**. Everything uses Go's standard library.

---

## CLI Commands

### 1. Start the HTTP Server

```bash
go run . serve
```

Starts the server on `http://localhost:8080` with all endpoints available.

### 2. Fetch Weather Data

```bash
# Fetch by city name
go run . fetch --key=YOUR_API_KEY --q=London

# Fetch by coordinates (lat,lon)
go run . fetch --key=YOUR_API_KEY --q=28.6139,77.2090

# Fetch by zip/postal code
go run . fetch --key=YOUR_API_KEY --q=10001

# Fetch by IP address
go run . fetch --key=YOUR_API_KEY --q=auto:ip

# Fetch with language (Hindi)
go run . fetch --key=YOUR_API_KEY --q=Delhi --lang=hi
```

This fetches weather from WeatherAPI.com and saves it to `weather.json`.

### 3. Show Saved Data

```bash
go run . show
```

Reads `weather.json` and prints a formatted summary:

```
=== Weather Summary ===
City       : London, City of London, United Kingdom
Location   : 51.5171, -0.1062
Timezone   : Europe/London
Local Time : 2025-02-09 12:30
Temp       : 8.0°C / 46.4°F
Feels Like : 5.2°C / 41.4°F
Humidity   : 71%
Wind       : 16.9 kph W
Condition  : Partly cloudy
```

### 4. Export to CSV

```bash
go run . csv
```

Reads `weather.json` and creates `weather.csv` with 17 columns:

| City | Region | Country | Latitude | Longitude | Localtime | TempC | TempF | FeelsLikeC | Humidity | PressureMb | WindKph | WindDir | Cloud | VisKm | UV  | Condition |
| ---- | ------ | ------- | -------- | --------- | --------- | ----- | ----- | ---------- | -------- | ---------- | ------- | ------- | ----- | ----- | --- | --------- |

---

## HTTP API Endpoints

Start the server first with `go run . serve`, then use any HTTP client (browser, curl, Postman).

| Method   | Endpoint               | Description                                    |
| -------- | ---------------------- | ---------------------------------------------- |
| `GET`    | `/`                    | Welcome message                                |
| `GET`    | `/goodbye`             | Goodbye message                                |
| `GET`    | `/weather?key=..&q=..` | Fetch weather from WeatherAPI.com              |
| `GET`    | `/weather/local`       | Read saved `weather.json`                      |
| `POST`   | `/weather/local`       | Create new weather data (fails if file exists) |
| `PUT`    | `/weather/local`       | Replace entire `weather.json`                  |
| `PATCH`  | `/weather/local`       | Update specific fields only                    |
| `DELETE` | `/weather/local`       | Delete `weather.json` file                     |

### Postman Testing Guide

#### 1. Welcome Message

```
GET http://localhost:8080/
```

#### 2. Goodbye Message

```
GET http://localhost:8080/goodbye
```

#### 3. Fetch Weather from API

```
GET http://localhost:8080/weather?key=YOUR_API_KEY&q=London
```

Optional `lang` param:

```
GET http://localhost:8080/weather?key=YOUR_API_KEY&q=Delhi&lang=hi
```

Query (`q`) accepts: city name, `lat,lon`, zip code, or IP address.

#### 4. Read Local Data

```
GET http://localhost:8080/weather/local
```

Returns the contents of `weather.json` (404 if file doesn't exist).

#### 5. Create Local Data (POST)

```
POST http://localhost:8080/weather/local
Content-Type: application/json

{
  "location": {
    "name": "Mumbai",
    "region": "Maharashtra",
    "country": "India",
    "lat": 19.0760,
    "lon": 72.8777
  },
  "current": {
    "temp_c": 32.0,
    "temp_f": 89.6,
    "humidity": 70,
    "condition": {
      "text": "Partly cloudy"
    }
  }
}
```

Returns `201 Created`. Returns `409 Conflict` if `weather.json` already exists.

#### 6. Replace All Data (PUT)

```
PUT http://localhost:8080/weather/local
Content-Type: application/json

{
  "location": {
    "name": "Paris",
    "country": "France",
    "lat": 48.8566,
    "lon": 2.3522
  },
  "current": {
    "temp_c": 18.0,
    "temp_f": 64.4,
    "humidity": 55,
    "condition": {
      "text": "Clear"
    }
  }
}
```

Overwrites the entire file regardless of existing content.

#### 7. Partial Update (PATCH)

```
PATCH http://localhost:8080/weather/local
Content-Type: application/json

{
  "current": {
    "temp_c": 25.0,
    "temp_f": 77.0
  }
}
```

Only updates the fields you send — everything else stays the same. Uses recursive deep merge, so nested objects are handled correctly.

#### 8. Delete Local Data

```
DELETE http://localhost:8080/weather/local
```

Deletes the `weather.json` file. Returns `404` if file doesn't exist.

### Example Request & Response

**Request:**

```
GET http://localhost:8080/weather?key=YOUR_KEY&q=London
```

**Response (200 OK):**

```json
{
  "location": {
    "name": "London",
    "region": "City of London, Greater London",
    "country": "United Kingdom",
    "lat": 51.52,
    "lon": -0.11,
    "tz_id": "Europe/London",
    "localtime_epoch": 1739098200,
    "localtime": "2025-02-09 12:30"
  },
  "current": {
    "last_updated": "2025-02-09 12:15",
    "temp_c": 8.0,
    "temp_f": 46.4,
    "is_day": 1,
    "condition": {
      "text": "Partly cloudy",
      "icon": "//cdn.weatherapi.com/weather/64x64/day/116.png",
      "code": 1003
    },
    "wind_mph": 10.5,
    "wind_kph": 16.9,
    "wind_degree": 270,
    "wind_dir": "W",
    "pressure_mb": 1020.0,
    "precip_mm": 0.0,
    "humidity": 71,
    "cloud": 50,
    "feelslike_c": 5.2,
    "feelslike_f": 41.4,
    "vis_km": 10.0,
    "uv": 2.0,
    "gust_mph": 14.3,
    "gust_kph": 23.0
  }
}
```

---

## CRUD Operations

All CRUD operations work on the `/weather/local` endpoint using different HTTP methods:

### POST — Create

- Creates a new `weather.json` with the data from the request body
- **Fails with `409 Conflict`** if the file already exists (use PUT to replace)
- Returns `201 Created` on success

### GET — Read

- Returns the contents of `weather.json` as formatted JSON
- Returns `404 Not Found` if no data has been saved yet

### PUT — Replace

- Completely replaces the contents of `weather.json`
- Always overwrites — idempotent operation
- Returns `200 OK` on success

### PATCH — Partial Update

- Updates **only the fields you send** in the request body
- Uses recursive deep merge — nested objects (like `current.temp_c`) merge correctly without losing sibling fields
- Returns `404` if no existing data to patch
- Returns `200 OK` on success

### DELETE — Remove

- Deletes the `weather.json` file from disk
- Returns `404 Not Found` if file doesn't exist
- Returns `200 OK` on success

---

## Running Tests

```bash
# Run all 16 tests with verbose output
go test ./... -v
```

**Tests included (16 total):**

| #   | Test                                 | What It Covers                                  |
| --- | ------------------------------------ | ----------------------------------------------- |
| 1   | `TestHandleHello`                    | `GET /` returns welcome message                 |
| 2   | `TestHandleGoodBye`                  | `GET /goodbye` returns goodbye message          |
| 3   | `TestHandleWeatherAPI_MissingParams` | `/weather` without params returns 400           |
| 4   | `TestSaveAndLoadJSON`                | Marshal → save → load → unmarshal round trip    |
| 5   | `TestParseFlags`                     | CLI flag parser (`--key=val` → map)             |
| 6   | `TestHandleSavedData_NoFile`         | `GET /weather/local` returns 404 when no file   |
| 7   | `TestLocalPost`                      | `POST /weather/local` creates file, returns 201 |
| 8   | `TestLocalPost_Conflict`             | `POST` returns 409 when file already exists     |
| 9   | `TestLocalPut`                       | `PUT /weather/local` replaces file content      |
| 10  | `TestLocalPatch`                     | `PATCH` updates only specified fields           |
| 11  | `TestLocalDelete`                    | `DELETE` removes the file                       |
| 12  | `TestLocalDelete_NotFound`           | `DELETE` returns 404 when no file exists        |
| 13  | `TestLoggingMiddleware`              | Logging middleware passes requests through      |
| 14  | `TestJSONContentTypeMiddleware`      | Sets `Content-Type: application/json` header    |
| 15  | `TestExportToCSV`                    | CSV export creates valid file                   |
| 16  | `TestMarshalUnmarshal`               | Direct `json.Marshal` / `json.Unmarshal` verify |

---

## File-by-File Breakdown

### `models.go`

Defines **Go structs** that map to the WeatherAPI.com JSON response using `json:"..."` struct tags. Three nested structs:

- `WeatherResponse` — Top-level (contains Location + Current)
- `Location` — City name, region, country, coordinates, timezone, local time
- `CurrentWeather` — Temperature, humidity, wind, pressure, UV, cloud cover, visibility, condition
- `Condition` — Weather description text, icon URL, condition code

### `weather.go`

The **HTTP client** that calls the WeatherAPI.com free endpoint (`/v1/current.json`). Builds the URL with query params, makes a `GET` request, reads the body, checks the status code, and unmarshals the JSON into our Go structs.

### `storage.go`

Handles **local persistence** using `json.MarshalIndent` (struct → pretty JSON → file) and `json.Unmarshal` (file → bytes → struct). Reads and writes `weather.json`.

### `csv_export.go`

Loads data from `weather.json` and writes it to `weather.csv` using Go's `encoding/csv` package. Exports 17 columns of weather data in a single row with headers.

### `middleware.go`

Two **HTTP middleware** functions using the wrapper pattern:

- `loggingMiddleware` — Logs every request (method, path, duration) using `slog`
- `jsonContentTypeMiddleware` — Sets `Content-Type: application/json` on all responses

### `main.go`

The **entry point** with two modes:

1. **CLI mode** — Parses `os.Args` to dispatch commands (`serve`, `fetch`, `show`, `csv`)
2. **Server mode** — Starts an HTTP server with routes, middleware, and full CRUD handlers for `/weather/local`

Includes a custom `parseFlags()` function and a recursive `mergeMap()` helper for PATCH operations.

### `main_test.go`

**16 unit tests** using `testing`, `net/http/httptest`, and temp file cleanup. Tests every handler, middleware, storage function, CLI parser, CSV export, and all CRUD operations.

---

## Tech Stack

| Component             | Technology                                                |
| --------------------- | --------------------------------------------------------- |
| Language              | Go 1.25                                                   |
| API                   | [WeatherAPI.com](https://www.weatherapi.com/) (Free Plan) |
| HTTP Server           | `net/http` (standard library)                             |
| JSON                  | `encoding/json`                                           |
| CSV                   | `encoding/csv`                                            |
| Logging               | `log/slog`                                                |
| Testing               | `testing` + `net/http/httptest`                           |
| External Dependencies | **None** — 100% standard library                          |

---

_Built as part of the [GoLang-101](https://github.com/) learning series._
