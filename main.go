package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"
)

var port = flag.Int("p", 8080, "Port to use")

func main() {
	flag.Parse()

	// Setup static folder
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Setup logic route
	http.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(M)
		_, _ = w.Write(data)
	})

	go ProbeLoop(M, time.Minute)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
