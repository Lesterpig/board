package manager

import (
	"context"
	"errors"
	"time"

	"github.com/Lesterpig/board/alert"
	"github.com/Lesterpig/board/config"
	"github.com/sirupsen/logrus"

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
	logger     *logrus.Logger
	kubeClient *KubeClient
	LastUpdate time.Time             `json:"LastUpdate"`
	Services   map[string][]*Service `json:"Services,omitempty"`
	Alerts     []alert.Alerter       `json:"Alerts,omitempty"`
}

func NewManager(cfg *config.Config, log *logrus.Logger, client *KubeClient) (*Manager, error) {
	manager := Manager{
		logger:     log,
		kubeClient: client,
	}

	manager.Services = make(map[string][]*Service)
	for _, c := range cfg.Probes {

		prober, err := manager.createProberFromConfig(c)
		if err != nil {
			return nil, err
		}
		manager.appendService(c.Category, &Service{
			Prober: prober,
			Name:   c.Name,
			Target: c.Config.Target,
		})
	}

	manager.Alerts = make([]alert.Alerter, len(cfg.Alerts))
	for _, c := range cfg.Alerts {
		constructor := alert.AlertConstructors[c.Type]
		if constructor == nil {
			return nil, errors.New("unknown alert type: " + c.Type)
		}

		manager.Alerts = append(manager.Alerts, constructor(c))
	}

	m := &manager

	ctx := context.Background()
	if cfg.AutoDiscover != nil && len(cfg.AutoDiscover) > 0 {
		for _, autoDiscoverConfig := range cfg.AutoDiscover {
			fetcher, err := m.kubeClient.Fetch(autoDiscoverConfig.KubernetesResource)

			if err != nil {
				log.Errorf(
					"Error fetching kubernetes resource: %v %v",
					autoDiscoverConfig.KubernetesResource,
					err,
				)
				continue
			}

			adc := autoDiscoverConfig
			m.logger.Infof("fetching resources %v", autoDiscoverConfig.KubernetesResource)
			go func() {
				m.logger.Debug("inside function")
				m.logger.Infof("resource: %v loop interval: %v", adc.KubernetesResource, adc.LoopInterval)
				for {
					ticker := time.NewTicker(adc.LoopInterval)

					select {
					case <-ctx.Done():
						ticker.Stop()
						return

					case <-ticker.C:
						log.Debug("before fetcher()")
						resources := <-fetcher(ctx)
						log.Infof("Fetched resource: %v found: %v", adc.KubernetesResource, len(resources))

						newAutoDiscoveredServices := make([]*Service, 0)
						for _, resource := range resources {
							mappedService, err := manager.mapKubernetesResource(adc.KubernetesResource, resource)
							if err != nil {
								log.Error("Could not create a map the resource to a service", err)
								continue
							}
							m.logger.Debugf("mapped %v resource to service %v", adc.KubernetesResource, mappedService)

							serviceExists := false
							for cat, srvs := range m.Services {
								if serviceExists {
									break // The service exists by configuration, no need to check other categories
								}

								if adc.KubernetesResource == cat {
									continue // go to next category
								}

								for _, srv := range srvs {
									if srv.Target == mappedService.Target {
										m.logger.Infof("service with target %v already exist in map in category %v", mappedService.Target, cat)
										serviceExists = true
										break
									}
								}
							}
							if !serviceExists {
								newAutoDiscoveredServices = append(newAutoDiscoveredServices, mappedService)
							}
						}
						// by replacing the category entirely we will remove any services that are no longer present
						m.Services[adc.KubernetesResource] = newAutoDiscoveredServices
					}
				}
			}()
		}
	}
	m.logger.Info("Done initializaing autodicover routines")
	return m, nil

}

func (m *Manager) appendService(cat string, service *Service) {
	m.Services[cat] = append(m.Services[cat], service)
}

func (m *Manager) mapKubernetesResource(resource string, res Ingress) (*Service, error) {
	opts := make(map[string]interface{})
	opts["VerifyCertificate"] = res.tls

	probeConfig := probe.Config{
		Type: "http",
		Config: probe.ProberConfig{
			Options: opts,
			Target:  res.host + res.path},
		//Category: "AutoDiscovered",
	}

	proper, err := m.createProberFromConfig(probeConfig)

	if err != nil {
		return nil, err
	} // empty configuration will result in defaults

	service := Service{
		Prober: proper,
		Name:   resource + ": " + res.name,
		Target: probeConfig.Config.Target,
	}
	return &service, nil

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
	m.logger.Debug("Probing all")

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

func (m *Manager) createProberFromConfig(c probe.Config) (probe.Prober, error) {
	proberConstructor := probe.ProberConstructors[c.Type]
	if proberConstructor == nil {
		return nil, errors.New("unknown probe type: " + c.Type)
	}

	c.Config = config.SetProbeConfigDefaults(c.Config)

	prober := proberConstructor()
	err := prober.Init(c.Config)
	return prober, err

}
