package web

import (
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
	mux := http.NewServeMux()
    mux.HandleFunc("/api/block/list", s.blockList)
    mux.HandleFunc("/api/block/add", s.blockAdd)
    mux.HandleFunc("/api/block/delete", s.blockDelete)
	// API
	mux.HandleFunc("/api/stats", statsHandler)

	// 预留：block 管理 / 登录 / 配置
	// mux.HandleFunc("/api/block", s.blockHandler)

	// 首页（占位）
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("DNS Core Web Panel"))
	})

	go func() {
		log.Println("Web panel listening on", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatal(err)
		}
	}()
}
