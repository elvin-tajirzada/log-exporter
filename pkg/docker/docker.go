package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/elvin-tacirzade/log-exporter/pkg/config"
	"io"
	"log"
	"time"
)

type Docker struct {
	Client        *client.Client
	ContainerName string
}

func New(conf *config.Config) (*Docker, error) {
	cli, cliErr := client.NewClientWithOpts(client.WithHost(conf.Docker.SocketPath), client.WithVersion(conf.Docker.CliApiVersion))
	if cliErr != nil {
		return nil, fmt.Errorf("failed to create a new client: %v", cliErr)
	}

	return &Docker{
		Client:        cli,
		ContainerName: conf.Docker.ContainerName,
	}, nil
}

func (d *Docker) getContainerID() (string, error) {
	var containerID string

	containers, containersErr := d.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if containersErr != nil {
		return "", fmt.Errorf("failed to get container list: %v", containersErr)
	}

	for _, container := range containers {
		if container.Names[0] == "/"+d.ContainerName {
			containerID = container.ID
			break
		}
	}

	if containerID == "" {
		return "", fmt.Errorf("container with name '%s' not found", d.ContainerName)
	}

	return containerID, nil
}

func (d *Docker) GetLogs() (io.ReadCloser, error) {
	containerID, containerIDErr := d.getContainerID()
	if containerIDErr != nil {
		return nil, fmt.Errorf("failed to get container id: %v", containerIDErr)
	}

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Tail:       "0",
	}

	logs, logsErr := d.Client.ContainerLogs(context.Background(), containerID, options)
	if logsErr != nil {
		return nil, fmt.Errorf("failed to get container logs: %v", logsErr)
	}

	return logs, nil
}

func (d *Docker) ReConnect() io.ReadCloser {
	for {
		log.Println("Reconnecting to container in 10 seconds...")
		time.Sleep(time.Second * 10)

		logs, logsErr := d.GetLogs()
		if logsErr == nil {
			log.Println("connection successful")
			return logs
		}

		log.Println(logsErr)
	}
}
