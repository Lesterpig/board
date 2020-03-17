package main

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/Lesterpig/board/alert"
	"github.com/Lesterpig/board/probe"

	"github.com/spf13/viper"
)

type serviceConfig struct {
	probe.Config
	Name     string
	Category string
	Probe    string
}

type alertConfig struct {
	Type    string
	Token   string
	Webhook string
	Channel string
}

var probeConstructors = map[string](func() probe.Prober){
	"dns":       func() probe.Prober { return &probe.DNS{} },
	"http":      func() probe.Prober { return &probe.HTTP{} },
	"minecraft": func() probe.Prober { return &probe.Minecraft{} },
	"port":      func() probe.Prober { return &probe.Port{} },
	"smtp":      func() probe.Prober { return &probe.SMTP{} },
}

var alertConstructors = map[string](func(c alertConfig) alert.Alerter){
	"pushbullet": func(c alertConfig) alert.Alerter {
		return alert.NewPushbullet(c.Token)
	},
	"slack": func(c alertConfig) alert.Alerter {
		return alert.NewSlack(c.Webhook, c.Channel)
	},
}

var alerters []alert.Alerter

func parseConfigString(cnf string) (dir string, name string) {
	dir = filepath.Dir(cnf)
	basename := filepath.Base(cnf)
	name = strings.TrimSuffix(basename, filepath.Ext(basename))

	return
}

func loadConfig(configPath, configName string) (*Manager, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	sc := make([]serviceConfig, 0)

	err = viper.UnmarshalKey("services", &sc)
	if err != nil {
		return nil, err
	}

	manager := Manager{}
	manager.Services = make(map[string]([]*Service))

	for _, c := range sc {
		constructor := probeConstructors[c.Probe]
		if constructor == nil {
			return nil, errors.New("unknown probe type: " + c.Probe)
		}

		c.Config = setProbeConfigDefaults(c.Config)

		prober := constructor()

		err = prober.Init(c.Config)
		if err != nil {
			return nil, err
		}

		manager.Services[c.Category] = append(manager.Services[c.Category], &Service{
			Prober: prober,
			Name:   c.Name,
			Target: c.Target,
		})
	}

	ac := make([]alertConfig, 0)

	err = viper.UnmarshalKey("alerts", &ac)
	if err != nil {
		return nil, err
	}

	alerters = make([]alert.Alerter, 0)

	for _, c := range ac {
		constructor := alertConstructors[c.Type]
		if constructor == nil {
			return nil, errors.New("unknown alert type: " + c.Type)
		}

		alerters = append(alerters, constructor(c))
	}

	return &manager, err
}

func setProbeConfigDefaults(c probe.Config) probe.Config {
	if c.Warning == 0 {
		c.Warning = 500 * time.Millisecond
	}

	if c.Fatal == 0 {
		c.Fatal = time.Minute
	}

	return c
}
