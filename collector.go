package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type SuccintCollector struct {
	metrics map[string]*prometheus.Desc
	logger  log.Logger
}

func newSuccintCollector(logger log.Logger) *SuccintCollector {
	c := &SuccintCollector{
		logger: logger,
	}
	return c
}

// Collect implements prometheus.Collector.
func (c SuccintCollector) Collect(ch chan<- prometheus.Metric) {

	response, err := http.Get(*succintUrl)
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to make the API request", "err", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to read response body", "err", err)
	}

	var proofs []SuccintProof
	if err := json.Unmarshal(body, &proofs); err != nil {
		level.Error(c.logger).Log("msg", "Error parsing JSON", "err", err)
	}

	for _, record := range proofs {
		fmt.Printf("ID: %s, Status: %s, CreatedAt: %s\n", record.ID, record.Status, record.CreatedAt)
	}

	fmt.Println("----")
	fmt.Println(len(proofs))
	fmt.Println("----")

	// data := []*Datum{}
	// var err error
	// scrape
	// if len(e.TargetURLs()) > 0 {
	// 	data, err = e.gatherData()
	// 	if err != nil {
	// 		log.Errorf("Error gathering Data from remote API: %v", err)
	// 		return
	// 	}
	// }

	// rates, err := e.getRates()
	// if err != nil {
	// 	log.Errorf("Error gathering Rates from remote API: %v", err)
	// 	return
	// }

	// Set prometheus gauge metrics using the data gathered
	// err = e.processMetrics(data, rates, ch)

	// if err != nil {
	// 	log.Error("Error Processing Metrics", err)
	// 	return
	// }

	// log.Info("All Metrics successfully collected")

}

// Describe implements prometheus.Collector.
func (c SuccintCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}

}
