package metrics

import (
	"sort"
	"sync"
	"time"
)

type TopItem struct {
	Key   string `json:"key"`
	Count uint64 `json:"count"`
}

type counter struct {
	count uint64
	last  int64
}

type TopCounter struct {
	mu   sync.Mutex
	data map[string]*counter
	ttl  int64
}

func NewTopCounter(ttl time.Duration) *TopCounter {
	return &TopCounter{
		data: make(map[string]*counter),
		ttl:  int64(ttl.Seconds()),
	}
}

func (t *TopCounter) Add(key string) {
	now := time.Now().Unix()

	t.mu.Lock()
	defer t.mu.Unlock()

	if c, ok := t.data[key]; ok {
		c.count++
		c.last = now
		return
	}

	t.data[key] = &counter{
		count: 1,
		last:  now,
	}
}

func (t *TopCounter) cleanup(now int64) {
	for k, v := range t.data {
		if now-v.last > t.ttl {
			delete(t.data, k)
		}
	}
}

func (t *TopCounter) Top(n int) []TopItem {
	now := time.Now().Unix()

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cleanup(now)

	items := make([]TopItem, 0, len(t.data))
	for k, v := range t.data {
		items = append(items, TopItem{
			Key:   k,
			Count: v.count,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if len(items) > n {
		items = items[:n]
	}

	return items
}
