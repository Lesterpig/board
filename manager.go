package main

import (
	"time"

	"github.com/Lesterpig/board/probe"
)

// Service stores several information from a service, especially its last status.
type Service struct {
	Prober  probe.Prober `json:"-"`
	Name    string
	Status  probe.Status
	Message string
	Target  string
}

// Manager stores several services sorted by categories.
type Manager struct {
	LastUpdate time.Time
	Services   map[string]([]*Service)
}

// ProbeLoop starts the main loop that will call ProbeAll regularly.
func (manager *Manager) ProbeLoop(interval time.Duration) {
	manager.ProbeAll()
	c := time.Tick(interval)
	for range c {
		manager.ProbeAll()
	}
}

// ProbeAll triggers the Probe function for each registered service in the manager.
// Everything is done asynchronously.
func (manager *Manager) ProbeAll() {
	log.Debug("Probing all")
	manager.LastUpdate = time.Now()
	for category, services := range manager.Services {
		for _, service := range services {
			go func(category string, service *Service) {
				prevStatus := service.Status
				service.Status, service.Message = service.Prober.Probe()

				if prevStatus != service.Status {
					if service.Status == probe.StatusError {
						AlertAll(category, service)
					} else if prevStatus == probe.StatusError {
						AlertAll(category, service)
					}
				}
			}(category, service)
		}
	}
}

// AlertAll sends an alert signaling the provided service is DOWN.
// It uses global configuration for list of alert (`A` variable).
func AlertAll(category string, service *Service) {
	date := time.Now().Format("15:04:05 MST")
	for _, alert := range alerters {
		alert.Alert(service.Status, category, service.Name, service.Message, service.Target, date)
	}
}
