package main

import (
	"github.com/docker/docker/client"
	"log"
	"net/http"
	"time"
)

func main() {
	registry := ServiceRegistry{}
	registry.Init()

	dockerCLI, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	registrar := Registrar{
		Interval:  3 * time.Second,
		SRegistry: &registry,
		DockerCLI: dockerCLI,
	}
	go registrar.Observe()

	app := Application{SRegistry: &registry}
	http.HandleFunc("/", app.Handle)

	log.Fatalln(http.ListenAndServe(":3000", nil))
}
