package dnsserver

import (
	"log"
	"time"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"
	"dns-core/internal/web"

	"github.com/miekg/dns"
)

type Server struct {
	addr    string
	blocker *blocker.Blocker
	cache   *cache.Cache
}

func New(addr string) *Server {
	b := blocker.New()
	_ = b.LoadFile("assets/block/block.txt")
	_ = b.Watch()

	w := web.New(b)
	w.Start(":8080")

	return &Server{
		addr:    addr,
		blocker: b,
		cache:   cache.New(),
	}
}

func (s *Server) Start() error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handle)

	go func() {
		log.Println("TCP DNS listening on", s.addr)
		_ = (&dns.Server{
			Addr:    s.addr,
			Net:     "tcp",
			Handler: handler,
		}).ListenAndServe()
	}()

	log.Println("UDP DNS listening on", s.addr)
	return (&dns.Server{
		Addr:    s.addr,
		Net:     "udp",
		Handler: handler,
	}).ListenAndServe()
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

func (s *Server) handle(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	ensureEDNS0(r)
	q := r.Question[0]

	if s.blocker.IsBlocked(q.Name) {
		_ = w.WriteMsg(blocker.BlockResponse(r))
		return
	}

	if msg := s.cache.Get(q); msg != nil {
		_ = w.WriteMsg(msg)
		return
	}

	c := &dns.Client{Net: w.LocalAddr().Network(), Timeout: 5 * time.Second}
	resp, _, err := c.Exchange(r, "8.8.8.8:53")
	if err != nil {
		return
	}

	s.cache.Set(q, resp)
	_ = w.WriteMsg(resp)
}
