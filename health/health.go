package health

import (
	"context"
	"sync"
)

var (
	Up   = "up"
	Down = "down"
)

var (
	healthMu  sync.RWMutex
	healthily = make([]map[string]Healthier, 0)
)

type Healthier interface {
	Ping(ctx context.Context) error
}

type Detail struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Health struct {
	Status  string   `json:"status"`
	Details []Detail `json:"details,omitempty"`
}

// Register  a service to monitor health status
func Register(name string, health Healthier) {
	healthMu.Lock()
	defer healthMu.Unlock()
	if health == nil {
		panic("health: Register health is nil")
	}
	for _, nh := range healthily {
		if _, dup := nh[name]; dup {
			panic("health: Register called twice for health " + name)
		}
	}
	healthily = append(healthily, map[string]Healthier{name: health})
}

// Ping return Health status
func Ping(ctx context.Context) (health Health) {
	healthMu.RLock()
	defer healthMu.RUnlock()
	health.Status = Up

	var details []Detail
	for _, nh := range healthily {
		for n, h := range nh {
			var errStr string
			status := Up
			if err := h.Ping(ctx); err != nil {
				health.Status = Down
				status = Down
				errStr = err.Error()
			}
			details = append(details, Detail{Name: n, Status: status, Error: errStr})
		}
	}
	health.Details = details[:]
	return health
}
