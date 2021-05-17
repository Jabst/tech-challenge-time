// +build integrationdb

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"pento/code-challenge/domain"
	"pento/code-challenge/domain/tracker/models"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/tkuchiki/faketime"
)

var ERROR = fmt.Errorf("expected error")

func initTrackerStore() (*TrackerStore, error) {

	connString := fmt.Sprintf("host=localhost port=5434 user=postgres password=postgres dbname=postgres sslmode=disable")

	pool, err := sql.Open("pgx", connString)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(`delete from time_tracker;
		ALTER SEQUENCE time_tracker_id_seq RESTART WITH 1;
		INSERT INTO time_tracker(started, ended, name, created_at, updated_at, version)
		VALUES ('2020-05-15 00:00:00', '2020-05-15 10:00:00', 'test_time_tracker_1', '2020-01-01 00:00:01', '2020-01-01 00:00:00', 1),
			('2020-05-16 00:00:00', '2020-05-16 10:00:00', 'test_time_tracker_2', '2020-02-01 00:00:01', '2020-01-01 00:00:00', 1);
		`)
	if err != nil {
		panic(err)
	}

	store := NewTrackerStore(pool)

	return store, nil
}

func Test_TrackerStore_Store(t *testing.T) {

	sampleMeta := domain.NewMeta()
	sampleMeta2 := domain.NewMeta()

	sampleMeta.HydrateMeta(false, time.Now(), time.Now(), 1)
	sampleMeta2.HydrateMeta(false, time.Now(), time.Now(), 2)

	f := faketime.NewFaketime(2021, time.May, 1, 1, 0, 0, 0, time.UTC)
	defer f.Undo()
	f.Do()

	type testInput struct {
		tracker models.TimeTracker
	}

	type testExpectation struct {
		err    error
		result models.TimeTracker
	}

	testCases := []struct {
		description string
		input       testInput
		expected    testExpectation
	}{
		{
			description: "when creating a time tracker",
			input: testInput{
				tracker: models.TimeTracker{
					Name:  "test_tracker_3",
					Start: time.Now(),
					Meta:  domain.NewMeta(),
				},
			},
			expected: testExpectation{
				result: models.TimeTracker{
					Name:  "test_tracker_3",
					Start: time.Now(),
					End:   time.Time{},
					ID:    3,
					Meta:  sampleMeta,
				},
				err: nil,
			},
		},
		{
			description: "when updating a time tracker",
			input: testInput{
				tracker: models.TimeTracker{
					ID:    1,
					Name:  "test_tracker_3",
					Start: time.Now(),
					End:   time.Now().Add(time.Duration(1) * time.Hour),
					Meta:  sampleMeta,
				},
			},
			expected: testExpectation{
				result: models.TimeTracker{
					Name:  "test_tracker_3",
					Start: time.Now(),
					End:   time.Now().Add(time.Duration(1) * time.Hour),
					ID:    1,
					Meta:  sampleMeta2,
				},
				err: nil,
			},
		},
		{
			description: "when updating a time tracker but version is wrong",
			input: testInput{
				tracker: models.TimeTracker{
					ID:    1,
					Name:  "test_tracker_3",
					Start: time.Now(),
					End:   time.Now().Add(time.Duration(1) * time.Hour),
					Meta:  sampleMeta2,
				},
			},
			expected: testExpectation{
				result: models.TimeTracker{},
				err:    ErrWrongVersion,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			var ctx = context.TODO()
			defer ctx.Done()

			repo, err := initTrackerStore()
			defer repo.pool.Close()
			g.Expect(err).ToNot(HaveOccurred(), "should not return an error setting up the repository")

			result, err := repo.Store(ctx, tc.input.tracker, tc.input.tracker.Meta.GetVersion())

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				g.Expect(result.Start).To(Equal(tc.expected.result.Start), "should be the same start timestamp")
				g.Expect(result.End).To(Equal(tc.expected.result.End), "should be the same end timestamp")
				g.Expect(result.Name).To(Equal(tc.expected.result.Name), "should be the same name")
				g.Expect(result.ID).To(Equal(tc.expected.result.ID), "should be the same id")
				g.Expect(result.Meta.GetVersion()).To(Equal(tc.expected.result.Meta.GetVersion()), "should be the same version")
			}
		})
	}
}

func Test_TrackerStore_Delete(t *testing.T) {

	type testInput struct {
		id uint64
	}

	type testExpectation struct {
		err    error
		result models.TimeTracker
	}

	expectedMeta := domain.NewMeta()
	expectedMeta.SetDeleted(true)

	testCases := []struct {
		description string
		input       testInput
		expected    testExpectation
	}{
		{
			description: "when deleting a user",
			input: testInput{
				id: 1,
			},
			expected: testExpectation{
				result: models.TimeTracker{
					Meta: expectedMeta,
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			var ctx = context.TODO()
			defer ctx.Done()

			repo, err := initTrackerStore()
			defer repo.pool.Close()
			g.Expect(err).ToNot(HaveOccurred(), "should not return an error setting up the repository")

			err = repo.Delete(ctx, tc.input.id)

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
			}
		})
	}
}

