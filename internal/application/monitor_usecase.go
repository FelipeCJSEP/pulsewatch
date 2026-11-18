package application

import (
	"context"

	"github.com/FelipeCJSEP/pulsewatch/internal/domain"
)

type CreateMonitorInput struct {
	Name             string
	TargetURL        string
	IntervalSeconds  int
	FailureThreshold int
}

type MonitorUseCase struct {
	monitorRepo  domain.MonitorRepository
	incidentRepo domain.IncidentRepository
}

func NewMonitorUseCase(monitorRepo domain.MonitorRepository, incidentRepo domain.IncidentRepository) *MonitorUseCase {
	return &MonitorUseCase{
		monitorRepo:  monitorRepo,
		incidentRepo: incidentRepo,
	}
}

func (uc *MonitorUseCase) CreateMonitor(ctx context.Context, input CreateMonitorInput) (*domain.Monitor, error) {
	monitor, err := domain.NewMonitor(
		input.Name,
		input.TargetURL,
		input.IntervalSeconds,
		input.FailureThreshold,
	)

	if err != nil {
		return nil, err
	}

	if err := uc.monitorRepo.Create(ctx, monitor); err != nil {
		return nil, err
	}

	return monitor, nil
}

func (uc *MonitorUseCase) ListMonitors(ctx context.Context) ([]*domain.Monitor, error) {
	return uc.monitorRepo.ListAll(ctx)
}

func (uc *MonitorUseCase) ListIncidentsByMonitor(ctx context.Context, monitorID string) ([]*domain.Incident, error) {
	return uc.incidentRepo.ListByMonitorID(ctx, monitorID)
}

func (uc *MonitorUseCase) ProcessCheckResult(ctx context.Context, monitorID string, isSuccess bool, failureReason string) error {
	monitor, err := uc.monitorRepo.GetByID(ctx, monitorID)

	if err != nil {
		return err
	}

	if isSuccess {
		wasDown := monitor.Status == domain.StatusDown
		monitor.RecordSuccess()

		if wasDown {
			incident, err := uc.incidentRepo.GetOpenByMonitorID(ctx, monitorID)

			if err == nil && incident != nil {
				_ = incident.Resolve()

				if err := uc.incidentRepo.Update(ctx, incident); err != nil {
					return err
				}
			}
		}
	} else {
		justFell := monitor.RecordFailed()

		if justFell {
			incident, err := domain.NewIncident(monitorID, failureReason)

			if err != nil {
				return err
			}

			if err := uc.incidentRepo.Create(ctx, incident); err != nil {
				return err
			}
		}
	}

	return uc.monitorRepo.Update(ctx, monitor)
}

func (uc *MonitorUseCase) AcknowledgeIncident(ctx context.Context, incidentID string) error {
	incident, err := uc.incidentRepo.GetByID(ctx, incidentID)

	if err != nil {
		return err
	}

	if err := incident.Acknowledge(); err != nil {
		return err
	}

	return uc.incidentRepo.Update(ctx, incident)
}
