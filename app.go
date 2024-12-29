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

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/ping", PingHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/robots.txt", RobotsHandler)
	http.HandleFunc("/sitemap.xml", SitemapHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func RootHandler(w http.ResponseWriter, _ *http.Request) {
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
}

func PingHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("I'm healthy"))
}

func RobotsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	data := map[string]string{
		"URL": domain + "/sitemap.xml",
	}

	_ = t.ExecuteTemplate(w, "robots.txt.tmpl", data)
}

func SitemapHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	data := map[string]string{
		"URL": domain,
	}

	_ = t.ExecuteTemplate(w, "sitemap.xml.tmpl", data)
}

func getYrData() Yr {
	client := createInsecureHTTPClient()

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

// Create a custom HTTP client with disabled SSL verification.
func createInsecureHTTPClient() *http.Client {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: customTransport}
}

type Yr struct {
	Created       string        `json:"created"`
	Precipitation Precipitation `json:"precipitation"`
}

type Precipitation struct {
	Value float32 `json:"value"`
}
