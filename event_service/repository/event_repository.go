package repository

import (
	"context"
	"database/sql"
	"time"

	"muchway/event_service/domain"
)

type EventRepository interface {
	Create(ctx context.Context, e *domain.Event) (*domain.Event, error)
	Get(ctx context.Context, id string) (*domain.Event, error)
	Update(ctx context.Context, e *domain.Event) (*domain.Event, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Event, error)
}

type pgEventRepo struct{ db *sql.DB }

func NewPostgresEventRepository(db *sql.DB) EventRepository {
	return &pgEventRepo{db: db}
}

func (r *pgEventRepo) Create(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	now := time.Now()
	e.CreatedAt, e.UpdatedAt = now, now
	const stmt = `
    INSERT INTO events (name, start_time, status, winner_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, created_at, updated_at
    `
	err := r.db.QueryRowContext(ctx, stmt,
		e.Name, e.StartTime, e.Status, e.WinnerID, e.CreatedAt, e.UpdatedAt,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *pgEventRepo) Get(ctx context.Context, id string) (*domain.Event, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,name,start_time,status,winner_id,created_at,updated_at
     FROM events WHERE id=$1`, id,
	)
	var e domain.Event
	if err := row.Scan(&e.ID, &e.Name, &e.StartTime, &e.Status, &e.WinnerID, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *pgEventRepo) Update(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	e.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE events SET name=$2,start_time=$3,status=$4,winner_id=$5,updated_at=$6 WHERE id=$1`,
		e.ID, e.Name, e.StartTime, e.Status, e.WinnerID, e.UpdatedAt,
	)
	return e, err
}

func (r *pgEventRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id=$1`, id)
	return err
}

func (r *pgEventRepo) List(ctx context.Context) ([]*domain.Event, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,name,start_time,status,winner_id,created_at,updated_at FROM events`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.Name, &e.StartTime, &e.Status, &e.WinnerID, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, nil
}
