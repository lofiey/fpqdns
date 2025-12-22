package dnsserver

import (
	"log"
	"time"

	"github.com/miekg/dns"
)

type Server struct {
	addr   string
	server *dns.Server
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Start() error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleQuery)

	s.server = &dns.Server{
		Addr:         s.addr,
		Net:          "udp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop() {
	if s.server != nil {
		_ = s.server.Shutdown()
	}
}

func (s *Server) handleQuery(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	c.Net = "udp"

	resp, _, err := c.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Println("upstream error:", err)
		return
	}

	_ = w.WriteMsg(resp)
}
