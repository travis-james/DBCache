package metrics

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/travis-james/DBCache/internal/datastore"
)

func RunPrometheusServer() {
	http.Handle("/metrics", promhttp.Handler())

	addr := ":2112"
	log.Printf("prometheus server listening at %v", addr)
	log.Fatal(http.ListenAndServe(":2112", nil))
}

type MetricsManager struct {
	CacheMisses    prometheus.Counter
	CacheHits      prometheus.Counter
	CacheItemGauge prometheus.Gauge
	ticker         *time.Ticker
	cancelFunc     context.CancelFunc
}

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
	go m.startGaugeUpdater(ctx, cache)
	return m
}

func (m *MetricsManager) Close() {
	m.cancelFunc()
	m.ticker.Stop()
}

func (m *MetricsManager) startGaugeUpdater(ctx context.Context, cache datastore.Cache) {
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
