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
	Store(ctx context.Context, tracker models.TimeTracker, version uint32) (models.TimeTracker, error)
	Delete(ctx context.Context, id uint64) error
}

type TrackerService struct {
	store TrackerStore
}

type CreateTrackerParams struct {
	Start time.Time
	End   time.Time
	Name  string
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

func (s TrackerService) GetTimeTracker(ctx context.Context, id uint64) (models.TimeTracker, error) {
	user, err := s.store.Get(ctx, id)
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to get user", err)
	}

	if user.IsZero() {
		return models.TimeTracker{}, ErrTrackerNotFound
	}

	return user, nil
}

func (s TrackerService) CreateUser(ctx context.Context, params CreateTrackerParams) (models.TimeTracker, error) {
	timeTracker := models.NewTimeTracker(0, params.Start, params.End, params.Name)

	timeTracker, err := s.store.Store(ctx, timeTracker, 0)
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to store time tracker", err)
	}

	return timeTracker, nil
}

func (s TrackerService) UpdateUser(ctx context.Context, params UpdateTrackerParams) (models.TimeTracker, error) {
	timeTracker, err := s.GetTimeTracker(ctx, params.ID)
	if err != nil {
		return models.TimeTracker{}, err
	}

	if timeTracker.IsZero() {
		return models.TimeTracker{}, ErrTrackerNotFound
	}

	timeTracker, err = s.store.Store(ctx, timeTracker, params.Version)
	if err != nil {
		return models.TimeTracker{}, fmt.Errorf("%w failed to store user", err)
	}

	return timeTracker, nil
}

func (s TrackerService) DeleteUser(ctx context.Context, params DeleteTrackerParams) error {
	err := s.store.Delete(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("%w failed to delete user", err)
	}

	return nil
}
