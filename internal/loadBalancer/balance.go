package loadbalancer

type Backend string

type Balancer interface {
	NextBackend() (Backend, error)
}
