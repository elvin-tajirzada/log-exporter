package containerization

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/elvin-tajirzada/log-exporter/internal/config"
	"io"
	"log"
	"time"
)

type (
	DockerAPI interface {
		GetLogs() (io.ReadCloser, error)
		ReConnect() io.ReadCloser
	}

	Docker struct {
		Client        *client.Client
		ContainerName string
	}
)

func NewDocker(conf *config.Docker) (*Docker, error) {
	cli, cliErr := client.NewClientWithOpts(client.WithHost(conf.SocketPath), client.WithVersion(conf.CliApiVersion))
	if cliErr != nil {
		return nil, fmt.Errorf("failed to create a new client: %v", cliErr)
	}

	return &Docker{
		Client:        cli,
		ContainerName: conf.ContainerName,
	}, nil
}

func (d *Docker) GetLogs() (io.ReadCloser, error) {
	// get container id
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

	// get container logs
	logs, logsErr := d.Client.ContainerLogs(context.Background(), containerID, options)
	if logsErr != nil {
		return nil, fmt.Errorf("failed to get container logs: %v", logsErr)
	}

	return logs, nil
}

func (d *Docker) ReConnect() io.ReadCloser {
	for {
		log.Println("reconnecting to container in 10 seconds...")
		time.Sleep(time.Second * 10)

		logs, logsErr := d.GetLogs()
		if logsErr == nil {
			log.Println("connection successful")
			return logs
		}

		log.Println(logsErr)
	}
}

func (d *Docker) getContainerID() (string, error) {
	var containerID string

	// get containers list
	containers, containersErr := d.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if containersErr != nil {
		return "", fmt.Errorf("failed to get container list: %v", containersErr)
	}

	// get container id
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
