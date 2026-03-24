package loadbalancer

import (
	"context"

	"github.com/ahmed-cmyk/GopherGate/internal/health"
)

type Selector interface {
	Next(ctx context.Context) (Backend, error)
}

type HealthAware interface {
	SetHealth(backend string, healthy bool)
	IsHealthy(backend string) bool
}

type LoadBalancer interface {
	Selector
	HealthAware
}

type BalancerConfig struct {
	Path     string
	Balancer string
	Servers  []health.Target
}
