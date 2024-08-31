package launch

import (
	"fmt"
	"testing"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

const (
	NewChat = iota
	UpsertTemplate
	LastChat
	ViewTemplates
	DeleteTemplate
)

var optionMap = map[int]string{
	NewChat:        "NewChat",
	UpsertTemplate: "UpsertTemplate",
	LastChat:       "LastChat",
	ViewTemplates:  "ViewTemplates",
	DeleteTemplate: "DeleteTemplate",
}

type testCLIOptions struct {
	lastRecorded map[string]bool
}

func (o *testCLIOptions) reset() {
	o.lastRecorded = map[string]bool{
		optionMap[NewChat]:        false,
		optionMap[UpsertTemplate]: false,
		optionMap[LastChat]:       false,
		optionMap[ViewTemplates]:  false,
		optionMap[DeleteTemplate]: false,
	}
}

func (o *testCLIOptions) record(option string) {
	fmt.Printf("o.lastRecorded: %v\n", o.lastRecorded)
	if _, ok := o.lastRecorded[option]; !ok {
		panic(fmt.Sprintf("option does not exist, test is invalid - %q", option))
	}
	o.lastRecorded[option] = true
}

func (o *testCLIOptions) getLast() string {
	last := "none"
	for key, val := range o.lastRecorded {
		if val {
			last = key
			break
		}
	}
	return last
}

func TestCLIStrategy(t *testing.T) {
	type testCase struct {
		name        string
		config      []settings.Config
		expected    string
		shouldError bool
	}

	options := testCLIOptions{}
	lauch := newLaunchCLI(&testStrategies{options: &options})

	table := []testCase{
		{
			name:     "should call start new config when no special flags were provided",
			expected: optionMap[NewChat],
		},
		{
			name: "should call upsert template if UpsertTemplate value was provided",
			config: []settings.Config{{
				TemplateFlags: settings.TemplateFlags{UpsertTemplate: "upsert"},
			}},
			expected: optionMap[UpsertTemplate],
		},
		{
			name: "should call view templates if ViewTemplates value was provided",
			config: []settings.Config{{
				TemplateFlags: settings.TemplateFlags{ListTemplates: "view"},
			}},
			expected: optionMap[ViewTemplates],
		},
		{
			name: "should call last chat if Last was true",
			config: []settings.Config{{
				CLIFlags: settings.CLIFlags{Last: true},
			}},
			expected: optionMap[LastChat],
		},
		{
			name: "should error on conflicting flags",
			config: []settings.Config{{
				TemplateFlags: settings.TemplateFlags{
					ListTemplates:  "some",
					UpsertTemplate: "some",
				},
			},

				{
					TemplateFlags: settings.TemplateFlags{
						ListTemplates:  "some",
						DeleteTemplate: "some",
					},
				},

				{
					TemplateFlags: settings.TemplateFlags{
						DeleteTemplate: "some",
						UpsertTemplate: "some",
					},
				},
			},
			expected:    "no result, should error instead",
			shouldError: true,
		},
		{
			name: "should call delete if DeleteTemplate has value",
			config: []settings.Config{{
				TemplateFlags: settings.TemplateFlags{DeleteTemplate: "name"},
			}},
			expected: optionMap[DeleteTemplate],
		},
	}

	for _, test := range table {
		for _, config := range test.config {
			options.reset()
			err := lauch.Execute(&core.Core{}, &config)
			actual := options.getLast()
			if test.shouldError {
				if err == nil {
					t.Errorf("%q expected to error, but didn't. Last called %q", test.name, actual)
				}
				continue
			}
			if err != nil {
				t.Errorf("%q unexpected error: %v", test.name, err)
				continue
			}
			if actual != test.expected {
				t.Errorf(
					"\n%q bad return\nexpected: %s\nactual:   %s\n\n",
					test.name,
					test.expected,
					actual,
				)
			}
		}
	}
}

type testStrategies struct{ options *testCLIOptions }

func (ts *testStrategies) NewChat(*core.Core, *settings.Config) error {
	ts.options.record(optionMap[NewChat])
	return nil
}
func (ts *testStrategies) LastChat(*core.Core, *settings.Config) error {
	ts.options.record(optionMap[LastChat])
	return nil
}
func (ts *testStrategies) ListTemplates(*core.Core, *settings.Config) error {
	ts.options.record(optionMap[ViewTemplates])
	return nil
}
func (ts *testStrategies) UpsertTemplate(*core.Core, *settings.Config) error {
	ts.options.record(optionMap[UpsertTemplate])
	return nil
}
func (ts *testStrategies) DeleteTemplate(*core.Core, *settings.Config) error {
	ts.options.record(optionMap[DeleteTemplate])
	return nil
}
