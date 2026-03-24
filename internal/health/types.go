package health

import (
	"time"
)

type StatusChangeCallback func(url string, healthy bool)

type HealthChecker struct {
	routes    RouteStore
	interval  time.Duration
	pinger    *Pinger
	callbacks []StatusChangeCallback
}

type RouteStore interface {
	GetPaths() []string
	GetServersForPath(path string) []Target
}

type Target interface {
	GetURL() string
	GetStatus() bool
	SetStatus(alive bool)
}

func (hc *HealthChecker) OnStatusChange(cb StatusChangeCallback) {
	hc.callbacks = append(hc.callbacks, cb)
}

func (hc *HealthChecker) notifyStatusChange(url string, healthy bool) {
	for _, cb := range hc.callbacks {
		cb(url, healthy)
	}
}
