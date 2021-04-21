package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/cfunkhouser/httpdumper"
)

var address = flag.String("address", ":8080", "Bind address in host:port format.")

func main() {
	flag.Parse()
	if *address == "" {
		log.Fatal("Can't bind an empty address")
	}

	http.HandleFunc("/", httpdumper.Echo)
	http.Handle("/metrics", promhttp.Handler())
	log.WithField("address", *address).Info("Server starting")
	if err := http.ListenAndServe(*address, nil); err != nil {
		log.WithError(err).Error("Server has exited")
	}
}
