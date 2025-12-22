package dnsserver

import (
	"log"
	"time"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"

	"github.com/miekg/dns"
)

type Server struct {
	addr      string
	udpServer *dns.Server
	tcpServer *dns.Server
	blocker   *blocker.Blocker
	cache     *cache.Cache
}

func New(addr string) *Server {
	b := blocker.New()
	_ = b.LoadFile("assets/block/block.txt")

	return &Server{
		addr:    addr,
		blocker: b,
		cache:   cache.New(),
	}
}

func (s *Server) Start() error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleQuery)

	s.udpServer = &dns.Server{
		Addr:         s.addr,
		Net:          "udp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.tcpServer = &dns.Server{
		Addr:         s.addr,
		Net:          "tcp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("TCP DNS listening on", s.addr)
		_ = s.tcpServer.ListenAndServe()
	}()

	log.Println("UDP DNS listening on", s.addr)
	return s.udpServer.ListenAndServe()
}

func (s *Server) Stop() {
	if s.udpServer != nil {
		_ = s.udpServer.Shutdown()
	}
	if s.tcpServer != nil {
		_ = s.tcpServer.Shutdown()
	}
}

func ensureEDNS0(req *dns.Msg) {
	if opt := req.IsEdns0(); opt != nil {
		opt.SetDo()
		return
	}

	opt := &dns.OPT{
		Hdr: dns.RR_Header{
			Name:   ".",
			Rrtype: dns.TypeOPT,
		},
	}
	opt.SetDo()
	opt.SetUDPSize(1232)
	req.Extra = append(req.Extra, opt)
}

func (s *Server) handleQuery(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	ensureEDNS0(r)

	q := r.Question[0]

	// 1️⃣ 先查缓存（含 CNAME 链）
	if msg := s.cache.Get(q); msg != nil {
		_ = w.WriteMsg(msg)
		return
	}
	if msg := cache.ResolveCNAME(q, s.cache); msg != nil {
		_ = w.WriteMsg(msg)
		return
	}

	// 2️⃣ Blocker
	if s.blocker != nil && s.blocker.IsBlocked(q.Name) {
		resp := blocker.BlockResponse(r)
		_ = w.WriteMsg(resp)
		return
	}

	// 3️⃣ 上游查询
	client := &dns.Client{
		Net:     w.LocalAddr().Network(),
		Timeout: 5 * time.Second,
	}

	resp, _, err := client.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Println("upstream error:", err)
		return
	}

	// 4️⃣ 写入缓存（包括 CNAME）
	for _, rr := range resp.Answer {
		switch rr.Header().Rrtype {
		case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME:
			s.cache.Set(dns.Question{
				Name:   rr.Header().Name,
				Qtype: rr.Header().Rrtype,
				Qclass: dns.ClassINET,
			}, resp)
		}
	}

	_ = w.WriteMsg(resp)
}
