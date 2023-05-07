package main

import (
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name:        "http_requests_total",
		Help:        "Total number of HTTP requests",
		ConstLabels: prometheus.Labels{"server": "api"},
	},
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.Inc()
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func main() {
	prometheus.MustRegister(httpRequestsTotal)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.Inc()
		w.Write([]byte("Hello, world!"))
	})

	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.HandleFunc("/healthz", HealthCheck)

	http.ListenAndServe(":8080", nil)
}
