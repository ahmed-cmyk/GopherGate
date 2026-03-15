package middleware

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func NewLimiter(limit rate.Limit, bucket int) IPRateLimiter {
	return IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   limit,
		b:   bucket,
	}
}

func (l *IPRateLimiter) fetchLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok := l.ips[ip]
	if !ok {
		l.ips[ip] = rate.NewLimiter(l.r, l.b)
		return l.ips[ip]
	}

	return limiter
}

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func RateLimit(manager *IPRateLimiter) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := getIP(r)
			if err != nil {
				http.Error(w, "Identity Unknown", http.StatusInternalServerError)
				return
			}

			// Now 'manager' is accessible here!
			limiter := manager.fetchLimiter(ip)

			if !limiter.Allow() {
				http.Error(w, "Slow down, Gopher.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
