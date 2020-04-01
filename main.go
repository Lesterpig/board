package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gobuffalo/envy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var port = flag.Int("p", 8080, "Port to use")
var intervalCli = flag.Int("i", 0, "Interval in minutes")
var configPath = flag.String("f", "./board.yaml", "Path to config file")

var log = logrus.StandardLogger()

func main() {
	flag.Parse()

	interval := getInterval()

	log.Infof("Probe interval: %d", interval)

	manager, err := loadConfig(parseConfigString(*configPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Setup static folder
	http.Handle("/", http.FileServer(rice.MustFindBox("static").HTTPBox()))

	// Setup logic route
	http.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(manager)
		_, _ = w.Write(data)
	})

	http.Handle("/metrics", promhttp.Handler())

	go manager.ProbeLoop(time.Duration(int64(interval)) * time.Minute)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func getInterval() int {
	intervalEnv := getInt(envy.Get("INTERVAL", ""))

	if *intervalCli != 0 {
		return *intervalCli
	} else if intervalEnv != 0 {
		return intervalEnv
	}

	return 10
}

func getInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if nil != err {
		return 0
	}

	return int(i)
}
