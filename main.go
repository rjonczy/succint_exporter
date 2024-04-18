package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	registry := prometheus.NewRegistry()

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "my_metric",
		Help: "This is a simple gauge metric",
	})

	registry.MustRegister(gauge)

	gauge.Set(1)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Println("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
