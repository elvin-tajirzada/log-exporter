package containerization

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/elvin-tajirzada/log-exporter/internal/config"
)

type Docker struct {
	Client        *dockerClient.Client
	containerName string
	reconnectTime time.Duration
}

func NewDocker(conf *config.Docker) (*Docker, error) {
	client, err := dockerClient.NewClientWithOpts(dockerClient.WithHost(conf.SocketPath), dockerClient.WithVersion(conf.CliApiVersion))
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}

	return &Docker{
		Client:        client,
		containerName: conf.ContainerName,
		reconnectTime: conf.ReconnectTime,
	}, nil
}

func (d *Docker) GetLogs(ctx context.Context) (io.ReadCloser, error) {
	// get container id
	containerID, err := d.getContainerID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get container id: %v", err)
	}

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Tail:       "0",
	}

	// get container logs
	logs, err := d.Client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (d *Docker) getContainerID(ctx context.Context) (string, error) {
	var containerID string

	// get containers list
	containers, err := d.Client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return "", fmt.Errorf("unable to list container: %v", err)
	}

	// get container id
	for _, container := range containers {
		containerName := fmt.Sprintf("/%s", d.containerName)
		if container.Names[0] == containerName {
			containerID = container.ID
			break
		}
	}

	if containerID == "" {
		return "", fmt.Errorf("container not found: %v", d.containerName)
	}

	return containerID, nil
}
