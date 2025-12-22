package web

import (
	"log"
	"net/http"

	"dns-core/internal/blocker"
)

type Server struct {
	blocker *blocker.Blocker
	auth    *Auth
}

func New(b *blocker.Blocker) *Server {
	return &Server{
		blocker: b,
		auth: &Auth{
			Username: "admin",
			Password: "admin123", // ← 后面你可以配置化
		},
	}
}

func (s *Server) Start(addr string) {
	mux := http.NewServeMux()

	// API（全部需要鉴权）
	mux.HandleFunc("/api/stats", s.auth.Middleware(statsHandler))

	mux.HandleFunc("/api/block/list", s.auth.Middleware(s.blockList))
	mux.HandleFunc("/api/block/add", s.auth.Middleware(s.blockAdd))
	mux.HandleFunc("/api/block/delete", s.auth.Middleware(s.blockDelete))

	// 首页
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
