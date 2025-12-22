package cache

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type entry struct {
	msg    *dns.Msg
	expire time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]*entry
}

func New() *Cache {
	return &Cache{
		items: make(map[string]*entry),
	}
}

func key(q dns.Question) string {
	return q.Name + "|" + dns.TypeToString[q.Qtype]
}

func minTTL(msg *dns.Msg) uint32 {
	min := uint32(0)
	for _, rr := range msg.Answer {
		if min == 0 || rr.Header().Ttl < min {
			min = rr.Header().Ttl
		}
	}
	if min == 0 {
		min = 30
	}
	return min
}

func (c *Cache) Get(q dns.Question) *dns.Msg {
	k := key(q)

	c.mu.RLock()
	e, ok := c.items[k]
	c.mu.RUnlock()

	if !ok || time.Now().After(e.expire) {
		return nil
	}

	// 必须拷贝，避免并发修改
	return e.msg.Copy()
}

func (c *Cache) Set(q dns.Question, msg *dns.Msg) {
	ttl := minTTL(msg)
	if ttl == 0 {
		return
	}

	c.mu.Lock()
	c.items[key(q)] = &entry{
		msg:    msg.Copy(),
		expire: time.Now().Add(time.Duration(ttl) * time.Second),
	}
	c.mu.Unlock()
}
