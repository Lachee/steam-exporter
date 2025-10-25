package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	steamMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "steam_exporter_metric",
			Help: "A sample metric for Steam exporter",
		},
		[]string{"label"},
	)
)

func init() {
	prometheus.MustRegister(steamMetric)
}

func main() {
	// Set a sample metric value
	steamMetric.WithLabelValues("example").Set(1)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting Steam exporter on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
