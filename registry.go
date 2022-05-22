package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type backend struct {
	proxy       *httputil.ReverseProxy
	containerID string
}

type ServiceRegistry struct {
	mu       sync.RWMutex
	backends []backend
}

func (s *ServiceRegistry) Init() {
	s.mu = sync.RWMutex{}
	s.backends = []backend{}
}

func (s *ServiceRegistry) Add(containerID, addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	URL, _ := url.Parse(addr)

	s.backends = append(s.backends, backend{
		proxy:       httputil.NewSingleHostReverseProxy(URL),
		containerID: containerID,
	})
}

func (s *ServiceRegistry) GetByContainerID(containerID string) (backend, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, b := range s.backends {
		if b.containerID == containerID {
			return b, true
		}
	}

	return backend{}, false
}

func (s *ServiceRegistry) GetByIndex(index int) backend {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.backends[index]
}

func (s *ServiceRegistry) RemoveByContainerID(containerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var backends []backend
	for _, b := range s.backends {
		if b.containerID == containerID {
			continue
		}
		backends = append(backends, b)
	}

	s.backends = backends
}

func (s *ServiceRegistry) RemoveAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.backends = []backend{}
}

func (s *ServiceRegistry) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.backends)
}

func (s *ServiceRegistry) List() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.backends {
		fmt.Println(s.backends[i].containerID)
	}
}
