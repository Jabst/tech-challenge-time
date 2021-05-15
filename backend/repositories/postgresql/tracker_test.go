// +build integrationdb

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	_ "github.com/jackc/pgx/stdlib"
)

var ERROR = fmt.Errorf("expected error")

var medias = []models.TimeTracker{
	{
		Country:   "uk",
		Email:     "example@example.qqq",
		FirstName: "Test",
		LastName:  "Test",
		Nickname:  "testuser",
		Password:  "qwerty",
		Meta:      domain.NewMeta(),
	},
}

func initTrackerStore() (*TrackerStore, error) {

	connString := fmt.Sprintf("host=localhost port=5434 user=postgres password=postgres dbname=postgres sslmode=disable")

	pool, err := sql.Open("pgx", connString)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(`delete from users;
		ALTER SEQUENCE users_id_seq RESTART WITH 1;
		INSERT INTO users(start, end, name, created_at, updated_at, version)
		VALUES ('2020-01-01 00:00:00', '2020-01-01 10:00:00', 'test_time_tracker_1', '2020-01-01 00:00:00', '2020-01-01 00:00:00', 1),
			('2020-01-01 00:00:00', '2020-01-01 10:00:00', 'test_time_tracker_2', '2020-01-01 00:00:00', '2020-01-01 00:00:00', 1);
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

	sampleMeta.HydrateMeta(1, time.Now(), time.Now(), false)
	sampleMeta2.HydrateMeta(2, time.Now(), time.Now(), false)

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
					ID:    1,
					Start: time.Time,

					Meta: domain.NewMeta(),
				},
			},
			expected: testExpectation{
				result: models.TimeTracker{
					Country:   "uk",
					Email:     "user1@qweqwe.com",
					FirstName: "test",
					LastName:  "test",
					Meta:      sampleMeta,
					Nickname:  "testUser1",
					Password:  "aue8r9gau98e",
					ID:        3,
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

			result, err := repo.Store(ctx, tc.input.user, tc.input.user.Meta.GetVersion())

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				g.Expect(result.Country).To(Equal(tc.expected.result.Country), "should be the same country")
				g.Expect(result.Nickname).To(Equal(tc.expected.result.Nickname), "should be the same nickname")
				g.Expect(result.FirstName).To(Equal(tc.expected.result.FirstName), "should be the same first name")
				g.Expect(result.LastName).To(Equal(tc.expected.result.LastName), "should be the same last name")
				g.Expect(result.Email).To(Equal(tc.expected.result.Email), "should be the same email")
				g.Expect(result.Password).To(Equal(tc.expected.result.Password), "should be the same password")
				g.Expect(result.ID).To(Equal(tc.expected.result.ID), "should be the same id")
				g.Expect(result.Meta.GetVersion()).To(Equal(tc.expected.result.Meta.GetVersion()), "should be the same version")
			}
		})
	}
}

/*func Test_TrackerStore_Delete(t *testing.T) {

	type testInput struct {
		id int
	}

	type testExpectation struct {
		err    error
		result models.TimeTracker
	}

	expectedMeta := domain.NewMeta()
	expectedMeta.SetDisabled(true)

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

			user, err := repo.Delete(ctx, tc.input.id)

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				g.Expect(user.Meta.GetDisabled()).To(Equal(tc.expected.result.Meta.GetDisabled()), "should be disabled")
			}
		})
	}
}

func Test_TrackerStore_Get(t *testing.T) {

	type testInput struct {
		id int
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
				result: models.TimeTracker{},
				err:    nil,
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

			result, err := repo.Get(ctx, tc.input.id)

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				g.Expect(result.Country).To(Equal(tc.expected.result.Country), "should be the same country")
				g.Expect(result.Nickname).To(Equal(tc.expected.result.Nickname), "should be the same nickname")
				g.Expect(result.FirstName).To(Equal(tc.expected.result.FirstName), "should be the same first name")
				g.Expect(result.LastName).To(Equal(tc.expected.result.LastName), "should be the same last name")
				g.Expect(result.Email).To(Equal(tc.expected.result.Email), "should be the same email")
				g.Expect(result.Password).To(Equal(tc.expected.result.Password), "should be the same password")
				g.Expect(result.ID).To(Equal(tc.expected.result.ID), "should be the same id")
			}
		})
	}
}

func Test_TrackerStore_List(t *testing.T) {

	type testInput struct {
		queryTerms map[string]string
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
			description: "when searching for users with country uk",
			input: testInput{
				queryTerms: map[string]string{
					"country": "uk",
				},
			},
			expected: testExpectation{
				result: []models.TimeTracker{
					{},
				},
				err: nil,
			},
		},
		{
			description: "when listing users",
			input: testInput{
				queryTerms: nil,
			},
			expected: testExpectation{
				result: []models.TimeTracker{
					{
						FirstName: "Test",
						LastName:  "Test",
						Nickname:  "testuser",
						Password:  "qwerty",
						Email:     "example@example.qqq",
						Country:   "uk",
						ID:        1,
					},
					{
						FirstName: "Test",
						LastName:  "Test",
						Nickname:  "testuser-2",
						Password:  "qwerty",
						Email:     "example-db@example.qqq",
						Country:   "ab",
						ID:        2,
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

			result, err := repo.List(ctx, tc.input.queryTerms)

			if tc.expected.err != nil {
				g.Expect(err).To(Equal(tc.expected.err), "should return the expected error")
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "should not return an error")
				for index := range result {
					g.Expect(result[index].Country).To(Equal(tc.expected.result[index].Country), "should be the same country")
					g.Expect(result[index].Nickname).To(Equal(tc.expected.result[index].Nickname), "should be the same nickname")
					g.Expect(result[index].FirstName).To(Equal(tc.expected.result[index].FirstName), "should be the same first name")
					g.Expect(result[index].LastName).To(Equal(tc.expected.result[index].LastName), "should be the same last name")
					g.Expect(result[index].Email).To(Equal(tc.expected.result[index].Email), "should be the same email")
					g.Expect(result[index].Password).To(Equal(tc.expected.result[index].Password), "should be the same password")
					g.Expect(result[index].ID).To(Equal(tc.expected.result[index].ID), "should be the same id")
				}
			}
		})
	}
}
*/
