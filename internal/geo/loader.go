package geo

import (
	"log"
)

type Config struct {
	GeoIPFile           string
	GeoSiteFile         string
	GeoIPDownloadURL    string
	GeoSiteDownloadURL  string
}

func Load(cfg Config) (*Matcher, error) {
	dl := &Downloader{
		GeoIPFile:   cfg.GeoIPFile,
		GeoSiteFile: cfg.GeoSiteFile,
		GeoIPURL:    cfg.GeoIPDownloadURL,
		GeoSiteURL:  cfg.GeoSiteDownloadURL,
	}

	if err := dl.Ensure(); err != nil {
		log.Println("[geo] download failed:", err)
	}

	m := NewMatcher()

	if err := m.LoadGeoIP(cfg.GeoIPFile); err != nil {
		log.Println("[geo] load geoip failed:", err)
	}

	if err := m.LoadGeoSite(cfg.GeoSiteFile); err != nil {
		log.Println("[geo] load geosite failed:", err)
	}

	return m, nil
}
