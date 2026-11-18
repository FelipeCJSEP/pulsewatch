package postgres

import (
	"context"
	"errors"

	"github.com/FelipeCJSEP/pulsewatch/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MonitorRepository struct {
	pool *pgxpool.Pool
}

func NewMonitorRepository(pool *pgxpool.Pool) *MonitorRepository {
	return &MonitorRepository{pool: pool}
}

func (mr *MonitorRepository) Create(ctx context.Context, m *domain.Monitor) error {
	query := `
		INSERT INTO monitors (id, name, target_url, interval_seconds, failure_threshold,
		consecutive_failures, status, last_checked_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := mr.pool.Exec(
		ctx,
		query,
		m.ID,
		m.Name,
		m.TargetURL,
		m.IntervalSeconds,
		m.FailureThreshold,
		m.ConsecutiveFailures,
		m.Status,
		m.LastCheckedAt,
		m.CreatedAt,
	)

	return err
}

func (mr *MonitorRepository) GetByID(ctx context.Context, id string) (*domain.Monitor, error) {
	query := `
		SELECT id, name, target_url, interval_seconds, failure_threshold,
			consecutive_failures, status, last_checked_at, created_at
		FROM monitors
		WHERE id = $1
	`

	var m domain.Monitor

	err := mr.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Name,
		&m.TargetURL,
		&m.IntervalSeconds,
		&m.FailureThreshold,
		&m.ConsecutiveFailures,
		&m.Status,
		&m.LastCheckedAt,
		&m.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &m, nil
}

func (mr *MonitorRepository) ListAll(ctx context.Context) ([]*domain.Monitor, error) {
	query := `
		SELECT id, name, target_url, interval_seconds, failure_threshold,
			consecutive_failures, status, last_checked_at, created_at
		FROM monitors
		ORDER BY created_at DESC
	`

	rows, err := mr.pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	monitors := []*domain.Monitor{}

	for rows.Next() {
		var m domain.Monitor

		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.TargetURL,
			&m.IntervalSeconds,
			&m.FailureThreshold,
			&m.ConsecutiveFailures,
			&m.Status,
			&m.LastCheckedAt,
			&m.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		monitors = append(monitors, &m)
	}

	return monitors, rows.Err()
}

func (mr *MonitorRepository) Update(ctx context.Context, m *domain.Monitor) error {
	query := `
		UPDATE monitors SET
			consecutive_failures = $1,
			status = $2,
			last_checked_at = $3
		WHERE id = $4
	`

	_, err := mr.pool.Exec(
		ctx,
		query,
		m.ConsecutiveFailures,
		m.Status,
		m.LastCheckedAt,
		m.ID,
	)

	return err
}
