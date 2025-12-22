package dnsserver

import (
	"log"

	"github.com/miekg/dns"
)

func (s *Server) listenTCP() {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleDNS)

	server := &dns.Server{
		Addr:    s.addr,
		Net:     "tcp",
		Handler: handler,
	}

	log.Println("DNS TCP listening on", s.addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("TCP server failed:", err)
	}
}
