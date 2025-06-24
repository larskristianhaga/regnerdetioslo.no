package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
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
var domain = "https://regnerdetioslo.no"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("App live and listening on port:", port)

	http.HandleFunc("/", loggingMiddleware(RootHandler))
	http.HandleFunc("/health", loggingMiddleware(HealthHandler))
	http.HandleFunc("/robots.txt", loggingMiddleware(RobotsHandler))
	http.HandleFunc("/sitemap.xml", loggingMiddleware(SitemapHandler))
	http.HandleFunc("/links", loggingMiddleware(LinksHandler))
	http.HandleFunc("/.well-known/security.txt", loggingMiddleware(SecurityHandler))

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

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("I'm healthy"))
}

func RobotsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, `User-agent: *
Allow: /
Allow: /links`)
}

func SitemapHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>`+domain+`</loc>
    </url>
    <url>
        <loc>`+domain+`/links</loc>
    </url>
</urlset>`)
}

func LinksHandler(w http.ResponseWriter, _ *http.Request) {
	_ = t.ExecuteTemplate(w, "index.html.tmpl", nil)
}

func SecurityHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, `Contact: mailto:larskhaga@gmail.com
Expires: 2030-12-31T23:59:00.000Z
Canonical: https://regnerdetioslo.no/.well-known/security.txt
`)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		userAgent := r.Header.Get("User-Agent")
		event := r.Method + " " + r.URL.Path + " " + r.Proto

		logEntry := LogEntry{
			Message:   "Request incoming",
			IP:        ip,
			Event:     event,
			Status:    "-",
			UserAgent: userAgent,
		}

		jsonLog, _ := json.Marshal(logEntry)
		log.SetFlags(0)
		log.Println(string(jsonLog))

		next(w, r)
	}
}

type LogEntry struct {
	Message   string `json:"message"`
	IP        string `json:"ip"`
	Event     string `json:"event"`
	Status    string `json:"status"`
	UserAgent string `json:"user_agent"`
}
