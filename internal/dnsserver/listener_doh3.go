package dnsserver

import (
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

func (s *Server) listenDoH3(certFile, keyFile string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/dns-query", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	server := http3.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: s.tlsConfig(certFile, keyFile),
	}

	log.Println("DoH3 listening on :443")
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatal("DoH3 failed:", err)
	}
}
