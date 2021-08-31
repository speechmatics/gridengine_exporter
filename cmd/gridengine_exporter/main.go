package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/speechmatics/gridengine_exporter/collector"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var filter = flag.String("host-filter", "*", "Host to filter for, if not wildcard, only metrics for specified host is collected.")

func main() {
	flag.Parse()

	gridengineCollector := collector.NewGridengineCollector()
	gridengineCollector.Filter = *filter
	prometheus.MustRegister(gridengineCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
