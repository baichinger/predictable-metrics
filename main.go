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
		Name: "uptime_in_seconds",
		Help: "The uptime of this process in seconds",
	})

	currentMinute = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "current_minute",
		Help: "The current minute of this hour",
	}, func() float64 { return float64(time.Now().Minute()) })

	secondsDistribution = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "seconds_distribution",
		Help:    "The distribution of seconds across runtime",
		Buckets: prometheus.LinearBuckets(10, 10, 5),
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

func main() {
	recordMetricUptimeInSeconds()
	recordMetricSecondsDistribution()

	registry := prometheus.NewRegistry()
	registry.MustRegister(uptimeInSeconds, currentMinute, secondsDistribution)

	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", metricsHandler)

	log.Fatal(http.ListenAndServe(":2112", nil))
}
