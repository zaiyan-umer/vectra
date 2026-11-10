package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

// Any type that has a Select() method satisfies BackendSelector interface.
// So each algorithm must implement Select() so it's Select() method can be used here
type BackendSelector interface {
	Select() int
}

type Proxy struct {
	selector BackendSelector
	backends []*backend.Backend
	proxies  []*httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index := p.selector.Select()

	if index < 0 || index >= len(p.proxies) {
		log.Printf("invalid backend index: %d", index)
		http.Error(w, "no backend is available", http.StatusServiceUnavailable)
		return
	}

	target := p.backends[index]
	target.Acquire()
	defer target.Release()

	// we need to verify that the backend is healthy
	p.proxies[index].ServeHTTP(w, r)
}

func NewProxy(selector BackendSelector, targetBackends []*backend.Backend) (*Proxy, error) {
	if selector == nil {
		return nil, fmt.Errorf("backend selector must not be nil")
	}

	backends := make([]*httputil.ReverseProxy, 0, len(targetBackends))

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,

		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		IdleConnTimeout:       90 * time.Second,
	}

	for _, rawTarget := range targetBackends {
		target, err := parseBackendURL(rawTarget.URL)
		if err != nil {
			return nil, err
		}

		reverseProxy := httputil.NewSingleHostReverseProxy(target)
		reverseProxy.Transport = transport
		reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			handleUnhealthyBackendResponse(w, r, err, rawTarget)
		}

		backends = append(backends, reverseProxy)
	}

	return &Proxy{selector: selector, backends: targetBackends, proxies: backends}, nil
}

func parseBackendURL(raw string) (*url.URL, error) {
	target, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse backend URL %q: %w", raw, err)
	}

	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, fmt.Errorf("backend URL %q must use http or https", raw)
	}

	if target.Host == "" {
		return nil, fmt.Errorf("backend URL %q must include a host", raw)
	}

	return target, nil
}

func handleUnhealthyBackendResponse(w http.ResponseWriter, r *http.Request, err error, rawTarget *backend.Backend) {
	rawTarget.SetHealthy(false)
	log.Printf("backend %s failed: %v", rawTarget.URL, err)
	http.Error(w, "backend unavailable", http.StatusBadGateway)
}
