package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type UpstreamStat struct {
	Addr string

	Success uint64
	Failure uint64

	TotalLatency int64 // 纳秒累计
}

type UpstreamMetrics struct {
	mu sync.Mutex
	m  map[string]*UpstreamStat
}

var Upstreams = &UpstreamMetrics{
	m: make(map[string]*UpstreamStat),
}

func (u *UpstreamMetrics) get(addr string) *UpstreamStat {
	u.mu.Lock()
	defer u.mu.Unlock()

	if s, ok := u.m[addr]; ok {
		return s
	}

	s := &UpstreamStat{Addr: addr}
	u.m[addr] = s
	return s
}

func RecordUpstream(addr string, latency time.Duration, success bool) {
	s := Upstreams.get(addr)

	if success {
		atomic.AddUint64(&s.Success, 1)
		atomic.AddInt64(&s.TotalLatency, latency.Nanoseconds())
	} else {
		atomic.AddUint64(&s.Failure, 1)
	}
}

func SnapshotUpstreams() []map[string]any {
	Upstreams.mu.Lock()
	defer Upstreams.mu.Unlock()

	out := make([]map[string]any, 0, len(Upstreams.m))
	for _, s := range Upstreams.m {
		total := atomic.LoadUint64(&s.Success)
		var avgMs int64
		if total > 0 {
			avgMs = atomic.LoadInt64(&s.TotalLatency) / int64(total) / 1e6
		}

		out = append(out, map[string]any{
			"addr":       s.Addr,
			"success":    atomic.LoadUint64(&s.Success),
			"failure":    atomic.LoadUint64(&s.Failure),
			"avg_latency_ms": avgMs,
		})
	}
	return out
}
