package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	uptimeInSecondsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "predictable",
		Subsystem: "demo",
		Name:      "uptime_in_seconds_counter",
		Help:      "The uptime of this process in seconds",
	})

	currentMinuteGauge = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "predictable",
		Subsystem: "demo",
		Name:      "current_minute_gauge",
		Help:      "The current minute of this hour",
	}, func() float64 { return float64(time.Now().Minute()) })

	secondsDistributionHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "predictable",
		Subsystem: "demo",
		Name:      "seconds_distribution_histogram",
		Help:      "The distribution of seconds across runtime",
		Buckets:   prometheus.LinearBuckets(10, 10, 5),
	})

	secondsDistributionSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:  "predictable",
		Subsystem:  "demo",
		Name:       "seconds_distribution_summary",
		Help:       "The distribution of seconds across runtime",
		Objectives: map[float64]float64{},
	})
)

func recordMetricCounter() {
	go func() {
		ticker := time.Tick(time.Second)
		for range ticker {
			uptimeInSecondsCounter.Inc()
		}
	}()
}

func recordMetricHistogram() {
	go func() {
		ticker := time.Tick(time.Second)
		for now := range ticker {
			secondsDistributionHistogram.Observe(float64(now.Second()))
		}
	}()
}

func recordMetricSummary() {
	go func() {
		ticker := time.Tick(time.Second)
		for now := range ticker {
			secondsDistributionSummary.Observe(float64(now.Second()))
		}
	}()
}

func main() {
	recordMetricCounter()
	recordMetricHistogram()
	recordMetricSummary()

	registry := prometheus.NewRegistry()
	registry.MustRegister(uptimeInSecondsCounter, currentMinuteGauge, secondsDistributionHistogram, secondsDistributionSummary)

	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", metricsHandler)

	log.Println("serving metrics at http://localhost:2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
