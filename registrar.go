package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
)

type Registrar struct {
	Interval  time.Duration
	DockerCLI *client.Client
	SRegistry *ServiceRegistry
}

const (
	HelloServiceImageName = "hello"
	ContainerRunningState = "running"
)

func (r *Registrar) Observe() {
	for range time.Tick(r.Interval) {
		cList, _ := r.DockerCLI.ContainerList(context.Background(), types.ContainerListOptions{
			All: true,
		})

		if len(cList) == 0 {
			r.SRegistry.RemoveAll()
			continue
		}

		for _, c := range cList {
			if c.Image != HelloServiceImageName {
				continue
			}

			_, exist := r.SRegistry.GetByContainerID(c.ID)

			if c.State == ContainerRunningState {
				if !exist {
					addr := fmt.Sprintf("http://localhost:%d", c.Ports[0].PublicPort)
					r.SRegistry.Add(c.ID, addr)
				}
			} else {
				if exist {
					r.SRegistry.RemoveByContainerID(c.ID)
				}
			}
		}
	}
}
