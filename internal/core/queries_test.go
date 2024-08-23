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
		expectedResult models.Template
		shouldError    bool
	}

	coreInstance, db := createCoreAndDB()
	var query *core.GetTemplateByNameQuery

	table := []testCase{
		{
			name: "should get template by name",
			init: func() {
				_, err := db.CreateTemplate(
					context.Background(),
					"hello",
					`{{define "hello"}}hello world{{end}}`,
				)
				if err != nil {
					panic(err)
				}
				query = core.NewGetTemplateByNameQuery(coreInstance, "hello")
			},

			shouldError: false,
			expectedResult: models.Template{
				ID:      1,
				Name:    "hello",
				Content: `{{define "hello"}}hello world{{end}}`,
			},
		},
		{
			name: "should error if template does not exist",
			init: func() {
				query = core.NewGetTemplateByNameQuery(coreInstance, "does not exist")
			},

			shouldError:    true,
			expectedResult: models.Template{},
		},
	}

	for _, test := range table {
		test.init()
		err := query.Execute(context.Background())
		test.expectedResult.TimestampsToNilForTest__()
		query.Result.TimestampsToNilForTest__()
		if test.shouldError && err == nil {
			t.Errorf("%s expected to error, but did not\n", test.name)
			continue
		}
		if !test.shouldError && err != nil {
			t.Errorf("%s unexpected error: %v\n", test.name, err)
			continue
		}
		if !reflect.DeepEqual(test.expectedResult, *query.Result) {
			t.Errorf(
				"%s - bad result\nexpected: %+v\nactual:   %+v",
				test.name,
				test.expectedResult,
				*query.Result,
			)
		}
	}
}
