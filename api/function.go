package regnerdetioslo

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("RegnerDetIOsloFunction", RegnerDetIOsloFunction)
}

func RegnerDetIOsloFunction(w http.ResponseWriter, r *http.Request) {

	// API endpoint.
	endpoint := 
"https://www.yr.no/api/v0/locations/1-72837/forecast/currenthour"

	// Get data from the API.
	response, err := http.Get(endpoint)

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

	// Return the object.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		DoesItRain   bool
		DataFromTime string
	}{DoesItRain: doesItRain, DataFromTime: lastUpdatedAt})
}

type Yr struct {
	Created       string        `json:"created"`
	Precipitation Precipitation `json:"precipitation"`
}

type Precipitation struct {
	Value float32 `json:"value"`
}

