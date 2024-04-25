package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricPrefix prefix for all metrics
const MetricPrefix = "succint"

// SuccintProject project name at succint
const SuccintProject = "blasrodri"

type SuccintCollector struct {
	metrics map[string]*prometheus.Desc
	logger  log.Logger
	mu      *sync.Mutex
}

// create metrics array
func addMetrics() map[string]*prometheus.Desc {

	metrics := make(map[string]*prometheus.Desc)

	metrics["proofs"] = prometheus.NewDesc(
		prometheus.BuildFQName(MetricPrefix, "proofs", "total"),
		"Total number of all proofs",
		[]string{"project", "status"}, nil,
	)
	metrics["timestamp"] = prometheus.NewDesc(
		prometheus.BuildFQName(MetricPrefix, "proof", "timestamp"),
		"Timestamp of latest proof",
		[]string{"project", "status"}, nil,
	)
	return metrics
}

func newSuccintCollector(logger log.Logger) *SuccintCollector {
	c := &SuccintCollector{
		logger:  logger,
		mu:      &sync.Mutex{},
		metrics: addMetrics(),
	}
	return c
}

// Collect implements prometheus.Collector.
func (c *SuccintCollector) Collect(ch chan<- prometheus.Metric) {

	c.mu.Lock()
	defer c.mu.Unlock()

	level.Debug(c.logger).Log("msg", "Started metrics collection")

	succintUrl := fmt.Sprintf("https://alpha.succinct.xyz/api/proofs?project=%s&limit=0&offset=0&status=all", *succintProject)
	level.Debug(c.logger).Log("msg", "Calling url: ", "succintUrl", succintUrl)
	response, err := http.Get(succintUrl)
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to make the API request", "err", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to read response body", "err", err)
		return
	}

	var proofs []SuccintProof
	if err := json.Unmarshal(body, &proofs); err != nil {
		level.Error(c.logger).Log("msg", "Error parsing JSON", "err", err)
		return
	}

	c.processMetrics(proofs, ch)

	level.Debug(c.logger).Log("msg", "Stopped metrics collection")

}

// Describe implements prometheus.Collector.
func (c *SuccintCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}

}

// processMetrics - processes the response data and sets the metrics
func (c *SuccintCollector) processMetrics(data []SuccintProof, ch chan<- prometheus.Metric) error {

	// total no of proofs
	ch <- prometheus.MustNewConstMetric(c.metrics["proofs"], prometheus.GaugeValue, float64(len(data)), *succintProject, "ALL")

	// total no of failed proofs
	ch <- prometheus.MustNewConstMetric(c.metrics["proofs"], prometheus.GaugeValue, float64(countFailedProofs(data)), *succintProject, "FAILED")

	// total no of running proofs
	ch <- prometheus.MustNewConstMetric(c.metrics["proofs"], prometheus.GaugeValue, float64(countRunningProofs(data)), *succintProject, "RUNNING")

	// timestamp of latest SUCCESS
	ch <- prometheus.MustNewConstMetric(c.metrics["timestamp"], prometheus.GaugeValue, float64(getLatestSuccessTimestamp(data).Unix()), *succintProject, "SUCCESS")

	// timestamp of latest FAILED
	ch <- prometheus.MustNewConstMetric(c.metrics["timestamp"], prometheus.GaugeValue, float64(getLatestFailureTimestamp(data).Unix()), *succintProject, "FAILED")

	return nil
}
