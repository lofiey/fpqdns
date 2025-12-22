package dnsserver

import (
	"log"

	"github.com/miekg/dns"
)

func (s *Server) listenDoQ(certFile, keyFile string) {
	server := &dns.Server{
		Addr:      ":853",
		Net:       "quic",
		Handler:   dns.HandlerFunc(s.handleDNS),
		TLSConfig: s.tlsConfig(certFile, keyFile),
	}

	log.Println("DoQ listening on :853 (QUIC)")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("DoQ failed:", err)
	}
}
