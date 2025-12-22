package web

import (
	"encoding/json"
	"log"
	"net/http"

	"dns-core/internal/blocker"
)

type Server struct {
	blocker *blocker.Blocker
}

func New(b *blocker.Blocker) *Server {
	return &Server{blocker: b}
}

func (s *Server) Start(addr string) {
	http.HandleFunc("/", s.index)
	http.HandleFunc("/block/list", s.list)
	http.HandleFunc("/block/add", s.add)
	http.HandleFunc("/block/del", s.del)

	log.Println("Web panel listening on", addr)
	go http.ListenAndServe(addr, nil)
}

func (s *Server) index(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("dns-core web panel"))
}

func (s *Server) list(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(s.blocker.List())
}

func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "missing domain", 400)
		return
	}
	_ = s.blocker.Add(domain)
	w.WriteHeader(200)
}

func (s *Server) del(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "missing domain", 400)
		return
	}
	_ = s.blocker.Remove(domain)
	w.WriteHeader(200)
}