func Test_TrackerStore_Get(t *testing.T) {

	sampleMeta := domain.NewMeta()
	sampleMeta2 := domain.NewMeta()

	sampleMeta.HydrateMeta(false, time.Now(), time.Now(), 1)
	sampleMeta2.HydrateMeta(false, time.Now(), time.Now(), 2)

	f := faketime.NewFaketime(2021, time.May, 1, 1, 0, 0, 0, time.UTC)
	defer f.Undo()
	f.Do()

	type testInput struct {
		id uint64
	}

	type testExpectation struct {
		err    error
		result models.TimeTracker
	}

	testCases := []struct {
		description string
		input       testInput
		expected    testExpectation
	}{
		{
			description: "when getting a tracker",
			input: testInput{
				id: 1,
			},
			expected: testExpectation{
				result: models.TimeTracker{
					Start: time.Date(2020, time.May, 15, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2020, time.May, 15, 10, 0, 0, 0, time.UTC),
					Name:  "test_time_tracker_1",
					Meta:  sampleMeta,
					ID:    1,
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			var ctx = context.TODO()
			defer ctx.Done()
			// '2020-05-15 00:00:00', '2020-05-15 10:00:00', 'test_time_tracker_1'

			repo, err := initTrackerStore()
			defer repo.pool.Close()
			g.Expect(err).ToNot(HaveOccurred(), "should not return an error setting up the repository")

			result, err := repo.Get(ctx, tc.input.id)

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				g.Expect(result.Start).To(Equal(tc.expected.result.Start), "should be the same start date")
				g.Expect(result.End).To(Equal(tc.expected.result.End), "should be the same end date")
				g.Expect(result.Name).To(Equal(tc.expected.result.Name), "should be the same name")
				g.Expect(result.Meta.GetVersion()).To(Equal(tc.expected.result.Meta.GetVersion()), "should be the same version")
				g.Expect(result.Meta.GetDeleted()).To(Equal(tc.expected.result.Meta.GetDeleted()), "should be the same deleted status")
				g.Expect(result.ID).To(Equal(tc.expected.result.ID), "should be the same id")
			}
		})
	}
}

func Test_TrackerStore_List(t *testing.T) {

	sampleMeta := domain.NewMeta()

	sampleMeta.HydrateMeta(false, time.Now(), time.Now(), 1)

	f := faketime.NewFaketime(2021, time.May, 1, 1, 0, 0, 0, time.UTC)
	defer f.Undo()
	f.Do()

	type testInput struct {
		start time.Time
		end   time.Time
	}

	type testExpectation struct {
		err    error
		result []models.TimeTracker
	}

	testCases := []struct {
		description string
		input       testInput
		expected    testExpectation
	}{
		{
			description: "when listing all time trackers",
			input: testInput{
				start: time.Time{},
				end:   time.Time{},
			},
			expected: testExpectation{
				result: []models.TimeTracker{
					{
						Start: time.Date(2020, time.May, 15, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2020, time.May, 15, 10, 0, 0, 0, time.UTC),
						Name:  "test_time_tracker_1",
						Meta:  sampleMeta,
						ID:    1,
					},
					{
						Start: time.Date(2020, time.May, 16, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2020, time.May, 16, 10, 0, 0, 0, time.UTC),
						Name:  "test_time_tracker_2",
						Meta:  sampleMeta,
						ID:    2,
					},
				},
				err: nil,
			},
		},
		{
			description: "when listing with time window",
			input: testInput{
				start: time.Date(2020, time.May, 15, 0, 0, 1, 0, time.UTC),
				end:   time.Date(2020, time.May, 16, 10, 0, 0, 1, time.UTC),
			},
			expected: testExpectation{
				result: []models.TimeTracker{
					{
						Start: time.Date(2020, time.May, 16, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2020, time.May, 16, 10, 0, 0, 0, time.UTC),
						Name:  "test_time_tracker_2",
						Meta:  sampleMeta,
						ID:    2,
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			var ctx = context.TODO()
			defer ctx.Done()

			repo, err := initTrackerStore()
			defer repo.pool.Close()
			g.Expect(err).ToNot(HaveOccurred(), "should not return an error setting up the repository")

			result, err := repo.List(ctx, tc.input.start, tc.input.end)

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				for index := range result {
					g.Expect(result[index].Start).To(Equal(tc.expected.result[index].Start), "should be the same start date")
					g.Expect(result[index].End).To(Equal(tc.expected.result[index].End), "should be the same end date")
					g.Expect(result[index].Name).To(Equal(tc.expected.result[index].Name), "should be the same name")
					g.Expect(result[index].Meta.GetVersion()).To(Equal(tc.expected.result[index].Meta.GetVersion()), "should be the same version")
					g.Expect(result[index].Meta.GetDeleted()).To(Equal(tc.expected.result[index].Meta.GetDeleted()), "should be the same deleted status")
					g.Expect(result[index].ID).To(Equal(tc.expected.result[index].ID), "should be the same id")
				}
			}
		})
	}
}
