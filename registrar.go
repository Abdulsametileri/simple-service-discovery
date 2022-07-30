package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type Registrar struct {
	DockerClient *DockerClient
	SRegistry    *ServiceRegistry
}

func (r *Registrar) Init() error {
	cList, err := r.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("ancestor", HelloServiceImageName),
			filters.Arg("status", ContainerRunningState),
		),
	})
	if err != nil {
		return err
	}

	for _, c := range cList {
		r.SRegistry.Add(c.ID, findContainerAddress(c.Ports[0].PublicPort))
	}

	return nil
}

func (r *Registrar) Observe() {
	msgCh, errCh := r.DockerClient.Events(context.Background(), types.EventsOptions{
		Filters: filters.NewArgs(
			filters.Arg("type", "container"),
			filters.Arg("image", HelloServiceImageName),
			filters.Arg("event", "start"),
			filters.Arg("event", "kill"),
		),
	})

	for {
		select {
		case c := <-msgCh:
			fmt.Printf("State of the container %d is %s\n", c.ID, c.Status)
			if c.Status == ContainerKillState {
				r.SRegistry.RemoveByContainerID(c.ID)
			} else if c.Status == ContainerStartState {
				port, err := r.DockerClient.GetContainerPort(context.Background(), c.ID)
				if err != nil {
					fmt.Printf("err getting newly started container port %s\n", err.Error())
					continue
				}
				r.SRegistry.Add(c.ID, findContainerAddress(port))
			}
		case err := <-errCh:
			fmt.Println("Error Docker Event Chan", err.Error())
		}
	}
}

func findContainerAddress(cPort uint16) string {
	return fmt.Sprintf("http://localhost:%d", cPort)
}
