package geo

import (
	"os"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

type GeoSite struct {
	cn map[string]struct{}
}

func LoadGeoSite(path string) (*GeoSite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list routercommon.GeoSiteList
	if err := proto.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	cn := make(map[string]struct{})
	for _, site := range list.Entry {
		if strings.EqualFold(site.CountryCode, "CN") {
			for _, d := range site.Domain {
				cn[strings.ToLower(d.Value)] = struct{}{}
			}
		}
	}

	return &GeoSite{cn: cn}, nil
}

func (g *GeoSite) IsCN(domain string) bool {
	domain = strings.TrimSuffix(strings.ToLower(domain), ".")
	_, ok := g.cn[domain]
	return ok
}
