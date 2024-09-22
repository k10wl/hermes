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

	coreInstance, _ := createCoreAndDB()
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

	coreInstance, _ := createCoreAndDB()
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

	table := []testCase{
		{
			name: "should return all chats is limit overshoots",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 20, 0)
				return c.GetDB().Close
			},
			expected: db_helpers.GenerateChatsSliceN(10),
		},
		{
			name: "should return only 5 first results",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 5, 0)
				return c.GetDB().Close
			},
			expected: db_helpers.GenerateChatsSliceN(5),
		},
		{
			name: "should return 5 last results if start after was specified",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 5, 5)
				return c.GetDB().Close
			},
			expected: db_helpers.GenerateChatsSliceN(10)[5:],
		},
		{
			name: "should return empty slice if pagination overshoots data",
			init: func() closeDB {
				c := prepare(10)
				query = core.NewGetChatsQuery(c, 10, 10)
				return c.GetDB().Close
			},
			expected: db_helpers.GenerateChatsSliceN(0),
		},
	}

	for _, test := range table {
		cleanup := test.init()
		defer cleanup()
		ctx := context.Background()
		query.Execute(ctx)
		for _, c := range query.Result {
			c.TimestampsToNilForTest__()
		}
		expected := test_helpers.UnpointerSlice(test.expected)
		actual := test_helpers.UnpointerSlice(query.Result)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf(
				"Bed result %s\nexpected: %+v\nactual:   %+v\n",
				test.name,
				expected,
				actual,
			)
			continue
		}
	}
}
