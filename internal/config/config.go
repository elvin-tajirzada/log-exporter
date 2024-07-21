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
	containerName := os.Getenv("CONTAINER_NAME")
	if containerName == "" {
		return nil, fmt.Errorf("environment parameter `CONTAINER_NAME` can't be empty")
	}

	lokiURL := os.Getenv("LOKI_URL")
	if lokiURL == "" {
		return nil, fmt.Errorf("environment parameter `LOKI_URL` can't be empty")
	}

	u, err := url.ParseRequestURI(lokiURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse LOKI_URL: %s", err)
	}

	reconnTime := dockerReconnectTime
	reconnectionTime := os.Getenv("DOCKER_RECONNECTION_TIME")
	if reconnectionTime != "" {
		reconnTime, err = time.ParseDuration(reconnectionTime)
		if err != nil {
			return nil, fmt.Errorf("unable to parse DOCKER_RECONNECTION_TIME: %s", err)
		}
	}

	return &Config{
		Docker: &Docker{
			SocketPath:    dockerSocketPath,
			CliApiVersion: dockerCliApiVersion,
			ContainerName: containerName,
			ReconnectTime: reconnTime,
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
