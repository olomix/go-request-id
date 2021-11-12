package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	go_request_id "github.com/olomix/go-request-id"
)

func TestInstrumentClient(t *testing.T) {
	client := &http.Client{}
	InstrumentClient(client)

	var m sync.Mutex
	var rid string
	var expectedRId = "XXX-111"

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			rid = r.Header.Get("X-Request-Id")
			m.Unlock()
		}),
	)
	defer srv.Close()

	ctx := go_request_id.WithRequestID(context.Background(), expectedRId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	m.Lock()
	defer m.Unlock()
	if rid != expectedRId {
		t.Fatalf("expected request-id '%v', got '%v'", expectedRId, rid)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	client := &http.Client{}
	InstrumentClient(client)

	var m sync.Mutex
	var rid string
	var expectedRId = "XXX-111"

	fn := func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		rid = go_request_id.ExtractRequestID(r.Context())
		m.Unlock()
	}
	handler := RequestIDMiddleware(http.HandlerFunc(fn))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Request-Id", expectedRId)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	m.Lock()
	defer m.Unlock()
	if rid != expectedRId {
		t.Fatalf("expected request-id '%v', got '%v'", expectedRId, rid)
	}
}
