package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

type fixedSelector int

func (s fixedSelector) Select() int {
	return int(s)
}

func TestNewProxyRejectsInvalidBackendURL(t *testing.T) {
	_, err := NewProxy(fixedSelector(0), []*backend.Backend{{URL: "ftp://backend"}})
	if err == nil {
		t.Fatal("NewProxy returned nil error for unsupported URL scheme")
	}
}

func TestProxyReturnsUnavailableWhenSelectorIsInvalid(t *testing.T) {
	proxy, err := NewProxy(fixedSelector(-1), nil)
	if err != nil {
		t.Fatalf("NewProxy returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://load-balancer/", nil)

	proxy.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
}

func TestProxyForwardsRequestToSelectedBackend(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" {
			t.Errorf("backend path = %q, want /hello", r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("backend response"))
	}))
	defer backendServer.Close()

	proxy, err := NewProxy(fixedSelector(0), []*backend.Backend{{URL: backendServer.URL}})
	if err != nil {
		t.Fatalf("NewProxy returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://load-balancer/hello", nil)
	proxy.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}
	if !strings.Contains(recorder.Body.String(), "backend response") {
		t.Fatalf("body = %q, want backend response", recorder.Body.String())
	}
}

func TestNewProxyConfiguresBackendResponseTimeout(t *testing.T) {
	backendServer := httptest.NewServer(http.NotFoundHandler())
	defer backendServer.Close()

	proxy, err := NewProxy(fixedSelector(0), []*backend.Backend{{URL: backendServer.URL}})
	if err != nil {
		t.Fatalf("NewProxy returned error: %v", err)
	}

	transport, ok := proxy.proxies[0].Transport.(*http.Transport)
	if !ok {
		t.Fatalf("proxy transport type = %T, want *http.Transport", proxy.proxies[0].Transport)
	}
	if transport.ResponseHeaderTimeout <= 0 {
		t.Fatal("proxy response-header timeout is not configured")
	}
}

func TestProxyReturnsWhenBackendStalls(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer backendServer.Close()

	proxy, err := NewProxy(fixedSelector(0), []*backend.Backend{{URL: backendServer.URL}})
	if err != nil {
		t.Fatalf("NewProxy returned error: %v", err)
	}

	start := time.Now()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://load-balancer/", nil)
	proxy.ServeHTTP(recorder, request)

	if elapsed := time.Since(start); elapsed > 6*time.Second {
		t.Fatalf("stalled backend request took %s, want at most 6s", elapsed)
	}
	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadGateway)
	}
}

func TestProxyCancelsBackendRequestWhenClientCancels(t *testing.T) {
	backendStarted := make(chan struct{})
	backendCanceled := make(chan struct{})
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(backendStarted)
		<-r.Context().Done()
		close(backendCanceled)
	}))
	defer backendServer.Close()

	proxy, err := NewProxy(fixedSelector(0), []*backend.Backend{{URL: backendServer.URL}})
	if err != nil {
		t.Fatalf("NewProxy returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	request := httptest.NewRequest(http.MethodGet, "http://load-balancer/", nil).WithContext(ctx)
	done := make(chan struct{})

	go func() {
		proxy.ServeHTTP(httptest.NewRecorder(), request)
		close(done)
	}()

	select {
	case <-backendStarted:
	case <-time.After(time.Second):
		t.Fatal("backend request did not start")
	}

	cancel()

	select {
	case <-backendCanceled:
	case <-time.After(time.Second):
		t.Fatal("backend request was not canceled")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("proxy did not return after client cancellation")
	}
}
