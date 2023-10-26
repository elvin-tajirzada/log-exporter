package service

import (
	"bufio"
	"encoding/json"
	"github.com/elvin-tajirzada/log-exporter/internal/models"
	"github.com/elvin-tajirzada/log-exporter/pkg/containerization"
	"github.com/elvin-tajirzada/log-exporter/pkg/db"
	"github.com/elvin-tajirzada/log-exporter/pkg/extraction"
	"log"
)

type LogExporter struct {
	Loki          db.LokiAPI
	Docker        containerization.DockerAPI
	ContainerName string
}

func NewLogExporter(loki db.LokiAPI, docker containerization.DockerAPI, containerName string) *LogExporter {
	return &LogExporter{
		Loki:          loki,
		Docker:        docker,
		ContainerName: containerName,
	}
}

func (l *LogExporter) Start() {
	containerLog := models.ContainerLog{
		Container: l.ContainerName,
	}

	// get docker container logs
	logs, logsErr := l.Docker.GetLogs()
	if logsErr != nil {
		log.Fatal(logsErr)
	}

	defer logs.Close()

	for {
		// create a new reader for logs
		reader := bufio.NewReader(logs)

		for {
			// get log line
			line, lineErr := reader.ReadString('\n')
			if lineErr != nil {
				log.Printf("failed to read log: %v\n", lineErr)
				break
			}

			// extraction json value from log line
			jsonLog, jsonLogErr := extraction.JSON(line)
			if jsonLogErr != nil {
				log.Println(jsonLogErr)
				continue
			}

			// decode json log
			unmarshalErr := json.Unmarshal([]byte(jsonLog), &containerLog)
			if unmarshalErr != nil {
				log.Printf("failed to unmarshall container log: log: %s, err: %v\n", jsonLog, unmarshalErr)
				continue
			}

			// insert log to loki
			lokiPushErr := l.Loki.Push(&containerLog)
			if lokiPushErr != nil {
				log.Fatalf("failed to push to loki: %v", lokiPushErr)
			}
		}

		// reconnect docker
		logs = l.Docker.ReConnect()
	}
}
