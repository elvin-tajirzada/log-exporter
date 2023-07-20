package reader

import (
	"encoding/json"
	"fmt"
	"github.com/elvin-tacirzade/log-exporter/pkg/db"
	"github.com/elvin-tacirzade/log-exporter/pkg/models"
	"github.com/hpcloud/tail"
	"log"
)

type Reader struct {
	Tail *tail.Tail
	Loki *db.Loki
}

func New(filePath string, loki *db.Loki) (*Reader, error) {
	t, err := tail.TailFile(filePath, tail.Config{
		Follow: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begins tailing the file. file path: %v, err: %v", filePath, err)
	}

	return &Reader{
		Tail: t,
		Loki: loki,
	}, nil
}

func (r *Reader) Handle() {
	var (
		dockerLog    models.DockerLog
		containerLog models.ContainerLog
	)

	for line := range r.Tail.Lines {

		unMarshalDockerLogErr := json.Unmarshal([]byte(line.Text), &dockerLog)
		if unMarshalDockerLogErr != nil {
			log.Printf("failed to unmarshall docker log: log: %s, err: %v", line.Text, unMarshalDockerLogErr)
			continue
		}

		unMarshalContainerLogErr := json.Unmarshal([]byte(dockerLog.Log), &containerLog)
		if unMarshalContainerLogErr != nil {
			log.Printf("failed to unmarshall container log: log: %s, err: %v", dockerLog.Log, unMarshalContainerLogErr)
			continue
		}

		lokiPushErr := r.Loki.Push(&containerLog)
		if lokiPushErr != nil {
			log.Fatalf("failed to push to loki: %v", lokiPushErr)
		}

	}
}
