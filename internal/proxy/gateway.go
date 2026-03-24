package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"

	config "github.com/ahmed-cmyk/GopherGate/internal"
	"github.com/ahmed-cmyk/GopherGate/internal/health"
	loadbalancer "github.com/ahmed-cmyk/GopherGate/internal/loadBalancer"
	"github.com/ahmed-cmyk/GopherGate/internal/middleware"
	"github.com/charmbracelet/log"
)

type routeEntry struct {
	balancer loadbalancer.Balancer
	handler  http.Handler
	methods  []string
}

type Gateway struct {
	routes map[string]routeEntry
}

func NewGateway(cfg *config.Config, routeMap *Routes) *Gateway {
	gw := &Gateway{
		routes: make(map[string]routeEntry),
	}

	for _, route := range cfg.Routes {
		servers := routeMap.GetServersForPath(route.Path)

		entry, err := gw.buildRouteEntry(&route, servers)
		if err != nil {
			log.Errorf("Skipping route %s due to error: %v", route.Path, err)
			continue
		}

		gw.routes[route.Path] = entry
	}
	return gw
}

func (gw *Gateway) buildRouteEntry(route *config.Route, servers []health.Target) (routeEntry, error) {
	targetURL, err := url.Parse(route.Targets[0])
	if err != nil {
		return routeEntry{}, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = AddDirector(route, originalDirector)

	finalHandler := middleware.AddMiddlewareChain(proxy, route.Middlewares)

	balancer := loadbalancer.ResolveBalancer(loadbalancer.BalancerConfig{
		Path:     route.Path,
		Balancer: route.Balancer,
		Servers:  servers,
	})

	return routeEntry{
		balancer: balancer,
		handler:  finalHandler,
		methods:  route.Methods,
	}, nil
}

func (gw *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var found bool
	var matched routeEntry

	for path, entry := range gw.routes {
		if strings.HasPrefix(r.URL.Path, path) {
			matched = entry
			found = true
			break
		}
	}

	// If no route has been found return a 404 error
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Check if method is supported
	if len(matched.methods) > 0 && !slices.Contains(matched.methods, r.Method) {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	// Select next backend for the matched route, if nothing matches return 500 error
	host, err := matched.balancer.NextBackend()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	r.URL.Host = string(host)
	r.URL.Scheme = "http"
	r.Host = string(host)

	log.Infof("Routing %s request to %s", r.URL.Path, host)

	matched.handler.ServeHTTP(w, r)
}

// UpdateBackendHealth updates the health status of a backend across all routes
func (gw *Gateway) UpdateBackendHealth(backendURL string, healthy bool) {
	for _, entry := range gw.routes {
		entry.balancer.SetHealth(backendURL, healthy)
	}
}
