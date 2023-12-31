package models

import (
	"net"
	"time"
)

type ContainerLog struct {
	IP        net.IP    `json:"ip" db:"ip"`
	Container string    `json:"container" db:"container"`
	Caller    string    `json:"caller" db:"caller"`
	Path      string    `json:"path" db:"path"`
	Level     string    `json:"level" db:"level"`
	Method    string    `json:"method" db:"method"`
	Status    int       `json:"status" db:"status"`
	Message   string    `json:"msg" db:"msg"`
	Device    string    `json:"dt" db:"dt"`
	Timing    float64   `json:"timing" db:"timing"`
	Time      time.Time `json:"ts" db:"ts"`
}
