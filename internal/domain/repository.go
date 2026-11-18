package domain

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("registro não encontrado")

type MonitorRepository interface {
	Create(ctx context.Context, monitor *Monitor) error
	GetByID(ctx context.Context, id string) (*Monitor, error)
	ListAll(ctx context.Context) ([]*Monitor, error)
	Update(ctx context.Context, monitor *Monitor) error
}

type IncidentRepository interface {
	Create(ctx context.Context, incident *Incident) error
	GetByID(ctx context.Context, id string) (*Incident, error)
	GetOpenByMonitorID(ctx context.Context, monitorID string) (*Incident, error)
	ListByMonitorID(ctx context.Context, monitorID string) ([]*Incident, error)
	Update(ctx context.Context, incident *Incident) error
}
