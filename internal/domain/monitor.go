package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MonitorStatus string

const (
	StatusPending MonitorStatus = "PENDING"
	StatusUp      MonitorStatus = "UP"
	StatusDown    MonitorStatus = "DOWN"
	StatusPaused  MonitorStatus = "PAUSED"
)

var (
	ErrMonitorNameRequired    = errors.New("o nome do monitor é obrigatório")
	ErrMonitorInvalidURL      = errors.New("a URL deve começar com http:// ou https://")
	ErrMonitorInvalidInterval = errors.New("o intervalo mínimo é de 5 segundos")
)

type Monitor struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	TargetURL           string        `json:"targetUrl"`
	IntervalSeconds     int           `json:"intervalSeconds"`
	FailureThreshold    int           `json:"failureThreshold"`
	ConsecutiveFailures int           `json:"consecutiveFailures"`
	Status              MonitorStatus `json:"status"`
	LastCheckedAt       *time.Time    `json:"lastCheckedAt"`
	CreatedAt           time.Time     `json:"createdAt"`
}

func NewMonitor(name string, targetURL string, intervalSeconds int, failureThreshold int) (*Monitor, error) {
	name = strings.TrimSpace(name)
	targetURL = strings.TrimSpace(targetURL)

	if name == "" {
		return nil, ErrMonitorNameRequired
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		return nil, ErrMonitorInvalidURL
	}

	if intervalSeconds < 5 {
		return nil, ErrMonitorInvalidInterval
	}

	if failureThreshold <= 0 {
		failureThreshold = 3
	}

	return &Monitor{
		ID:               uuid.NewString(),
		Name:             name,
		TargetURL:        targetURL,
		IntervalSeconds:  intervalSeconds,
		FailureThreshold: failureThreshold,
		Status:           StatusPending,
		CreatedAt:        time.Now().UTC(),
	}, nil
}

func (m *Monitor) RecordSuccess() {
	now := time.Now().UTC()
	m.LastCheckedAt = &now
	m.ConsecutiveFailures = 0
	m.Status = StatusUp
}

func (m *Monitor) RecordFailed() bool {
	now := time.Now().UTC()
	m.LastCheckedAt = &now
	m.ConsecutiveFailures++

	if m.ConsecutiveFailures >= m.FailureThreshold && m.Status != StatusDown {
		m.Status = StatusDown
		return true
	}

	return false
}

func (m *Monitor) Pause() {
	m.Status = StatusPaused
}
