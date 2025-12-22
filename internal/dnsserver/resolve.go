package dnsserver

import (
	"time"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"

	"github.com/miekg/dns"
)

type Resolver struct {
	blocker *blocker.Blocker
	cache   *cache.Cache
}

func NewResolver(b *blocker.Blocker, c *cache.Cache) *Resolver {
	return &Resolver{blocker: b, cache: c}
}

func ensureEDNS0(m *dns.Msg) {
	if opt := m.IsEdns0(); opt != nil {
		opt.SetDo()
		return
	}
	o := &dns.OPT{}
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetDo()
	o.SetUDPSize(1232)
	m.Extra = append(m.Extra, o)
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) == 0 {
		return nil, nil
	}

	ensureEDNS0(req)
	q := req.Question[0]

	// Block
	if r.blocker.IsBlocked(q.Name) {
		return blocker.BlockResponse(req), nil
	}

	// Cache
	if msg := r.cache.Get(q); msg != nil {
		return msg, nil
	}

	c := &dns.Client{Net: "udp", Timeout: 5 * time.Second}
	resp, _, err := c.Exchange(req, "8.8.8.8:53")
	if err != nil {
		return nil, err
	}

	r.cache.Set(q, resp)
	return resp, nil
}
