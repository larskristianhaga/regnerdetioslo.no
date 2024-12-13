package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var yrEndpoint = "https://www.yr.no/api/v0/locations/1-72837/forecast/currenthour"
var domain = "https://www.regnerdetioslo.no"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("App live and listening on port:", port)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		yrResponse := getYrData()

		isRaining := "Nei"
		if yrResponse.Precipitation.Value > 0 {
			isRaining = "Ja"
		}

		data := map[string]string{
			"isRaining":    isRaining,
			"dataFromTime": yrResponse.Created,
		}

		_ = t.ExecuteTemplate(w, "index.html.tmpl", data)
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		data := map[string]string{
			"URL": domain + "/sitemap.xml",
		}

		_ = t.ExecuteTemplate(w, "robots.txt.tmpl", data)
	})

	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")

		data := map[string]string{
			"URL": domain,
		}

		_ = t.ExecuteTemplate(w, "sitemap.xml.tmpl", data)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getYrData() Yr {
	// Create a custom HTTP client with disabled SSL verification.
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	response, err := client.Get(yrEndpoint)
	if err != nil {
		log.Fatal(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var yr Yr
	err = json.Unmarshal(responseData, &yr)
	if err != nil {
		log.Fatal(err.Error())
	}

	return yr
}

type Yr struct {
	Created       string        `json:"created"`
	Precipitation Precipitation `json:"precipitation"`
}

type Precipitation struct {
	Value float32 `json:"value"`
}
