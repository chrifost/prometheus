package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var count int

var (
	taskCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "http_handler",
		Name:      "completed_request_total",
		Help:      "Total number of request recieved.",
	})
	gaugeValue = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gauge_value",
		Help: "Gauge arugment provided to /counter",
	})
	summaryValue = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "temperature_celsius",
		Help:       "Temperatures from a system",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func counterHandler(w http.ResponseWriter, req *http.Request) {
	var gauge string
	fmt.Printf("Executing counterHandler.\n")

	// All we care about is the query param and parsing it
	m, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		fmt.Println("Error parsing req.URL.RawQuery")
	}

	for key, value := range m {
		switch key {
		case "gauge":
			gauge = value[0]
			if gauge == "" {
				fmt.Println("Error parsing gauge: ", err)
			}
		}
	}

	fmt.Println("gauge: ", gauge)

	f, err := strconv.ParseFloat(gauge, 64)
	if err != nil {
		fmt.Println("Error ParseFloat")
	}
	fmt.Println("gauge: ", f)
	gaugeValue.Set(f)

	count++
	taskCounter.Inc()

	fmt.Fprintf(w, "Executing counterHandler %v", count)
}

func main() {

	fmt.Printf("Starting p8s test service...\n")

	http.HandleFunc("/counter", counterHandler)
	//http.HandleFunc("/metrics", prometheusHandler)
	http.Handle("/metrics", promhttp.Handler())

	// Create my counter metric
	if err := prometheus.Register(taskCounter); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("taskCounter registered.")
	}

	// Create my counter metric
	if err := prometheus.Register(gaugeValue); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("gaugeValue registered.")
		gaugeValue.Set(0.0)
	}

	// Create my counter metric
	if err := prometheus.Register(summaryValue); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("summaryValue registered.")

		// Simulate some Temperatures
		for i := 0; i < 1000; i++ {
			summaryValue.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
		}
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
