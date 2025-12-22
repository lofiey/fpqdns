package doh

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func ListenHTTP3(addr string, handler http.Handler, tlsConfig *tls.Config) error {
	server := http3.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	log.Println("DoH HTTP/3 listening on", addr)
	return server.ListenAndServe()
}
