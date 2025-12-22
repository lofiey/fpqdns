package dnsserver

import (
	"log"
	"time"

	"dns-core/internal/blocker"

	"github.com/miekg/dns"
)

type Server struct {
	addr      string
	udpServer *dns.Server
	tcpServer *dns.Server
	blocker   *blocker.Blocker
}

func New(addr string) *Server {
	b := blocker.New()
	_ = b.LoadFile("assets/block/block.txt")

	return &Server{
		addr:    addr,
		blocker: b,
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
		if err := s.tcpServer.ListenAndServe(); err != nil {
			log.Println("TCP DNS stopped:", err)
		}
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

	// ✅ 确保 DNSSEC 能力透传
	ensureEDNS0(r)

	q := r.Question[0]
	domain := q.Name

	// 1️⃣ 拦截优先
	if s.blocker != nil && s.blocker.IsBlocked(domain) {
		resp := blocker.BlockResponse(r)
		_ = w.WriteMsg(resp)
		return
	}

	// 2️⃣ 上游转发（保持 DO 位）
	client := &dns.Client{
		Net:     w.LocalAddr().Network(),
		Timeout: 5 * time.Second,
	}

	resp, _, err := client.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Println("upstream exchange error:", err)
		return
	}

	_ = w.WriteMsg(resp)
}
