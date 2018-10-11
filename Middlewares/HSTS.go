package middlewares

import (
	"net/http"
)

func HSTS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// HSTS is simply a HTTP header (Strict-Transport-Security) that instructs the browser to change all http:// requests to https://.
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	})
}
