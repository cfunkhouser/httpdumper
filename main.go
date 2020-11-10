package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func dump(w http.ResponseWriter, r *http.Request) {
	rd, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request, couldn't dump")
		log.WithError(err).Warn("failed to dump request")
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(rd); err != nil {
		log.WithError(err).Warn("failed to write to client")
	}
	log.WithFields(log.Fields{
		"protocol": r.Proto,
		"method":   r.Method,
		"url":      r.URL.String(),
		"from":     r.RemoteAddr,
	}).Info("request dumped")
}

func main() {
	http.HandleFunc("/", dump)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
