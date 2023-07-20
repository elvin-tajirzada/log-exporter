package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elvin-tacirzade/log-exporter/pkg/config"
	"github.com/elvin-tacirzade/log-exporter/pkg/models"
	"io"
	"net/http"
	"time"
)

type Loki struct {
	URL         string
	StreamKey   string
	StreamValue string
}

func NewLoki(conf *config.Config) *Loki {
	return &Loki{
		URL:         conf.Loki.URL,
		StreamKey:   conf.Loki.Stream.Key,
		StreamValue: conf.Loki.Stream.Value,
	}
}

func (l *Loki) Push(containerLog *models.ContainerLog) error {
	entries := [][]interface{}{
		{
			time.Now().UnixNano(),
			containerLog,
		},
	}

	payload := l.getPayload(entries)

	payloadBytes, payloadBytesErr := json.Marshal(payload)
	if payloadBytesErr != nil {
		return fmt.Errorf("failed to convert payload to bytes: %v", payloadBytesErr)
	}

	req, reqErr := http.NewRequest("POST", l.URL, bytes.NewBuffer(payloadBytes))
	if reqErr != nil {
		return fmt.Errorf("failed to create a new request: %v", reqErr)
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient

	resp, respErr := client.Do(req)
	if respErr != nil {
		return fmt.Errorf("failed to send an HTTP request: %v", respErr)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return fmt.Errorf("failed to read HTTP response body: status code: %v, err: %v", resp.StatusCode, bodyErr)
		}
		return fmt.Errorf("failed to get correct HTTP response: status code: %v, err: %v", resp.StatusCode, string(body))
	}

	return nil
}

func (l *Loki) getPayload(entries [][]interface{}) map[string][]map[string]interface{} {
	return map[string][]map[string]interface{}{
		"streams": {
			{
				"stream": map[string]string{
					l.StreamKey: l.StreamValue,
				},
				"values": entries,
			},
		},
	}
}
