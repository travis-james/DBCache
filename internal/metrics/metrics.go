package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunPrometheusServer() {
	http.Handle("/metrics", promhttp.Handler())

	addr := ":2112"
	log.Printf("prometheus server listening at %v", addr)
	log.Fatal(http.ListenAndServe(":2112", nil))
}
