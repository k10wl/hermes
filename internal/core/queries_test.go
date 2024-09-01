package core_test

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
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
			expected := unpointerTemplateSlice(test.expectedResult)
			actual := unpointerTemplateSlice(query.Result)
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

func unpointerTemplateSlice(arg []*models.Template) []models.Template {
	res := make([]models.Template, len(arg))
	for i, v := range arg {
		res[i] = *v
	}
	return res
}
