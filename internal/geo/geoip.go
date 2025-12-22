package geo

type GeoIP struct {
	// 预留：后续按客户端 IP 判断
}

func LoadGeoIP(path string) (*GeoIP, error) {
	return &GeoIP{}, nil
}
