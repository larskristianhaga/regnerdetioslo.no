package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var endpoint = "https://www.yr.no/api/v0/locations/1-72837/forecast/currenthour"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("👋"))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		// Create a custom HTTP client with disabled SSL verification.
		customTransport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: customTransport}

		// Get data from the API using the custom client.
		response, err := client.Get(endpoint)

		// Guard for empty response.
		if err != nil {
			log.Fatal(err.Error())
		}

		// Get the JSON data out of the response.
		responseData, err := io.ReadAll(response.Body)

		// Guard for something wrong with reading the content.
		if err != nil {
			log.Fatal(err.Error())
		}

		var yr Yr
		// JSONIFY the string.
		_ = json.Unmarshal(responseData, &yr)

		// Extract the values I want.
		doesItRain := yr.Precipitation.Value > 0
		lastUpdatedAt := yr.Created

		// Enable CORS.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Set the content type to JSON.
		w.Header().Set("Content-Type", "application/json")

		// Return the object.
		_ = json.NewEncoder(w).Encode(struct {
			DoesItRain   bool
			DataFromTime string
		}{DoesItRain: doesItRain, DataFromTime: lastUpdatedAt})

	})

	log.Println("App live and listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Yr struct {
	Created       string        `json:"created"`
	Precipitation Precipitation `json:"precipitation"`
}

type Precipitation struct {
	Value float32 `json:"value"`
}
