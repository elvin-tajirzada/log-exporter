package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elvin-tacirzade/log-exporter/internal/config"
	"github.com/elvin-tacirzade/log-exporter/internal/models"
	"io"
	"net/http"
	"time"
)

type (
	LokiAPI interface {
		Push(containerLog *models.ContainerLog) error
	}

	Loki struct {
		URL         string
		StreamKey   string
		StreamValue string
	}
)

func NewLoki(conf *config.Loki) *Loki {
	return &Loki{
		URL:         conf.URL,
		StreamKey:   conf.Stream.Key,
		StreamValue: conf.Stream.Value,
	}
}

func (l *Loki) Push(containerLog *models.ContainerLog) error {
	entries := [][]interface{}{
		{
			time.Now().UnixNano(),
			containerLog,
		},
	}

	// create payload
	payload := l.getPayload(entries)

	// encode payload
	payloadBytes, payloadBytesErr := json.Marshal(payload)
	if payloadBytesErr != nil {
		return fmt.Errorf("failed to convert payload to bytes: %v", payloadBytesErr)
	}

	// create a new request
	req, reqErr := http.NewRequest("POST", l.URL, bytes.NewBuffer(payloadBytes))
	if reqErr != nil {
		return fmt.Errorf("failed to create a new request: %v", reqErr)
	}

	// add header to request
	req.Header.Add("Content-Type", "application/json")

	// send HTTP request
	client := http.DefaultClient

	resp, respErr := client.Do(req)
	if respErr != nil {
		return fmt.Errorf("failed to send an HTTP request: %v", respErr)
	}

	defer resp.Body.Close()

	// check HTTP status code
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
