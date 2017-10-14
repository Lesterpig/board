package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
)

var port = flag.Int("p", 8080, "Port to use")

func main() {
	flag.Parse()
	manager, err := loadConfig()
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

	go manager.ProbeLoop(10 * time.Minute)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
