package reader

import (
	"bufio"
	"encoding/json"
	"github.com/elvin-tacirzade/log-exporter/pkg/db"
	"github.com/elvin-tacirzade/log-exporter/pkg/docker"
	"github.com/elvin-tacirzade/log-exporter/pkg/models"
	"github.com/elvin-tacirzade/log-exporter/pkg/utils"
	"log"
)

type Reader struct {
	Loki   *db.Loki
	Docker *docker.Docker
}

func New(loki *db.Loki, dock *docker.Docker) (*Reader, error) {
	return &Reader{
		Loki:   loki,
		Docker: dock,
	}, nil
}

func (r *Reader) Handle() {
	containerLog := models.ContainerLog{
		Container: r.Docker.ContainerName,
	}

	logs, logsErr := r.Docker.GetLogs()
	if logsErr != nil {
		log.Fatal(logsErr)
	}

	defer logs.Close()

	for {
		reader := bufio.NewReader(logs)

		for {
			line, lineErr := reader.ReadString('\n')
			if lineErr != nil {
				log.Printf("failed to read log: %v\n", lineErr)
				break
			}

			jsonLog, jsonLogErr := utils.ExtractJSON(line)
			if jsonLogErr != nil {
				log.Println(jsonLogErr)
				continue
			}

			unmarshalErr := json.Unmarshal([]byte(jsonLog), &containerLog)
			if unmarshalErr != nil {
				log.Printf("failed to unmarshall container log: log: %s, err: %v\n", jsonLog, unmarshalErr)
				continue
			}

			lokiPushErr := r.Loki.Push(&containerLog)
			if lokiPushErr != nil {
				log.Fatalf("failed to push to loki: %v", lokiPushErr)
			}
		}

		logs = r.Docker.ReConnect()
	}
}
