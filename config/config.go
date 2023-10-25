package config

import (
	"os"
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
	KubeConfig   KubeClientConfig
}

type KubeClientConfig struct {
	Kubeconfig  string
	Kubecontext string
	Timeout     time.Duration
	Namespace   string
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

	loopInterval := viper.GetDuration("LoopInterval")

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
	for i := range adc {
		if adc[i].LoopInterval == 0 {
			log.Infof("No loop interval set for %v, using default of 5 minutes", adc[i].KubernetesResource)
			adc[i].LoopInterval = 5 * time.Minute
		}
	}

	ac := make([]alert.AlertConfig, 0)
	err = viper.UnmarshalKey("Alerts", &ac)
	if err != nil {
		return nil, err
	}
	kubeconf := KubeClientConfig{}
	err = viper.UnmarshalKey("KubeClient", &kubeconf)
	if err != nil {
		return nil, err
	}
	ns := os.Getenv("NAMESPACE")
	if ns != "" {
		kubeconf.Namespace = ns
	} else {
		log.Info("No namespace set, using default")
		kubeconf.Namespace = "default"
	}

	conf := &Config{
		LoopInterval: loopInterval,
		AutoDiscover: adc,
		Probes:       sc,
		Alerts:       ac,
		KubeConfig:   kubeconf,
	}
	err = viper.Unmarshal(conf)
	log.Info("ProberConfig Read: ", conf)

	if err != nil {
		return nil, err
	}
	log.Info(viper.ConfigFileUsed())

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
