package loadbalancer

import (
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
	bLen := len(rr.backends)

	if rr.counter == uint64(bLen) {
		atomic.StoreUint64(&rr.counter, 0)
	}

	backend := rr.backends[rr.counter]

	atomic.AddUint64(&rr.counter, 1)

	return backend, nil
}
