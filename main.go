package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	uptimeInSeconds = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "uptime_in_seconds_counter",
		Help:      "The uptime of this process in seconds",
		Subsystem: "demo",
		Namespace: "predictable",
	})

	currentMinute = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:      "current_minute_gauge",
		Help:      "The current minute of this hour",
		Subsystem: "demo",
		Namespace: "preditable",
	}, func() float64 { return float64(time.Now().Minute()) })

	secondsDistribution = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "seconds_distribution_histogram",
		Help:      "The distribution of seconds across runtime",
		Subsystem: "demo",
		Namespace: "preditable",
		Buckets:   prometheus.LinearBuckets(10, 10, 5),
	})

	summary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "seconds_distribution_summary",
		Help:       "The summary of seconds across runtime",
		Objectives: map[float64]float64{},
		Subsystem:  "demo",
		Namespace:  "preditable",
	})
)

func recordMetricUptimeInSeconds() {
	go func() {
		ticker := time.Tick(time.Second)
		for range ticker {
			uptimeInSeconds.Inc()
		}
	}()
}

func recordMetricSecondsDistribution() {
	go func() {
		ticker := time.Tick(time.Second)
		for now := range ticker {
			secondsDistribution.Observe(float64(now.Second()))
		}
	}()
}

func recordSummaryMetrics() {
	go func() {
		ticker := time.Tick(time.Second)
		for now := range ticker {
			summary.Observe(float64(now.Second()))
		}
	}()
}

func main() {
	recordMetricUptimeInSeconds()
	recordMetricSecondsDistribution()
	recordSummaryMetrics()

	registry := prometheus.NewRegistry()
	registry.MustRegister(uptimeInSeconds, currentMinute, secondsDistribution, summary)

	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", metricsHandler)

	log.Println("serving metrics at http://localhost:2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
