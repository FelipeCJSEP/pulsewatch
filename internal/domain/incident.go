package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IncidentStatus string

const (
	IncidentOpen         IncidentStatus = "OPEN"
	IncidentAcknowledged IncidentStatus = "ACKNOWLEDGED"
	IncidentResolved     IncidentStatus = "RESOLVED"
)

var (
	ErrIncidentMonitorRequired = errors.New("o ID do monitor é obrigatório")
	ErrIncidentReasonRequired  = errors.New("o motivo da falha é obrigatório")
	ErrIncidentAlreadyResolved = errors.New("o incidente já está resolvido e não pode ser alterado")
)

type Incident struct {
	ID         string         `json:"id"`
	MonitorID  string         `json:"monitorId"`
	Reason     string         `json:"reason"`
	Status     IncidentStatus `json:"status"`
	OpenedAt   time.Time      `json:"openedAt"`
	ResolvedAt *time.Time     `json:"resolvedAt"`
}

func NewIncident(monitorID string, reason string) (*Incident, error) {
	monitorID = strings.TrimSpace(monitorID)
	reason = strings.TrimSpace(reason)

	if monitorID == "" {
		return nil, ErrIncidentMonitorRequired
	}

	if reason == "" {
		return nil, ErrIncidentReasonRequired
	}

	return &Incident{
		ID:        uuid.NewString(),
		MonitorID: monitorID,
		Reason:    reason,
		Status:    IncidentOpen,
		OpenedAt:  time.Now().UTC(),
	}, nil
}

func (i *Incident) Acknowledge() error {
	if i.Status == IncidentResolved {
		return ErrIncidentAlreadyResolved
	}

	i.Status = IncidentAcknowledged
	return nil
}

func (i *Incident) Resolve() error {
	if i.Status == IncidentResolved {
		return ErrIncidentAlreadyResolved
	}

	i.Status = IncidentResolved
	now := time.Now().UTC()
	i.ResolvedAt = &now

	return nil
}
