package main

import (
	"encoding/json"
	"flag"
	rice "github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

var port = flag.Int("p", 8080, "Port to use")
var intervalCli = flag.Int("i", 0, "Interval in minutes")
var configPath = flag.String("f", "./board.yaml", "Path to config file")

var log = logrus.StandardLogger()

func main() {
	flag.Parse()

	config, err := loadConfig(parseConfigString(*configPath))
	if err != nil {
		log.Fatal(os.Stderr, err)
	}
	manager, err := NewManager(config)

	if err != nil {
		log.Fatal(os.Stderr, err)
	}

	// Setup static folder
	http.Handle("/", http.FileServer(rice.MustFindBox("static").HTTPBox()))

	// Setup logic route
	http.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(manager)
		_, _ = w.Write(data)
	})

	go manager.ProbeLoop(getLoopInterval(config.LoopInterval))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func getLoopInterval(duration time.Duration) time.Duration {
	if duration == 0 {
		log.Info("Using default interval of 10 minutes")
		return time.Minute * 10
	}
	log.Info("Using default interval of %v", duration)
	return duration
}
