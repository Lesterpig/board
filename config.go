package main

import (
	"github.com/Lesterpig/board/alert"
	"path/filepath"
	"strings"
	"time"

	"github.com/Lesterpig/board/probe"

	"github.com/spf13/viper"
)

type Config struct {
	AutoDiscover AutoDiscoverConfig
	Probes       []probe.Config
	Alerts       []alert.AlertConfig
}
type AutoDiscoverConfig struct {
	Ingres bool
}

func parseConfigString(cnf string) (dir string, name string) {
	dir = filepath.Dir(cnf)
	basename := filepath.Base(cnf)
	name = strings.TrimSuffix(basename, filepath.Ext(basename))
	return
}

func loadConfig(configPath, configName string) (*Config, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	sc := make([]probe.Config, 0)
	err = viper.UnmarshalKey("Probes", &sc)

	adc := AutoDiscoverConfig{}
	err = viper.UnmarshalKey("autodiscover", &adc)

	ac := make([]alert.AlertConfig, 0)
	err = viper.UnmarshalKey("Alerts", &ac)

	conf := &Config{
		AutoDiscover: adc,
		Probes:       sc,
		Alerts:       ac,
	}
	err = viper.Unmarshal(conf)
	log.Printf("ProberConfig Read %v", conf)

	if err != nil {
		return nil, err
	}
	log.Printf(viper.ConfigFileUsed())

	return conf, nil

}

func setProbeConfigDefaults(c probe.ProberConfig) probe.ProberConfig {
	if c.Warning == 0 {
		c.Warning = 500 * time.Millisecond
	}

	if c.Fatal == 0 {
		c.Fatal = time.Minute
	}

	return c
}
