package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/paulcager/osgridref"
	flag "github.com/spf13/pflag"
)

const (
	apiVersion = "v4"
)

var (
	staticCache time.Duration
	listenPort  string
)

func main() {
	flag.StringVar(&listenPort, "port", ":9090", "Port to listen on")
	flag.DurationVar(&staticCache, "static-cache-max-age", 1*time.Hour, "If not zero, the max-age property to set in Cache-Control for responses")
	flag.Parse()

	server := makeHTTPServer(listenPort)
	log.Fatal(server.ListenAndServe())
}

type Reply struct {
	OSGridRef string  `json:"osGridRef"`
	Easting   int     `json:"easting"`
	Northing  int     `json:"northing"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
}

func makeHTTPServer(listenPort string) *http.Server {
	http.Handle("/"+apiVersion+"/gridref/", makeCachingHandler(1*time.Hour, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			gridRefStr := r.URL.Path[len("/"+apiVersion+"/gridref/"):]
			gridRef, err := osgridref.ParseOsGridRef(gridRefStr)
			if err != nil {
				handleError(w, r, gridRefStr, err)
				return
			}

			lat, lon := gridRef.ToLatLon()
			handle(w, r, gridRef, lat, lon)
		})))

	http.Handle("/"+apiVersion+"/latlon/", makeCachingHandler(1*time.Hour, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			latLonStr := r.URL.Path[len("/"+apiVersion+"/latlon/"):]
			latLon, err := osgridref.ParseLatLon(latLonStr, 0, osgridref.WGS84)
			if err != nil {
				handleError(w, r, latLonStr, err)
				return
			}

			gridRef := latLon.ToOsGridRef()
			handle(w, r, gridRef, latLon.Lat, latLon.Lon)
		})))

	if !strings.Contains(listenPort, ":") {
		listenPort = ":" + listenPort
	}

	log.Println("Starting HTTP server on " + listenPort)
	s := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       10 * time.Minute,
		Handler:           makeLoggingHandler(http.DefaultServeMux),
		Addr:              listenPort,
	}

	return s
}

func handleError(w http.ResponseWriter, _ *http.Request, str string, _ error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Invalid request: %q\n", str)
}

func handle(w http.ResponseWriter, _ *http.Request, ref osgridref.OsGridRef, lat float64, lon float64) {
	reply := Reply{
		OSGridRef: ref.StringNCompact(8),
		Easting:   ref.Easting,
		Northing:  ref.Northing,
		Lat:       lat,
		Lon:       lon,
	}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	err := enc.Encode(reply)
	if err != nil {
		log.Printf("Failed to write response: %s", err)
		w.WriteHeader(http.StatusBadGateway)
	}
}

func makeCachingHandler(age time.Duration, h http.Handler) http.Handler {
	ageSeconds := int64(math.Round(age.Seconds()))
	if ageSeconds <= 0 {
		return h
	}

	header := fmt.Sprintf("public,max-age=%d", ageSeconds)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", header)
			h.ServeHTTP(w, r)
		})
}

func makeLoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			end := time.Now()

			uri := r.URL.String()
			method := r.Method
			fmt.Printf("%s %s %s %d\n", method, uri, r.RemoteAddr, end.Sub(start).Milliseconds())
		})
}
