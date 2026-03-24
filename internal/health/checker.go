package health

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
)

func NewHealthChecker(routes RouteStore, interval time.Duration, timeout time.Duration) *HealthChecker {
	pinger := NewPinger(timeout)

	return &HealthChecker{
		routes:    routes,
		interval:  interval,
		pinger:    pinger,
		callbacks: make([]StatusChangeCallback, 0),
	}
}

func (hc *HealthChecker) StartHealthChecker(ctx context.Context) {
	healthTicker := time.NewTicker(hc.interval)
	defer healthTicker.Stop()

	for {
		select {
		case <-healthTicker.C:
			paths := hc.routes.GetPaths()

			for _, path := range paths {
				servers := hc.routes.GetServersForPath(path)
				for _, server := range servers {
					// Capture previous status before ping
					prevStatus := server.GetStatus()

					// Launch each ping in its own goroutine
					go func(s Target, wasHealthy bool) {
						log.Debugf("Pinging: %s", s.GetURL())
						hc.pinger.Ping(ctx, s)

						// Notify callbacks if status changed
						isHealthy := s.GetStatus()
						if isHealthy != wasHealthy {
							hc.notifyStatusChange(s.GetURL(), isHealthy)
						}
					}(server, prevStatus)
				}
			}
		case <-ctx.Done():
			log.Info("Health checker stopped")
			return
		}
	}
}
