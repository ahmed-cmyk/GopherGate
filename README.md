# 🚀 GopherGate: A Go-Powered API Gateway

A lightweight, high-performance API Gateway built from scratch in Go. This project is a deep dive into Go's `net/http` package, concurrency patterns, and distributed systems architecture.

## 📖 Overview

**GopherGate** acts as the single entry point for your microservices. It handles the "cross-cutting concerns" like authentication, rate limiting, and request routing so your backend services don't have to.

### Why build this?

* To master the `httputil.ReverseProxy` standard library.
* To implement common distributed system patterns (Circuit Breakers, Retries).
* To practice writing high-performance middleware.

---

## 🛠 Feature Roadmap

I am building this incrementally. Here is the current status of the gateway:

### Phase 1: Core Routing

* [x] **Reverse Proxy:** Forwarding traffic to backend targets.
* [x] **Dynamic Path Matching:** Routing based on URL prefixes (e.g., `/users/*`).
* [x] **Method Filtering:** Restricting routes to specific HTTP verbs (GET, POST, etc.).
* [x] **Header Transformation:** Adding/Removing headers before forwarding.

### Phase 2: Traffic Control

* [x] **Load Balancing:** Round-robin distribution across multiple service instances.
* [x] **Rate Limiting:** Protect backends using the Token Bucket algorithm.
* [x] **Health Checks:** Automatically removing unhealthy backends from the pool.

### Phase 3: Security & Ops

* [x] **JWT Validation:** Centralized authentication at the edge.
* [x] **CORS Middleware:** Global configuration for cross-origin requests.
* [x] **Structured Logging:** Zap or Logrus integration for JSON logging.
* [x] **Prometheus Metrics:** Tracking request latency and 5xx errors.

### Phase 4: Continued Improvements

* [ ] **Rate Limiting:** Implement a "Token Bucket" or "Leaky Bucket" algorithm to prevent DDoS and API abuse (using x/time/rate).
* [ ] **Active/Passive Health Checks:** Moving beyond simple pings to observing real-time failures and automatically pulling bad backends from the rotation.
* [ ] **Retries & Timeouts:** Implementing "Smart Retries" with exponential backoff so a single blip in a backend doesn't result in a 502 Bad Gateway for the user.
* [ ] **Graceful Shutdown:** Ensuring the http.Server drains active connections before the process exits during a deployment.
* [ ] **Hot Reloading:** Updating the allowMap and routes dynamically from the config file without restarting the entire binary (using fsnotify).

---

## 🚦 Getting Started

### Prerequisites

* Go 1.21+
* Make (optional)

### Installation

```bash
git clone https://github.com/ahmed-cmyk/gophergate.git
cd gophergate
go mod download

```

### Running the Gateway

```bash
go run main.go --config config.yaml

```

---

## ⚙️ Configuration Example

The gateway is configured via a simple YAML file:

```yaml
server:
  port: 8080

routes:
  - path: /api/v1/users
    target: "http://user-service:8081"
    strip_prefix: true
    middlewares:
      - rate_limit
      - auth

```

---

## 🧪 Testing

To run the suite of unit and integration tests:

```bash
go test ./... -v

```

---
