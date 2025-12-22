package upstream

import (
	"time"

	"dns-core/internal/metrics"

	"github.com/miekg/dns"
)

type Upstream struct {
	Addr string
	Net  string
}

func (u *Upstream) Exchange(req *dns.Msg) (*dns.Msg, time.Duration, error) {
	c := &dns.Client{
		Net:     u.Net,
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, _, err := c.Exchange(req, u.Addr)
	latency := time.Since(start)

	if err != nil {
		metrics.RecordUpstream(u.Addr, latency, false)
		return nil, latency, err
	}

	metrics.RecordUpstream(u.Addr, latency, true)
	return resp, latency, nil
}
