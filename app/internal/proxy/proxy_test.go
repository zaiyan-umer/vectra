package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
