package application

import (
	"context"
	"testing"

	"github.com/FelipeCJSEP/pulsewatch/internal/domain"
)

type fakeMonitorRepo struct {
	monitors map[string]*domain.Monitor
}

func (f *fakeMonitorRepo) Create(ctx context.Context, m *domain.Monitor) error {
	f.monitors[m.ID] = m
	return nil
}

func (f *fakeMonitorRepo) GetByID(ctx context.Context, id string) (*domain.Monitor, error) {
	m, ok := f.monitors[id]

	if !ok {
		return nil, domain.ErrNotFound
	}

	return m, nil
}

func (f *fakeMonitorRepo) ListAll(ctx context.Context) ([]*domain.Monitor, error) {
	var list []*domain.Monitor

	for _, m := range f.monitors {
		list = append(list, m)
	}

	return list, nil
}

func (f *fakeMonitorRepo) Update(ctx context.Context, m *domain.Monitor) error {
	f.monitors[m.ID] = m
	return nil
}

type fakeIncidentRepo struct {
	incidents []*domain.Incident
}

func (f *fakeIncidentRepo) Create(ctx context.Context, i *domain.Incident) error {
	f.incidents = append(f.incidents, i)
	return nil
}

func (f *fakeIncidentRepo) GetByID(ctx context.Context, id string) (*domain.Incident, error) {
	for _, i := range f.incidents {
		if i.ID == id {
			return i, nil
		}
	}

	return nil, domain.ErrNotFound
}

func (f *fakeIncidentRepo) GetOpenByMonitorID(ctx context.Context, monitorID string) (*domain.Incident, error) {
	for _, i := range f.incidents {
		if i.MonitorID == monitorID && i.Status != domain.IncidentResolved {
			return i, nil
		}
	}

	return nil, domain.ErrNotFound
}

func (f *fakeIncidentRepo) ListByMonitorID(ctx context.Context, monitorID string) ([]*domain.Incident, error) {
	return f.incidents, nil
}

func (f *fakeIncidentRepo) Update(ctx context.Context, i *domain.Incident) error {
	return nil
}

func TestProcessCheckResult_OpenAndResolvesIncident(t *testing.T) {
	ctx := context.Background()
	monitorRepo := &fakeMonitorRepo{monitors: make(map[string]*domain.Monitor)}
	incidentRepo := &fakeIncidentRepo{}
	useCase := NewMonitorUseCase(monitorRepo, incidentRepo)

	monitor, err := useCase.CreateMonitor(ctx, CreateMonitorInput{
		Name:             "API Teste",
		TargetURL:        "https://api.teste.com",
		IntervalSeconds:  10,
		FailureThreshold: 2,
	})

	if err != nil {
		t.Fatalf("erro inesperado ao criar monitor: %v", err)
	}

	// simula duas falhas consecutivas (deve cair e abrir 1 incidente)
	_ = useCase.ProcessCheckResult(ctx, monitor.ID, false, "HTTP 500")
	_ = useCase.ProcessCheckResult(ctx, monitor.ID, false, "HTTP 500")

	if len(incidentRepo.incidents) != 1 {
		t.Fatalf("esperava 1 incidente aberto, mas encontrou %d", len(incidentRepo.incidents))
	}

	// Simula um sucesso (Auto-Recovery: deve resolver o incidente aberto)
	_ = useCase.ProcessCheckResult(ctx, monitor.ID, true, "")

	if incidentRepo.incidents[0].Status != domain.IncidentResolved {
		t.Fatalf("esperava incidente RESOLVED após recuperação, mas está %s", incidentRepo.incidents[0].Status)
	}
}
