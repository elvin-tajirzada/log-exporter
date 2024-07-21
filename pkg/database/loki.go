package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/elvin-tajirzada/log-exporter/internal/config"
)

type Loki struct {
	url, streamKey, streamValue string
}

func NewLoki(conf *config.Loki) *Loki {
	return &Loki{
		url:         conf.URL,
		streamKey:   conf.Stream.Key,
		streamValue: conf.Stream.Value,
	}
}

func (l *Loki) Push(log []byte) error {
	entries := [][]string{
		{
			fmt.Sprintf("%d", time.Now().UnixNano()),
			string(log),
		},
	}

	payload := l.getPayload(entries)

	// encode payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to encode payload: %v", err)
	}

	// TODO: write request package
	// create a new request
	req, err := http.NewRequest("POST", l.url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("unable to create new request: %v", err)
	}

	// add header to request
	req.Header.Add("Content-Type", "application/json")

	// send HTTP request
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send HTTP request: %v", err)
	}

	defer resp.Body.Close()

	// check HTTP status code
	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read HTTP response body: status code: %v, err: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("unable to get correct HTTP response: status code: %v, err: %v", resp.StatusCode, string(body))
	}

	return nil
}

func (l *Loki) getPayload(entries [][]string) map[string][]map[string]interface{} {
	return map[string][]map[string]interface{}{
		"streams": {
			{
				"stream": map[string]string{
					l.streamKey: l.streamValue,
				},
				"values": entries,
			},
		},
	}
}
