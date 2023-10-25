package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/Lesterpig/board/config"
	"github.com/Lesterpig/board/manager"
	"github.com/sirupsen/logrus"
)

var port = flag.Int("p", 8080, "Port to use")
var configPath = flag.String("f", "./board.yaml", "Path to config file")
var kubeconfig = flag.String("kubeconfig", "", "Path to kubeconfig file")
var kubecontext = flag.String("kubecontext", "", "Context to use in kubeconfig file")
var namespace = flag.String("namespace", "", "Namespace to watch")

var log = logrus.StandardLogger()

func main() {
	flag.Parse()

	cfgDir, cfgName := config.ParseConfigString(*configPath)
	cfg, err := config.LoadConfig(cfgDir, cfgName, log)
	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}

	var kubeClient *manager.KubeClient
	if len(cfg.AutoDiscover) > 0 {
		if *kubeconfig != "" {
			cfg.KubeConfig.Kubeconfig = *kubeconfig
		}
		if *kubecontext != "" {
			cfg.KubeConfig.Kubecontext = *kubecontext
		}
		if *namespace != "" {
			cfg.KubeConfig.Namespace = *namespace
		}
		kubeClient, err = manager.NewKubeClient(&cfg.KubeConfig)
		if err != nil {
			log.Fatal("Error initializing Kubernetes client: ", err)
		}
	}

	mng, err := manager.NewManager(cfg, log, kubeClient)
	if err != nil {
		log.Fatal("Error initializing new mng: ", err)
	}

	// Setup static folder
	http.Handle("/", http.FileServer(rice.MustFindBox("static").HTTPBox()))

	// Setup logic route
	http.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(mng)
		_, _ = w.Write(data)
	})

	go mng.ProbeLoop(getLoopInterval(cfg.LoopInterval))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func getLoopInterval(duration time.Duration) time.Duration {
	if duration == 0 {
		log.Info("Using default interval of 10 minutes")
		return time.Minute * 10
	}
	log.Info("Using default interval of ", duration)
	return duration
}
