// Package httpdumper contains utilities to inspect HTTP requests.
package httpdumper

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

// EchoHandler echos requests back to the client as plain text.
type EchoHandler struct {
	log logrus.FieldLogger
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request, couldn't dump")
		h.log.WithError(err).Warn("failed to dump request")
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(rd); err != nil {
		h.log.WithError(err).Warn("failed to write to client")
	}
	h.log.WithFields(logrus.Fields{
		"protocol": r.Proto,
		"method":   r.Method,
		"url":      r.URL.String(),
		"from":     r.RemoteAddr,
	}).Info("request dumped")
}

var defaultEchoHandler = &EchoHandler{
	log: logrus.StandardLogger(),
}

// Echo the request back to the client as plain text.
func Echo(w http.ResponseWriter, r *http.Request) {
	defaultEchoHandler.ServeHTTP(w, r)
}
