package core_test

import (
	"context"
	"reflect"
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
	if err := core.NewCreateTemplateCommand(coreInstance,
		`{{define "hello"}}hello world{{end}}`).Execute(context.Background()); err != nil {
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
					Content: `{{define "hello"}}hello world{{end}}`,
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
