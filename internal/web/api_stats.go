package web

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"dns-core/internal/metrics"
)

func statsHandler(w http.ResponseWriter, _ *http.Request) {
	m := metrics.Global

	resp := map[string]any{
		"uptime": time.Since(m.StartTime).String(),

		"query": map[string]uint64{
			"total":   m.Query.Total,
			"cn":      m.Query.CN,
			"foreign": m.Query.Foreign,
		},

		"cache": map[string]uint64{
			"hit":  m.CacheHit,
			"miss": m.CacheMiss,
		},

		"top": map[string]any{
			"domain": m.TopDomain.Top(5),
			"client": m.TopClient.Top(5),
		},

		"upstreams": metrics.SnapshotUpstreams(),

		"runtime": map[string]any{
			"goroutine": runtime.NumGoroutine(),
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
