package main

import (
	"errors"
	"time"

	"github.com/lesterpig/board/probe"

	"github.com/spf13/viper"
)

type config struct {
	probe.Config
	Name     string
	Category string
	Probe    string
}

var probeConstructors = map[string](func() probe.Prober){
	"dns":       func() probe.Prober { return &probe.DNS{} },
	"http":      func() probe.Prober { return &probe.HTTP{} },
	"minecraft": func() probe.Prober { return &probe.Minecraft{} },
	"port":      func() probe.Prober { return &probe.Port{} },
	"smtp":      func() probe.Prober { return &probe.SMTP{} },
}

func loadConfig() (Manager, error) {
	viper.SetConfigName("board")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	servicesConfig := make([]config, 0)
	err = viper.UnmarshalKey("services", &servicesConfig)
	if err != nil {
		return nil, err
	}

	manager := make(Manager)
	for _, c := range servicesConfig {
		constructor := probeConstructors[c.Probe]
		if constructor == nil {
			return nil, errors.New("unknown probe type: " + c.Probe)
		}

		// Update default fields
		if c.Warning == 0 {
			c.Warning = 500 * time.Millisecond
		}

		if c.Fatal == 0 {
			c.Fatal = time.Minute
		}

		prober := constructor()
		err = prober.Init(c.Config)
		if err != nil {
			return nil, err
		}

		manager[c.Category] = append(manager[c.Category], &Service{
			Prober: prober,
			Name:   c.Name,
		})
	}

	return manager, err
}
