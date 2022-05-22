package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Application struct {
	RequestCount uint64
	SRegistry    *ServiceRegistry
}

func (a *Application) Handle(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&a.RequestCount, 1)

	if a.SRegistry.Len() == 0 {
		w.Write([]byte(`No backend entry in the service registry`))
		return
	}

	backendIndex := int(atomic.LoadUint64(&a.RequestCount) % uint64(a.SRegistry.Len()))
	fmt.Printf("Request routing to instance %d\n", backendIndex)

	a.SRegistry.GetByIndex(backendIndex).
		proxy.
		ServeHTTP(w, r)
}
