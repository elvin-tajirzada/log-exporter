package models

import (
	"net"
	"time"
)

type DockerLog struct {
	Log    string    `json:"log"`
	Stream string    `json:"stream"`
	Time   time.Time `json:"time"`
}

type ContainerLog struct {
	IP      net.IP    `json:"ip" db:"ip"`
	Caller  string    `json:"caller" db:"caller"`
	Path    string    `json:"path" db:"path"`
	Level   string    `json:"level" db:"level"`
	Method  string    `json:"method" db:"method"`
	Status  int       `json:"status" db:"status"`
	Message string    `json:"msg" db:"msg"`
	Device  string    `json:"dt" db:"dt"`
	Timing  float64   `json:"timing" db:"timing"`
	Time    time.Time `json:"ts" db:"ts"`
}
