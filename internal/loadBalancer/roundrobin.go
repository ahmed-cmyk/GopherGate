package loadbalancer

import (
	"errors"
	"sync"
	"sync/atomic"
)

type RoundRobin struct {
	backends []Backend
	healthy  []bool
	counter  uint64
	mu       sync.RWMutex
}

func NewRoundRobin(servers []string) *RoundRobin {
	b := make([]Backend, len(servers))
	h := make([]bool, len(servers))
	for i, s := range servers {
		b[i] = Backend(s)
		h[i] = true // Default to healthy
	}
	return &RoundRobin{
		backends: b,
		healthy:  h,
		counter:  0,
	}
}

func (rr *RoundRobin) NextBackend() (Backend, error) {
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	if len(rr.backends) == 0 {
		return "", errors.New("no backends available")
	}

	// Try to find a healthy backend
	for attempts := 0; attempts < len(rr.backends); attempts++ {
		idx := int(atomic.AddUint64(&rr.counter, 1)-1) % len(rr.backends)
		if rr.healthy[idx] {
			return rr.backends[idx], nil
		}
	}

	return "", errors.New("no healthy backends available")
}

func (rr *RoundRobin) SetHealth(backend string, healthy bool) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	for i, b := range rr.backends {
		if string(b) == backend {
			rr.healthy[i] = healthy
			return
		}
	}
}
