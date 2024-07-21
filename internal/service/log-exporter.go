package service

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/elvin-tajirzada/log-exporter/pkg/containerization"
	"github.com/elvin-tajirzada/log-exporter/pkg/database"
	"github.com/elvin-tajirzada/log-exporter/pkg/extraction"
)

type LogExporter struct {
	Loki             database.Client
	Docker           containerization.Client
	ContainerName    string
	ReconnectionTime time.Duration
}

func NewLogExporter(
	loki database.Client,
	docker containerization.Client,
	containerName string,
	reconnectTime time.Duration,
) *LogExporter {
	return &LogExporter{
		Loki:             loki,
		Docker:           docker,
		ContainerName:    containerName,
		ReconnectionTime: reconnectTime,
	}
}

func (l *LogExporter) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Log exporter service shutdown successfully")
			return
		default:
			l.handleLogs(ctx)
		}
	}
}

func (l *LogExporter) handleLogs(ctx context.Context) {
	var containerLog map[string]interface{}

	// get docker container logs
	logs, err := l.Docker.GetLogs(ctx)
	if err != nil {
		log.Printf("Unable to get logs: %v\n", err)
		// reconnect docker client
		log.Printf("Trying to reconnect docker container in %s\n", l.ReconnectionTime.String())
		time.Sleep(l.ReconnectionTime)
		return
	}
	log.Println("Successfully connected to Docker container")
	defer logs.Close()

	// create a new reader for logs
	reader := bufio.NewReader(logs)

	for {
		// get log line
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Unable to read log: %v\n", err)
			time.Sleep(time.Second * 1)
			break
		}

		// extraction json value from log line
		jsonLog, err := extraction.JSON(line)
		if err != nil {
			log.Printf("Unable to extraction json from log line: %v\n", err)
			continue
		}

		if err = json.Unmarshal([]byte(jsonLog), &containerLog); err != nil {
			log.Printf("Unable to decode log: %v\n", err)
			continue
		}

		containerLog["container"] = l.ContainerName
		containerLogByte, err := json.Marshal(containerLog)
		if err != nil {
			log.Printf("Unable to encode log: %v\n", err)
			continue
		}

		// insert log to loki
		if err = l.Loki.Push(containerLogByte); err != nil {
			log.Fatalf("Unable to push to loki: %v\n", err)
		}
	}
}
