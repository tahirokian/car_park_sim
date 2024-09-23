package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	numberOfProcessedEvents  prometheus.CounterVec
	startupDelay             prometheus.Gauge
	eventProcessingDelay     prometheus.GaugeVec
	eventProcessingDelayHist prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		numberOfProcessedEvents: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "go_backend",
				Name:      "total_processed_events",
				Help:      "Total number of processed vehicle entry events for go backend service.",
			},
			[]string{"event_type"},
		),
		startupDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "go_backend",
				Name:      "startup_delay",
				Help:      "Startup delay for go backend service.",
			},
		),
		eventProcessingDelay: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "go_backend",
				Name:      "event_processing_delay",
				Help:      "Event processing delay for go backend service.",
			},
			[]string{"event_type"},
		),
		eventProcessingDelayHist: *prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "go_backend",
				Name:      "event_processing_delay_hist",
				Help:      "Event processing delay histogram for go backend service.",
			},
			[]string{"event_type"},
		),
	}

	reg.MustRegister(m.numberOfProcessedEvents)
	reg.MustRegister(m.startupDelay)
	reg.MustRegister(m.eventProcessingDelay)
	reg.MustRegister(m.eventProcessingDelayHist)
	return m
}

func setupPromethusEndpoint(reg *prometheus.Registry) {
	goBackendAddress := os.Getenv("GO_BACKEND_ADDR")
	goBackendPort := os.Getenv("GO_BACKEND_PORT")
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	http.Handle("/metrics", promHandler)
	log.Printf("Prometheus metrics available at http://%s:%s/metrics\n",
		goBackendAddress, goBackendPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", goBackendPort), nil))
}
