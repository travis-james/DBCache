package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/travis-james/DBCache/internal/config"
	"github.com/travis-james/DBCache/internal/datastore"
)

// RunPrometheusServer starts prometheus server at port 2112,
// at /metrics. Currently port is hardwired.
func RunPrometheusServer() {
	http.Handle("/metrics", promhttp.Handler())

	config, err := config.Load()
	if err != nil {
		panic(err) // TODO: Handle this more gracefully.
	}
	log.Printf("prometheus server listening at %s", config.PrometheusServerPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.PrometheusServerPort), nil)
	if err != nil {
		log.Fatalf("RunPrometheusServer error: %v", err)
	}
}

// MetricsManager contains the prometheus metrics for cache
// misses and hits, number of items in cache, and a ticker and
// cancel func to poll how many items are in the cache.
type MetricsManager struct {
	CacheMisses    prometheus.Counter
	CacheHits      prometheus.Counter
	CacheItemGauge prometheus.Gauge
	ticker         *time.Ticker
	cancelFunc     context.CancelFunc
}

// MetricsManagerInit sets up prometheus metrics for cache
// hit/miss/number of items. Takes a cache as an input to
// be polled for the size of the given cache.
func MetricsManagerInit(cache datastore.Cache) *MetricsManager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &MetricsManager{
		CacheMisses: promauto.NewCounter(prometheus.CounterOpts{
			Name: "dbcache_cache_total_misses",
			Help: "Total number of cache misses",
		}),
		CacheHits: promauto.NewCounter(prometheus.CounterOpts{
			Name: "dbcache_cache_total_hits",
			Help: "Total number of cache hits",
		}),
		CacheItemGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "dbcache_cache_items_total",
			Help: "Current number of items in the cache",
		}),
		ticker:     time.NewTicker(30 * time.Second),
		cancelFunc: cancel,
	}
	go m.pollCacheSize(ctx, cache)
	return m
}

// Close signals the cancel func, and stops the ticker.
func (m *MetricsManager) Close() {
	m.cancelFunc()
	m.ticker.Stop()
}

// pollCacheSize is an endless loop to check for a signal on the
// ticker, or a signal from the cancel func with the associated context.
func (m *MetricsManager) pollCacheSize(ctx context.Context, cache datastore.Cache) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.ticker.C: // Get number of items every 30 seconds.
			count, err := cache.NumberOfItems(ctx)
			if err == nil {
				m.CacheItemGauge.Set(float64(count))
			}
		}
	}
}
