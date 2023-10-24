package main

import (
	"errors"
	"github.com/Lesterpig/board/alert"
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

// Manager stores several Services sorted by categories.
type Manager struct {
	LastUpdate time.Time             `json:"LastUpdate"`
	Services   map[string][]*Service `json:"Probes,omitempty"`
	Alerts     []alert.Alerter       `json:"Alerts,omitempty"`
}

func NewManager(config *Config) (*Manager, error) {
	manager := Manager{}

	manager.Services = make(map[string][]*Service)
	for _, c := range config.Probes {
		probeConstructor := probe.ProbeConstructors[c.Type]
		if probeConstructor == nil {
			return nil, errors.New("unknown probe type: " + c.Type)
		}

		c.Config = setProbeConfigDefaults(c.Config)

		prober := probeConstructor()

		err := prober.Init(c.Config)
		if err != nil {
			return nil, err
		}

		manager.Services[c.Category] = append(manager.Services[c.Category], &Service{
			Prober: prober,
			Name:   c.Name,
			Target: c.Config.Target,
		})
	}

	manager.Alerts = make([]alert.Alerter, len(config.Alerts))
	for _, c := range config.Alerts {
		constructor := alert.AlertConstructors[c.Type]
		if constructor == nil {
			return nil, errors.New("unknown alert type: " + c.Type)
		}

		manager.Alerts = append(manager.Alerts, constructor(c))
	}

	m := &manager

	return m, nil

}

// ProbeLoop starts the main loop that will call ProbeAll regularly.
func (m *Manager) ProbeLoop(interval time.Duration) {
	m.probeAll()

	m.LastUpdate = time.Now()

	c := time.Tick(interval)
	for range c {
		m.probeAll()
	}
}

// ProbeAll triggers the probe function for each registered service in the manager.
// Everything is done asynchronously.
func (m *Manager) probeAll() {
	log.Debug("Probing all")

	m.LastUpdate = time.Now()

	for category, services := range m.Services {
		for _, service := range services {
			go func(category string, service *Service) {
				prevStatus := service.Status
				service.Status, service.Message = service.Prober.Probe()

				if prevStatus != service.Status {
					if service.Status == probe.StatusError {
						m.AlertAll(category, service)
					} else if prevStatus == probe.StatusError {
						m.AlertAll(category, service)
					}
				}
			}(category, service)
		}
	}
}

// AlertAll sends an alert signaling the provided service is DOWN.
// It uses global configuration for list of alert (`A` variable).
func (m *Manager) AlertAll(category string, service *Service) {
	date := time.Now().Format("15:04:05 MST")

	for _, alerter := range m.Alerts {
		alerter.Alert(service.Status, category, service.Name, service.Message, service.Target, date)
	}
}
