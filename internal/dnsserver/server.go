package dnsserver

import (
	"log"
	"time"

	"github.com/miekg/dns"
)

type Server struct {
	addr      string
	udpServer *dns.Server
	tcpServer *dns.Server
}

func New(addr string) *Server {
	return &Server{addr: addr}
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
		if err := s.tcpServer.ListenAndServe(); err != nil {
			log.Println("tcp dns failed:", err)
		}
	}()

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

func (s *Server) handleQuery(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	c.Net = w.RemoteAddr().Network()

	resp, _, err := c.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Println("upstream error:", err)
		return
	}

	_ = w.WriteMsg(resp)
}
