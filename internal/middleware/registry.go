package middleware

import (
	"net/http"

	loadbalancer "github.com/ahmed-cmyk/GopherGate/internal/loadBalancer"
)

type MiddlewareFunc func(http.Handler) http.Handler

var Registry = map[string]MiddlewareFunc{
	"logging": Logging,
}

func ResolveBalancer(path, balancer string, servers []string) loadbalancer.Balancer {
	// For now there is only one load balancer hence we only return one
	return loadbalancer.NewRoundRobin(servers)
}
