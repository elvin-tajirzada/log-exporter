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
		Docker *Docker
		Loki   *Loki
	}

	Docker struct {
		SocketPath, CliApiVersion, ContainerName string
		ReconnectTime                            time.Duration
	}

	Loki struct {
		URL    string
		Stream *Stream
	}

	Stream struct {
		Key   string
		Value string
	}
)

func New() (*Config, error) {
	// get container name
	containerName := os.Getenv("CONTAINER_NAME")
	if containerName == "" {
		return nil, fmt.Errorf("environment parameter `CONTAINER_NAME` can't be empty")
	}

	// get loki URL
	lokiURL := os.Getenv("LOKI_URL")
	u, uErr := url.ParseRequestURI(lokiURL)
	if uErr != nil {
		return nil, fmt.Errorf("failed to parse LOKI_URL: %s", uErr)
	}

	return &Config{
		Docker: &Docker{
			SocketPath:    dockerSocketPath,
			CliApiVersion: dockerCliApiVersion,
			ContainerName: containerName,
			ReconnectTime: dockerReconnectTime,
		},
		Loki: &Loki{
			URL: u.String(),
			Stream: &Stream{
				Key:   lokiStreamKey,
				Value: lokiStreamValue,
			},
		},
	}, nil
}
