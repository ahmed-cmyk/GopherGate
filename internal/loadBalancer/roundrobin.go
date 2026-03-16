package loadbalancer

import (
	"errors"
	"sync/atomic"
)

type RoundRobin struct {
	backends []Backend
	counter  uint64
}

func NewRoundRobin(servers []string) *RoundRobin {
	b := make([]Backend, len(servers))
	for i, s := range servers {
		b[i] = Backend(s)
	}
	return &RoundRobin{
		backends: b,
		counter:  0,
	}
}

func (rr *RoundRobin) NextBackend() (Backend, error) {
	if len(rr.backends) == 0 {
		return "", errors.New("No backends available")
	}

	i := atomic.AddUint64(&rr.counter, 1) - 1
	backend := rr.backends[i%uint64(len(rr.backends))]

	return backend, nil
}
