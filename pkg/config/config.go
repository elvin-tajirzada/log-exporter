package config

import (
	"fmt"
	"net/url"
	"os"
)

const (
	logFolderName        = "logs"
	logFileNameExtension = "json.log"
	streamKey            = "job"
	streamValue          = "log-exporter"
)

type (
	Config struct {
		Docker docker
		Loki   loki
	}

	docker struct {
		ContainerID, ContainerLogFilePath string
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
	containerID := os.Getenv("CONTAINER_ID")
	if containerID == "" {
		return nil, fmt.Errorf("environment parameter `CONTAINER_ID` can't be empty")
	}

	containerLogFilePath := fmt.Sprintf("/%s/%s-%s", logFolderName, containerID, logFileNameExtension)
	if _, err := os.Stat(containerLogFilePath); err != nil {
		return nil, fmt.Errorf("log file cannot be found. path: %s", containerLogFilePath)
	}

	lokiURL := os.Getenv("LOKI_URL")
	u, uErr := url.ParseRequestURI(lokiURL)
	if uErr != nil {
		return nil, fmt.Errorf("failed to parse LOKI_URL: %s", uErr)
	}

	return &Config{
		Docker: docker{
			ContainerID:          containerID,
			ContainerLogFilePath: containerLogFilePath,
		},
		Loki: loki{
			URL: u.String(),
			Stream: stream{
				Key:   streamKey,
				Value: streamValue,
			},
		},
	}, nil
}
