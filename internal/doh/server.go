package doh

import (
	"io"
	"net/http"

	"github.com/miekg/dns"
)

type Resolver func(*dns.Msg) (*dns.Msg, error)

type Server struct {
	resolve Resolver
}

func New(r Resolver) *Server {
	return &Server{resolve: r}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	// RFC 8484：浏览器直接 GET 必须拒绝
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("Content-Type") != "application/dns-message" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	req := new(dns.Msg)
	if err := req.Unpack(body); err != nil {
		http.Error(w, "invalid dns message", http.StatusBadRequest)
		return
	}

	resp, err := s.resolve(req)
	if err != nil {
		http.Error(w, "dns resolve failed", http.StatusInternalServerError)
		return
	}

	buf, err := resp.Pack()
	if err != nil {
		http.Error(w, "pack failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/dns-message")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
