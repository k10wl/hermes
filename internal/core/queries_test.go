package core_test

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestGetTemplateByNameQuery(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		expectedResult []*models.Template
		shouldError    bool
	}

	coreInstance, _ := test_helpers.CreateCore()
	var query *core.GetTemplatesByNamesQuery
	if err := core.NewUpsertTemplateCommand(coreInstance,
		`--{{define "hello"}}hello world--{{end}}`).Execute(context.Background()); err != nil {
		panic(err)
	}

	table := []testCase{
		{
			name: "should get template by name",
			init: func() {
				query = core.NewGetTemplatesByNamesQuery(coreInstance, []string{"hello"})
			},
			shouldError: false,
			expectedResult: []*models.Template{
				{
					ID:      1,
					Name:    "hello",
					Content: `--{{define "hello"}}hello world--{{end}}`,
				},
			},
		},
		{
			name: "should return empty array if template does not exist",
			init: func() {
				query = core.NewGetTemplatesByNamesQuery(coreInstance, []string{"does not exist"})
			},
			shouldError:    false,
			expectedResult: []*models.Template{},
		},
	}

	for _, test := range table {
		test.init()
		err := query.Execute(context.Background())
		for _, res := range test.expectedResult {
			res.TimestampsToNilForTest__()
		}
		for _, res := range query.Result {
			res.TimestampsToNilForTest__()
		}
		if test.shouldError {
			if err == nil {
				t.Errorf("%s expected to error, but did not\n", test.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s unexpected error: %v\n", test.name, err)
			continue
		}
		if !reflect.DeepEqual(test.expectedResult, query.Result) {
			t.Errorf(
				"%s - bad result\nexpected: %+v\nactual:   %+v",
				test.name,
				test.expectedResult,
				query.Result,
			)
		}
	}
}

func TestGetTemplatesByRegexp(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		expectedResult []*models.Template
		shouldError    bool
	}

	coreInstance, _ := test_helpers.CreateCore()
	var query *core.GetTemplatesByRegexp
	templates := map[string]string{
		"hello":   `--{{define "hello"}}hello--{{end}}`,
		"hi":      `--{{define "hi"}}hi--{{end}}`,
		"welcome": `--{{define "welcome"}}welcome--{{end}}`,
	}
	for _, template := range templates {
		if err := core.NewUpsertTemplateCommand(
			coreInstance,
			template,
		).Execute(context.Background()); err != nil {
			panic(err)
		}
	}

	table := []testCase{
		{
			name: "should return all templates",
			init: func() {
				query = core.NewGetTemplatesByRegexp(coreInstance, "%")
			},
			expectedResult: []*models.Template{
				{ID: 1, Name: "hello", Content: templates["hello"]},
				{ID: 2, Name: "hi", Content: templates["hi"]},
				{ID: 3, Name: "welcome", Content: templates["welcome"]},
			},
		},
		{
			name: "should return one template",
			init: func() {
				query = core.NewGetTemplatesByRegexp(coreInstance, "hi")
			},
			expectedResult: []*models.Template{
				{ID: 2, Name: "hi", Content: templates["hi"]},
			},
		},
		{
			name: "should return two matching templates",
			init: func() {
				query = core.NewGetTemplatesByRegexp(coreInstance, "h%")
			},
			expectedResult: []*models.Template{
				{ID: 1, Name: "hello", Content: templates["hello"]},
				{ID: 2, Name: "hi", Content: templates["hi"]},
			},
		},
		{
			name: "should return empty array if no matches",
			init: func() {
				query = core.NewGetTemplatesByRegexp(coreInstance, "nothing")
			},
		},
	}

	for _, test := range table {
		test.init()
		err := query.Execute(context.Background())
		for _, res := range test.expectedResult {
			res.TimestampsToNilForTest__()
		}
		for _, res := range query.Result {
			res.TimestampsToNilForTest__()
		}
		if test.shouldError {
			if err == nil {
				t.Errorf("%s expected to error, but did not\n", test.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s unexpected error: %v\n", test.name, err)
			continue
		}
		for _, e := range test.expectedResult {
			expected := test_helpers.UnpointerSlice(test.expectedResult)
			actual := test_helpers.UnpointerSlice(query.Result)
			if !slices.ContainsFunc(query.Result, func(el *models.Template) bool {
				return el.Name == e.Name
			}) {
				t.Errorf(
					"%s - bad result\nexpected: %+v\nactual:   %+v",
					test.name,
					expected,
					actual,
				)
			}
		}
	}
}

func TestGetChatsQuery(t *testing.T) {
	type closeDB func() error
	type testCase struct {
		name     string
		init     func() closeDB
		expected []*models.Chat
	}

	var query *core.GetChatsQuery

	prepare := func(n int64) *core.Core {
		db, err := sqlite3.NewSQLite3(":memory:")
		if err != nil {
			t.Fatalf("failed to create db - %s\n", err)
		}
		c := core.NewCore(db, &settings.Config{})
		ctx := context.Background()
		seeder := db_helpers.NewSeeder(db.DB, ctx)
		err = seeder.SeedChatsN(n)
		if err != nil {
			t.Fatalf("failed to seed db - %s\n", err)
		}
		return c
	}

	reversed := func(data []*models.Chat) []*models.Chat {
		slices.Reverse(data)
		return data
	}

	table := []testCase{
		{
			name: "should return all chats if limit unset",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, -1, -1)
				return c.GetDB().Close
			},
			expected: reversed(db_helpers.GenerateChatsSliceN(10)),
		},
		{
			name: "should return all chats if limit overshoots",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 100, -1)
				return c.GetDB().Close
			},
			expected: reversed(db_helpers.GenerateChatsSliceN(10)),
		},
		{
			name: "should return [10-8] results",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 3, -1)
				return c.GetDB().Close
			},
			expected: reversed(db_helpers.GenerateChatsSliceN(10))[0:3],
		},
		{
			name: "should return [7-5] results",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 3, 8)
				return c.GetDB().Close
			},
			expected: reversed(db_helpers.GenerateChatsSliceN(10))[3:6],
		},
		{
			name: "should return [4-2] results",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 3, 5)
				return c.GetDB().Close
			},
			expected: reversed(db_helpers.GenerateChatsSliceN(10))[6:9],
		},
		{
			name: "should return [1] result",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 3, 2)
				return c.GetDB().Close
			},
			expected: reversed(db_helpers.GenerateChatsSliceN(10))[9:],
		},
		{
			name: "should return empty slice if pagination overshoots data",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 10, 1)
				return c.GetDB().Close
			},
			expected: db_helpers.GenerateChatsSliceN(0),
		},
	}

	for _, test := range table {
		cleanup := test.init()
		defer cleanup()
		ctx := context.Background()
		err := query.Execute(ctx)
		if err != nil {
			t.Errorf("Failed to execute query - %s\n\n", err)
			continue
		}
		for _, c := range query.Result {
			c.TimestampsToNilForTest__()
		}
		expected := test_helpers.UnpointerSlice(test.expected)
		actual := test_helpers.UnpointerSlice(query.Result)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf(
				"Bed result %s\nexpected: %+v\nactual:   %+v\n\n",
				test.name,
				expected,
				actual,
			)
			continue
		}
	}
}

