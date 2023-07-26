package config

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

const (
	dockerSocketPath    = "unix:///var/run/docker.sock"
	dockerCliApiVersion = "1.41"
	dockerReconnectTime = time.Second * 10
	lokiStreamKey       = "job"
	lokiStreamValue     = "log-exporter"
)

type (
	Config struct {
		Docker docker
		Loki   loki
	}

	docker struct {
		SocketPath, CliApiVersion, ContainerName string
		ReconnectTime                            time.Duration
	}

	loki struct {
		URL    string
		Stream stream
	}

	stream struct {
		Key   string
		Value string
	}
)

func New() (*Config, error) {
	containerName := os.Getenv("CONTAINER_NAME")
	if containerName == "" {
		return nil, fmt.Errorf("environment parameter `CONTAINER_NAME` can't be empty")
	}

	lokiURL := os.Getenv("LOKI_URL")
	u, uErr := url.ParseRequestURI(lokiURL)
	if uErr != nil {
		return nil, fmt.Errorf("failed to parse LOKI_URL: %s", uErr)
	}

	return &Config{
		Docker: docker{
			SocketPath:    dockerSocketPath,
			CliApiVersion: dockerCliApiVersion,
			ContainerName: containerName,
			ReconnectTime: dockerReconnectTime,
		},
		Loki: loki{
			URL: u.String(),
			Stream: stream{
				Key:   lokiStreamKey,
				Value: lokiStreamValue,
			},
		},
	}, nil
}
