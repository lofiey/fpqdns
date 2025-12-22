package dnsserver

import (
	"io"
	"log"
	"net/http"

	"github.com/miekg/dns"
)

func (s *Server) listenDoH() {
	mux := http.NewServeMux()

	mux.HandleFunc("/dns-query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		req := new(dns.Msg)
		if err := req.Unpack(body); err != nil {
			http.Error(w, "bad dns message", http.StatusBadRequest)
			return
		}

		resp := s.resolve(req)
		raw, _ := resp.Pack()

		w.Header().Set("Content-Type", "application/dns-message")
		_, _ = w.Write(raw)
	})

	go func() {
		log.Println("DoH listening on :443")
		if err := http.ListenAndServe(":443", mux); err != nil {
			log.Fatal("DoH failed:", err)
		}
	}()
}
