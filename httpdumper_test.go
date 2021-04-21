package httpdumper

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func discardLogger(tb testing.TB) logrus.FieldLogger {
	tb.Helper()
	l := logrus.New()
	l.Out = io.Discard
	return l
}

func bodyWriter(tb testing.TB, body string) io.ReadCloser {
	var buf bytes.Buffer
	if _, err := (&buf).WriteString(body); err != nil {
		tb.Fatalf("error setting up test: %v", err)
	}
	return ioutil.NopCloser(&buf)
}

func TestEchoHandlerServeHTTP(t *testing.T) {
	r := httptest.NewRequest("GET", "https://example.com:9443/api", bodyWriter(t, "Hello there!"))
	r.Header.Add("X-Some-Thing", "testing")

	w := httptest.NewRecorder()

	(&EchoHandler{
		Log: discardLogger(t),
	}).ServeHTTP(w, r)

	resp := w.Result()

	if got := resp.StatusCode; got != 200 {
		t.Errorf("ServeHTTP(): StatusCode mismatch: got: %v, want: 200", got)
	}

	wantCTH := "text/plain; charset=utf-8"
	if got := resp.Header.Get("Content-Type"); got != wantCTH {
		t.Errorf("ServeHTTP(): Content-Type Header mismatch: got: %q, want: %q", got, wantCTH)
	}

	got, _ := io.ReadAll(resp.Body)

	want := "GET https://example.com:9443/api HTTP/1.1\r\nX-Some-Thing: testing\r\n\r\nHello there!"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("ServeHTTP(): response body mismatch (-want +got):\n%v", diff)
	}
}
