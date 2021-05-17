package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pento/code-challenge/domain/tracker/models"

	pgerr "github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

var (
	ErrWrongVersion        = errors.New("wrong version")
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrTimeTrackerNotFound = errors.New("user not found")
)

type TrackerStore struct {
	pool *sql.DB
}

func NewTrackerStore(pool *sql.DB) *TrackerStore {
	return &TrackerStore{pool}
}

func (s TrackerStore) Get(ctx context.Context, id uint64) (models.TimeTracker, error) {

	row := s.pool.QueryRowContext(ctx, `
		SELECT id, started, ended, name, created_at, updated_at, deleted, version
		FROM time_tracker
		WHERE id = $1 AND deleted = 'f' 
	`, id)

	return s.scan(row)
}

func queryComposer(start, end time.Time) string {
	if start.IsZero() && end.IsZero() {
		return ""
	}

	return fmt.Sprintf("started between $1 AND $2 AND")
}

func (s TrackerStore) List(ctx context.Context, start, end time.Time) ([]models.TimeTracker, error) {

	var tracker []models.TimeTracker = make([]models.TimeTracker, 0)

	queryArgs := make([]interface{}, 0)

	if !start.IsZero() && !end.IsZero() {
		queryArgs = append(queryArgs, start)
		queryArgs = append(queryArgs, end)
	}

	arguments := queryComposer(start, end)

	rows, err := s.pool.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, started, ended, name, created_at, updated_at, deleted, version
		FROM time_tracker
		WHERE %s deleted = 'f'
		order by created_at ASC
	`, arguments), queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("%w failed to query context", err)
	}

	defer rows.Close()

	tracker, err = s.scanMultipleRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w error scan multiple rows", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w rows returned error", err)
	}

	return tracker, nil
}

func (s TrackerStore) Store(ctx context.Context, tracker models.TimeTracker, version uint32) (models.TimeTracker, error) {
	var result models.TimeTracker

	tx, err := s.pool.Begin()
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to begin transaction", err)
	}

	current, err := s.lockForUpdate(ctx, tx, tracker.ID)
	if err != nil {
		tx.Rollback()
		return models.TimeTracker{}, err
	}

	if current != version {
		tx.Rollback()
		return models.TimeTracker{}, ErrWrongVersion
	}

	if current == 0 {
		result, err = s.create(ctx, tx, tracker)
	} else {
		result, err = s.update(ctx, tx, tracker, version)
	}
	if err != nil {
		tx.Rollback()
		return models.TimeTracker{}, err
	}

	tx.Commit()

	if err != nil {
		return models.TimeTracker{}, err
	}

	return result, nil
}

func (s TrackerStore) lockForUpdate(ctx context.Context, tx *sql.Tx, id uint64) (uint32, error) {
	var version uint32

	row := tx.QueryRowContext(ctx, `
		SELECT version
		FROM time_tracker
		WHERE id = $1 FOR UPDATE NOWAIT
	`, id)

	err := row.Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return version, nil
}

func (s TrackerStore) Delete(ctx context.Context, id uint64) error {
	_, err := s.pool.ExecContext(ctx, `
		UPDATE time_tracker
		SET deleted = 't', updated_at = NOW()
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("%w failed to set to deleted", err)
	}

	return nil
}

func (s TrackerStore) create(ctx context.Context, tx *sql.Tx, tracker models.TimeTracker) (models.TimeTracker, error) {

	row := tx.QueryRowContext(ctx, `
		INSERT INTO time_tracker(started, name)
		VALUES ($1, $2)
		RETURNING id, started, ended, name, created_at, updated_at, deleted, version
	`,
		tracker.Start,
		tracker.Name,
	)
	return s.scan(row)
}

func (s TrackerStore) update(ctx context.Context, tx *sql.Tx, tracker models.TimeTracker, version uint32) (models.TimeTracker, error) {

	row := tx.QueryRowContext(ctx, `
		UPDATE time_tracker
		SET started = $1, ended = $2, name = $3, version = $4, updated_at = NOW()
		WHERE id = $5 AND version = $6
		RETURNING id, started, ended, name, created_at, updated_at, deleted, version
	`,
		tracker.Start,
		tracker.End,
		tracker.Name,
		version+1,
		tracker.ID,
		tracker.Meta.GetVersion(),
	)
	return s.scan(row)
}

func (s TrackerStore) scan(row *sql.Row) (models.TimeTracker, error) {
	var (
		id        uint64
		start     time.Time
		end       sql.NullTime
		name      string
		deleted   bool
		version   uint32
		createdAt time.Time
		updatedAt time.Time
	)

	if err := row.Scan(
		&id,
		&start,
		&end,
		&name, &createdAt, &updatedAt, &deleted, &version); err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == pgerr.UniqueViolation {
				return models.TimeTracker{}, ErrUniqueViolation
			}
		}

		if err == sql.ErrNoRows {
			return models.TimeTracker{}, ErrTimeTrackerNotFound
		}

		return models.TimeTracker{}, err
	}

	return s.hydrateTimeTracker(id, start, end, name, deleted, version, createdAt, updatedAt), nil
}

func (s TrackerStore) scanMultipleRows(rows *sql.Rows) ([]models.TimeTracker, error) {
	var (
		tracker []models.TimeTracker = make([]models.TimeTracker, 0)
	)

	type timeTracker struct {
		id        uint64
		start     time.Time
		end       sql.NullTime
		name      string
		deleted   bool
		version   uint32
		createdAt time.Time
		updatedAt time.Time
	}

	for rows.Next() {
		var timetracker timeTracker
		if err := rows.Scan(
			&timetracker.id,
			&timetracker.start,
			&timetracker.end,
			&timetracker.name, &timetracker.createdAt, &timetracker.updatedAt, &timetracker.deleted, &timetracker.version); err != nil {
			if pgErr, ok := err.(pgx.PgError); ok {
				if pgErr.Code == pgerr.UniqueViolation {
					return nil, ErrUniqueViolation
				}
			}

			if err == sql.ErrNoRows {
				return nil, ErrTimeTrackerNotFound
			}

			return nil, err
		}

		user := s.hydrateTimeTracker(timetracker.id, timetracker.start, timetracker.end,
			timetracker.name, timetracker.deleted, timetracker.version, timetracker.createdAt, timetracker.updatedAt)

		tracker = append(tracker, user)
	}

	return tracker, nil
}

func (s TrackerStore) hydrateTimeTracker(id uint64, start time.Time, end sql.NullTime,
	name string, deleted bool, version uint32, createdAt, updatedAt time.Time) models.TimeTracker {

	var tracker models.TimeTracker

	if end.Valid {
		tracker = models.NewTimeTracker(id, start, end.Time, name)
	} else {
		tracker = models.NewTimeTracker(id, start, time.Time{}, name)
	}

	tracker.Meta.HydrateMeta(deleted, createdAt, updatedAt, version)

	return tracker
}
