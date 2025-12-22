package dnsserver

import (
	"log"

	"github.com/miekg/dns"
)

func (s *Server) listenUDP() {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleDNS)

	server := &dns.Server{
		Addr:    s.addr,
		Net:     "udp",
		Handler: handler,
	}

	log.Println("DNS UDP listening on", s.addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("UDP server failed:", err)
	}
}
