package postgres

import (
	"context"
	"errors"

	"github.com/FelipeCJSEP/pulsewatch/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository struct {
	pool *pgxpool.Pool
}

func NewIncidentRepository(pool *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{pool: pool}
}

func (ir *IncidentRepository) Create(ctx context.Context, i *domain.Incident) error {
	query := `
		INSERT INTO incidents (id, monitor_id, reason, status, opened_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := ir.pool.Exec(
		ctx,
		query,
		i.ID,
		i.MonitorID,
		i.Reason,
		i.Status,
		i.OpenedAt,
		i.ResolvedAt,
	)

	return err
}

func (ir *IncidentRepository) get(ctx context.Context, query string, args ...any) (*domain.Incident, error) {
	var i domain.Incident

	err := ir.pool.QueryRow(ctx, query, args...).Scan(
		&i.ID,
		&i.MonitorID,
		&i.Reason,
		&i.Status,
		&i.OpenedAt,
		&i.ResolvedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &i, nil
}

func (ir *IncidentRepository) GetByID(ctx context.Context, id string) (*domain.Incident, error) {
	query := `
		SELECT id, monitor_id, reason, status, opened_at, resolved_at
		FROM incidents
		WHERE id = $1
	`

	return ir.get(ctx, query, id)
}

func (ir *IncidentRepository) GetOpenByMonitorID(ctx context.Context, monitorID string) (*domain.Incident, error) {
	query := `
		SELECT id, monitor_id, reason, status, opened_at, resolved_at
		FROM incidents
		WHERE monitor_id = $1 AND status != 'RESOLVED'
		ORDER BY opened_at DESC
		LIMIT 1
	`

	return ir.get(ctx, query, monitorID)
}

func (ir *IncidentRepository) ListByMonitorID(ctx context.Context, monitorID string) ([]*domain.Incident, error) {
	query := `
		SELECT id, monitor_id, reason, status, opened_at, resolved_at
		FROM incidents
		WHERE monitor_id = $1
		ORDER BY opened_at DESC
	`

	rows, err := ir.pool.Query(ctx, query, monitorID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	incidents := []*domain.Incident{}

	for rows.Next() {
		var i domain.Incident

		err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.Reason,
			&i.Status,
			&i.OpenedAt,
			&i.ResolvedAt,
		)

		if err != nil {
			return nil, err
		}

		incidents = append(incidents, &i)
	}

	return incidents, rows.Err()
}

func (ir *IncidentRepository) Update(ctx context.Context, i *domain.Incident) error {
	query := `
		UPDATE incidents SET
			status = $1,
			resolved_at = $2
		WHERE id = $3
	`

	_, err := ir.pool.Exec(
		ctx,
		query,
		i.Status,
		i.ResolvedAt,
		i.ID,
	)

	return err
}
