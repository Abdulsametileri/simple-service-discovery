package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	HelloServiceImageName = "hello"
	ContainerRunningState = "running"
	ContainerKillState    = "kill"
	ContainerStartState   = "start"
)

type DockerClient struct {
	*client.Client
}

func NewDockerClient() (*DockerClient, error) {
	dockerCLI, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &DockerClient{dockerCLI}, nil
}

func (dc *DockerClient) GetContainerPort(ctx context.Context, id string) (uint16, error) {
	containers, err := dc.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("id", id),
		),
	})

	if len(containers) == 1 {
		return containers[0].Ports[0].PublicPort, nil
	}
	return 0, err
}
