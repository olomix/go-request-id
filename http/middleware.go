package http

import (
	"net/http"

	go_request_id "github.com/olomix/go-request-id"
)

// RequestIDMiddleware extracts request-id from incoming headers and attaches
// it to context.
func RequestIDMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := rIDFromRequest(r)
		if rid != "" {
			r = r.WithContext(go_request_id.WithRequestID(r.Context(), rid))
		}
		handler.ServeHTTP(w, r)
	})
}

// InstrumentClient replaces Transport in *http.Client to attach X-Request-Id
// header with request-id information to outgoing requests. If you need custom
// transport, put it into client before instrumenting, as overloading it after
// will remove instrumentation
func InstrumentClient(client *http.Client) {
	tr := client.Transport
	if tr == nil {
		tr = http.DefaultTransport
	}
	client.Transport = ridTransport{base: tr}
}

// ridTransport sets the X-Request-Id header before calling base.
type ridTransport struct {
	base http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t ridTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rid := go_request_id.ExtractRequestID(req.Context())
	if rid != "" {
		req.Header.Set("X-Request-Id", rid)
	}
	return t.base.RoundTrip(req)
}

// CloseIdleConnections closes any connections on its Transport which
// were previously connected from previous requests but are now
// sitting idle in a "keep-alive" state. It does not interrupt any
// connections currently in use.
//
// If the base Transport does not have a CloseIdleConnections method
// then this method does nothing.
func (t ridTransport) CloseIdleConnections() {
	type closeIdler interface {
		CloseIdleConnections()
	}
	if tr, ok := t.base.(closeIdler); ok {
		tr.CloseIdleConnections()
	}
}

// Looks like Go normalizes X-Request-ID to X-Request-Id. Should check both
// to be sure we found request-id.
var rIDHeaders = []string{"X-Request-ID", "X-Request-Id"}

func rIDFromRequest(r *http.Request) string {
	for _, header := range rIDHeaders {
		if h, ok := r.Header[header]; ok && len(h) > 0 {
			return h[0]
		}
	}
	return ""
}
