package services

import (
	"context"
	"errors"
	"fmt"
	"pento/code-challenge/domain/tracker/models"

	"time"
)

var (
	ErrTrackerNotFound = errors.New("tracker not found")
	ErrWrongVersion    = errors.New("wrong version provided")
)

type TrackerStore interface {
	Get(ctx context.Context, id uint64) (models.TimeTracker, error)
	List(ctx context.Context, start, end time.Time) ([]models.TimeTracker, error)
	Store(ctx context.Context, tracker models.TimeTracker, version uint32) (models.TimeTracker, error)
	Delete(ctx context.Context, id uint64) error
}

type TrackerService struct {
	store TrackerStore
}

type CreateTrackerParams struct {
	Start time.Time
	Name  string
}

type ListTimeTracker struct {
	Start time.Time
	End   time.Time
}

type UpdateTrackerParams struct {
	ID      uint64
	Start   time.Time
	End     time.Time
	Name    string
	Version uint32
}

type DeleteTrackerParams struct {
	ID uint64
}

func NewTrackerService(store TrackerStore) TrackerService {
	return TrackerService{
		store: store,
	}
}

func (s TrackerService) GetTracker(ctx context.Context, id uint64) (models.TimeTracker, error) {
	timeTracker, err := s.store.Get(ctx, id)
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to get tracker", err)
	}

	if timeTracker.IsZero() {
		return models.TimeTracker{}, ErrTrackerNotFound
	}

	return timeTracker, nil
}

func (s TrackerService) ListTrackers(ctx context.Context, params ListTimeTracker) ([]models.TimeTracker, error) {

	timeTrackers, err := s.store.List(ctx, params.Start, params.End)
	if err != nil {
		return nil, fmt.Errorf("%w failed to list trackers", err)
	}

	return timeTrackers, nil
}

func (s TrackerService) CreateTracker(ctx context.Context, params CreateTrackerParams) (models.TimeTracker, error) {
	timeTracker := models.NewTimeTracker(0, params.Start, time.Time{}, params.Name)

	timeTracker, err := s.store.Store(ctx, timeTracker, 0)
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to store time tracker", err)
	}

	return timeTracker, nil
}

func (s TrackerService) UpdateTracker(ctx context.Context, params UpdateTrackerParams) (models.TimeTracker, error) {
	timeTracker, err := s.GetTracker(ctx, params.ID)
	if err != nil {
		return models.TimeTracker{}, err
	}

	if timeTracker.IsZero() {
		return models.TimeTracker{}, ErrTrackerNotFound
	}

	if !params.End.IsZero() {
		timeTracker.End = params.End
	}

	if params.Name != "" {
		timeTracker.Name = params.Name
	}

	timeTracker, err = s.store.Store(ctx, timeTracker, params.Version)
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to store tracker", err)
	}

	return timeTracker, nil
}

func (s TrackerService) DeleteTracker(ctx context.Context, params DeleteTrackerParams) error {
	err := s.store.Delete(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("%w failed to delete user", err)
	}

	return nil
}
