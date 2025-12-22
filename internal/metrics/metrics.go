package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Counter struct {
	Total   uint64
	CN      uint64
	Foreign uint64
}

type TopMap struct {
	mu sync.Mutex
	m  map[string]uint64
}

func NewTopMap() *TopMap {
	return &TopMap{m: make(map[string]uint64)}
}

func (t *TopMap) Inc(key string) {
	t.mu.Lock()
	t.m[key]++
	t.mu.Unlock()
}

func (t *TopMap) Top(n int) map[string]uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	res := make(map[string]uint64)
	i := 0
	for k, v := range t.m {
		res[k] = v
		i++
		if i >= n {
			break
		}
	}
	return res
}

type Metrics struct {
	StartTime time.Time

	Query Counter

	CacheHit  uint64
	CacheMiss uint64

	TopDomain *TopMap
	TopClient *TopMap
}

var Global = &Metrics{
	StartTime: time.Now(),
	TopDomain: NewTopMap(),
	TopClient: NewTopMap(),
}

func IncQuery(isCN bool) {
	atomic.AddUint64(&Global.Query.Total, 1)
	if isCN {
		atomic.AddUint64(&Global.Query.CN, 1)
	} else {
		atomic.AddUint64(&Global.Query.Foreign, 1)
	}
}

func IncCache(hit bool) {
	if hit {
		atomic.AddUint64(&Global.CacheHit, 1)
	} else {
		atomic.AddUint64(&Global.CacheMiss, 1)
	}
}
