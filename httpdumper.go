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
	Log logrus.FieldLogger
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request, couldn't dump")
		h.Log.WithError(err).Warn("failed to dump request")
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(rd); err != nil {
		h.Log.WithError(err).Warn("failed to write to client")
	}
	h.Log.WithFields(logrus.Fields{
		"protocol": r.Proto,
		"method":   r.Method,
		"url":      r.URL.String(),
		"from":     r.RemoteAddr,
	}).Info("request dumped")
}

var defaultEchoHandler = &EchoHandler{
	Log: logrus.StandardLogger(),
}

// Echo the request back to the client as plain text.
func Echo(w http.ResponseWriter, r *http.Request) {
	defaultEchoHandler.ServeHTTP(w, r)
}

// LoggingTransport logs HTTP requests and responses through a http.Client.
type LoggingTransport struct {
	Transport http.RoundTripper
	Log       logrus.FieldLogger
}

func (t *LoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.Log.WithFields(logrus.Fields{
		"protocol": r.Proto,
		"method":   r.Method,
		"url":      r.URL.String(),
	}).Info("outgoing request")

	if rd, err := httputil.DumpRequestOut(r, true); err == nil {
		t.Log.Debug(rd)
	} else {
		t.Log.WithError(err).Warn("failed to dump outgoing request")
	}

	resp, err := t.Transport.RoundTrip(r)

	t.Log.WithFields(logrus.Fields{
		"code":     resp.StatusCode,
		"status":   resp.Status,
		"protocol": r.Proto,
		"method":   resp.Request.Method,
		"url":      resp.Request.URL,
	}).Info("incoming response")

	if rd, err := httputil.DumpResponse(resp, true); err == nil {
		t.Log.Debug(rd)
	} else {
		t.Log.WithError(err).Warn("failed to dump incoming response")
	}

	return resp, err
}

// DefaultTransport is a LoggingTransport with the defaults for underlying
// transport and logger.
func DefaultTransport() http.RoundTripper {
	return &LoggingTransport{
		Transport: http.DefaultTransport,
		Log:       logrus.StandardLogger(),
	}
}
