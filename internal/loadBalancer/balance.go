package loadbalancer

type Backend string

type Balancer interface {
	NextBackend() (Backend, error)
	SetHealth(backend string, healthy bool)
}

func ResolveBalancer(cfg BalancerConfig) Balancer {
	servers := make([]string, len(cfg.Servers))
	for i, s := range cfg.Servers {
		servers[i] = s.GetURL()
	}

	switch cfg.Balancer {
	case "roundrobin":
		return NewRoundRobin(servers)
	default:
		return NewRoundRobin(servers)
	}
}
