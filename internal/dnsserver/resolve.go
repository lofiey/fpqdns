package dnsserver

import (
	"time"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"
	"dns-core/internal/ecs"
	"dns-core/internal/geo"
	"dns-core/internal/metrics"
	"dns-core/internal/upstream"

	"github.com/miekg/dns"
)

type Resolver struct {
	blocker *blocker.Blocker
	cache   *cache.Cache
	geosite *geo.GeoSite
}

func NewResolver(b *blocker.Blocker, c *cache.Cache) *Resolver {
	gs, _ := geo.LoadGeoSite("assets/geo/geosite.dat")

	return &Resolver{
		blocker: b,
		cache:   c,
		geosite: gs,
	}
}

func ensureEDNS0(m *dns.Msg) {
	if opt := m.IsEdns0(); opt != nil {
		opt.SetDo()
		return
	}

	o := &dns.OPT{}
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetUDPSize(1232)
	o.SetDo()
	m.Extra = append(m.Extra, o)
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) == 0 {
		return nil, nil
	}

	ensureEDNS0(req)

	q := req.Question[0]
	domain := q.Name

	// Blocker
	if r.blocker.IsBlocked(domain) {
		return blocker.BlockResponse(req), nil
	}

	// Geo 判断
	isCN := r.geosite != nil && r.geosite.IsCN(domain)

	// Metrics：查询计数
	metrics.IncQuery(isCN)
	metrics.Global.TopDomain.Inc(domain)

	// Cache
	if msg := r.cache.Get(q); msg != nil {
		metrics.IncCache(true)
		return msg, nil
	}
	metrics.IncCache(false)

	// 国内上游
	if isCN {
		ecs.Attach(req, "1.2.4.8", 24) // 示例国内 ECS

		c := &dns.Client{
			Net:     "udp",
			Timeout: 5 * time.Second,
		}

		resp, _, err := c.Exchange(req, "223.5.5.5:53")
		if err != nil {
			return nil, err
		}

		r.cache.Set(q, resp)
		return resp, nil
	}

	// 海外并发竞速
	ecs.Attach(req, "8.8.8.8", 24)

	ups := []upstream.Upstream{
		{Addr: "8.8.8.8:53", Net: "udp"},
		{Addr: "1.1.1.1:53", Net: "udp"},
		{Addr: "9.9.9.9:53", Net: "udp"},
	}

	resp, err := upstream.Race(req, ups)
	if err != nil {
		return nil, err
	}

	r.cache.Set(q, resp)
	return resp, nil
}
