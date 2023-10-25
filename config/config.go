package config

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/Lesterpig/board/alert"
	"github.com/sirupsen/logrus"

	"github.com/Lesterpig/board/probe"

	"github.com/spf13/viper"
)

type Config struct {
	LoopInterval time.Duration
	AutoDiscover []AutoDiscoverConfig
	Probes       []probe.Config
	Alerts       []alert.AlertConfig
}

type AutoDiscoverConfig struct {
	LoopInterval       time.Duration
	KubernetesResource string
}

func ParseConfigString(cnf string) (dir string, name string) {
	dir = filepath.Dir(cnf)
	basename := filepath.Base(cnf)
	name = strings.TrimSuffix(basename, filepath.Ext(basename))
	return
}

func LoadConfig(configPath, configName string, log *logrus.Logger) (*Config, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	sc := make([]probe.Config, 0)
	err = viper.UnmarshalKey("Probes", &sc)
	if err != nil {
		return nil, err
	}

	adc := make([]AutoDiscoverConfig, 0)
	err = viper.UnmarshalKey("AutoDiscover", &adc)
	if err != nil {
		return nil, err
	}

	ac := make([]alert.AlertConfig, 0)
	err = viper.UnmarshalKey("Alerts", &ac)
	if err != nil {
		return nil, err
	}

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

func SetProbeConfigDefaults(c probe.ProberConfig) probe.ProberConfig {
	if c.Warning == 0 {
		c.Warning = 500 * time.Millisecond
	}

	if c.Fatal == 0 {
		c.Fatal = time.Minute
	}

	return c
}
