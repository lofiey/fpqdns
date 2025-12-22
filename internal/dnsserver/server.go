package dnsserver

import (
	"log"
	"sync"

	"dns-core/internal/geo"
	"dns-core/internal/web"
)

type Server struct {
	addr string

	geo *geo.Matcher

	once sync.Once
}

// New 创建 DNS Server
func New(addr string) *Server {
	// ===== Geo 自动下载 + 加载 =====
	g, err := geo.Load(geo.Config{
		GeoIPFile:          "data/geoip.dat",
		GeoSiteFile:        "data/geosite.dat",
		GeoIPDownloadURL:   "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/geoip.dat",
		GeoSiteDownloadURL: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/geosite.dat",
	})
	if err != nil {
		log.Println("[dns] geo load error:", err)
	}

	return &Server{
		addr: addr,
		geo:  g,
	}
}

// Start 启动所有服务（只会执行一次）
func (s *Server) Start() error {
	var err error

	s.once.Do(func() {
		// ===== DNS 核心 =====
		go s.listenUDP()
		go s.listenTCP()
		go s.listenDoH()

		// ===== Web 管理面板 =====
		go func() {
			w := web.New(nil) // blocker 已在内部管理
			w.Start(":8080")
		}()

		log.Println("[dns] server started")
	})

	return err
}

// GeoMatcher 暴露给 resolver 使用
func (s *Server) GeoMatcher() *geo.Matcher {
	return s.geo
}