func TestTemplatesQuery(t *testing.T) {
	type subject struct {
		cmd      *core.GetTemplatesQuery
		expected []models.Template
	}
	type testCase struct {
		name    string
		prepare func() subject
	}

	table := []testCase{
		{
			name: "should return all templates",
			prepare: func() subject {
				c, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				templates, err := seeder.SeedTemplatesN(2)
				if err != nil {
					t.Fatalf("failed to seed templates - %q", err)
				}
				expected := test_helpers.UnpointerSlice(templates)
				slices.Reverse(expected)
				return subject{
					cmd:      core.NewGetTemplatesQuery(c, -1, -1, ""),
					expected: expected,
				}
			},
		},

		{
			name: "should return limited amount of answers",
			prepare: func() subject {
				c, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				templates, err := seeder.SeedTemplatesN(100)
				if err != nil {
					t.Fatalf("failed to seed templates - %q", err)
				}
				expected := test_helpers.UnpointerSlice(templates)
				slices.Reverse(expected)
				return subject{
					cmd:      core.NewGetTemplatesQuery(c, -1, 10, ""),
					expected: expected[:10],
				}
			},
		},

		{
			name: "should return results after specified id",
			prepare: func() subject {
				c, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				templates, err := seeder.SeedTemplatesN(100)
				if err != nil {
					t.Fatalf("failed to seed templates - %q", err)
				}
				expected := test_helpers.UnpointerSlice(templates)
				slices.Reverse(expected)
				return subject{
					cmd:      core.NewGetTemplatesQuery(c, 91, 10, ""),
					expected: expected[10:20],
				}
			},
		},

		{
			name: "should return partial results if out of bounds",
			prepare: func() subject {
				c, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				templates, err := seeder.SeedTemplatesN(100)
				if err != nil {
					t.Fatalf("failed to seed templates - %q", err)
				}
				expected := test_helpers.UnpointerSlice(templates)
				slices.Reverse(expected)
				return subject{
					cmd:      core.NewGetTemplatesQuery(c, 6, 10, ""),
					expected: expected[95:],
				}
			},
		},

		{
			name: "should return all results if startBeforeID is -1 and limit is -1",
			prepare: func() subject {
				c, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				templates, err := seeder.SeedTemplatesN(100)
				if err != nil {
					t.Fatalf("failed to seed templates - %q", err)
				}
				expected := test_helpers.UnpointerSlice(templates)
				slices.Reverse(expected)
				return subject{
					cmd:      core.NewGetTemplatesQuery(c, -1, -1, ""),
					expected: expected,
				}
			},
		},

		{
			name: "should search name matches",
			prepare: func() subject {
				c, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				_, err := seeder.SeedTemplatesN(100)
				if err != nil {
					t.Fatalf("failed to seed templates - %q", err)
				}
				return subject{
					cmd: core.NewGetTemplatesQuery(c, -1, -1, "22"),
					expected: []models.Template{
						{
							Name:    "22",
							ID:      22,
							Content: `--{{template "22"}}22--{{end}}`,
						},
					},
				}
			},
		},
	}

	for _, test := range table {
		sub := test.prepare()
		err := sub.cmd.Execute(context.Background())
		if err != nil {
			t.Fatalf("unexpected error in %q - %q", test.name, err)
		}
		actual := test_helpers.UnpointerSlice(test_helpers.ResetSliceTime(sub.cmd.Result))
		if !reflect.DeepEqual(sub.expected, actual) {
			t.Fatalf(
				"did not get expected result in %q\nexpected: %+v\nactual:   %+v\n",
				test.name,
				sub.expected,
				actual,
			)
		}
	}
}
