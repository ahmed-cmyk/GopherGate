package middleware

import "net/http"

func CORS(allowlist map[string]struct{}) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Handle server-to-server and non-browser requests
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			// If the origin isn't in the allowlist return a status "forbidden" error
			_, found := allowlist[origin]
			if !found {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Set the essential headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Crucial: Tell caches the response depends on the Origin header
			w.Header().Set("Vary", "Origin")

			// Handle Preflight (OPTIONS)
			if r.Method == http.MethodOptions {
				// You MUST tell the browser which methods and headers are okay
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
				w.Header().Set("Access-Control-Max-Age", "300") // Cache preflight for 5 mins

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
