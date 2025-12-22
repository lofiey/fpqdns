package dot

import (
	"crypto/tls"
	"log"

	"github.com/miekg/dns"
)

type Resolver func(*dns.Msg) (*dns.Msg, error)

func Listen(addr string, tlsCfg *tls.Config, r Resolver) error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		resp, err := r(req)
		if err == nil && resp != nil {
			_ = w.WriteMsg(resp)
		}
	})

	server := &dns.Server{
		Addr:      addr,
		Net:       "tcp-tls",
		TLSConfig: tlsCfg,
		Handler:   handler,
	}

	log.Println("DoT listening on", addr)
	return server.ListenAndServe()
}
