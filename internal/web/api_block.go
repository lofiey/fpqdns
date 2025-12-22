package web

import (
	"encoding/json"
	"net/http"
)

type blockReq struct {
	Domain string `json:"domain"`
}

func (s *Server) blockList(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(s.blocker.List())
}

func (s *Server) blockAdd(w http.ResponseWriter, r *http.Request) {
	var req blockReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if err := s.blocker.Add(req.Domain); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) blockDelete(w http.ResponseWriter, r *http.Request) {
	var req blockReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if err := s.blocker.Delete(req.Domain); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
