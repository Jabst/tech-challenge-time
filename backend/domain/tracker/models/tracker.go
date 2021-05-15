package models

import (
	"pento/code-challenge/domain"
	"time"
)

type TimeTracker struct {
	ID    uint64
	Start time.Time
	End   time.Time
	Name  string
	Meta  domain.Meta
}

func NewTimeTracker(id uint64, start, end time.Time, name string) TimeTracker {
	return TimeTracker{
		ID:    id,
		Start: start,
		End:   end,
		Name:  name,
		Meta:  domain.NewMeta(),
	}
}

func (t TimeTracker) IsZero() bool {
	return t.ID == 0 &&
		t.Start.IsZero() &&
		t.End.IsZero() &&
		t.Name == ""
}
