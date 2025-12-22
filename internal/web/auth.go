package web

import (
	"crypto/subtle"
	"net/http"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) check(r *http.Request) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(u), []byte(a.Username)) != 1 {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(p), []byte(a.Password)) != 1 {
		return false
	}
	return true
}

func (a *Auth) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.check(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="DNS Core"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
