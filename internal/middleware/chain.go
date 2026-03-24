package middleware

import (
	"net/http"

	"github.com/charmbracelet/log"
)

func AddMiddlewareChain(target http.Handler, middlewares []string) http.Handler {
	current := target

	for _, middleware := range middlewares {
		if mwFunc, ok := DefaultRegistry.Get(middleware); ok {
			current = mwFunc(current)
		} else {
			log.Errorf("Warning: Middleware %s not found", middleware)
		}
	}

	return current
}
