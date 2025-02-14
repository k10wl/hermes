package template

import (
	"context"
	"strings"
	"testing"

	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/spf13/cobra"
)

func TestViewEmptyTemplatesList(t *testing.T) {
	type testSubject struct {
		cmd      *cobra.Command
		expected string
		out      *strings.Builder
	}
	type testCase struct {
		name    string
		prepare func() testSubject
	}

	table := []testCase{
		{
			name: "should notify that templates are empty",
			prepare: func() testSubject {
				coreInstance, _ := test_helpers.CreateCore()
				out := &strings.Builder{}
				coreInstance.GetConfig().Stdoout = out
				return testSubject{
					expected: "No templates found",
					cmd:      createViewCommand(coreInstance),
					out:      out,
				}
			},
		},

		{
			name: "should return one template if no query was specified and only one record in db exists",
			prepare: func() testSubject {
				coreInstance, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				if _, err := seeder.SeedTemplatesN(1); err != nil {
					t.Fatalf("failed to seed chats for templates test - %s\n", err.Error())
				}
				out := &strings.Builder{}
				coreInstance.GetConfig().Stdoout = out
				return testSubject{
					expected: `--{{define "1"}}1--{{end}}`,
					cmd:      createViewCommand(coreInstance),
					out:      out,
				}
			},
		},

		{
			name: "should return one template if only one matches search",
			prepare: func() testSubject {
				coreInstance, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				if _, err := seeder.SeedTemplatesN(9); err != nil {
					t.Fatalf("failed to seed chats for templates test - %s\n", err.Error())
				}
				out := &strings.Builder{}
				coreInstance.GetConfig().Stdoout = out
				cmd := createViewCommand(coreInstance)
				if err := cmd.Flags().Set("name", "1"); err != nil {
					t.Fatalf("failed to set flag for test - %s\n", err.Error())
				}
				return testSubject{
					expected: `--{{define "1"}}1--{{end}}`,
					cmd:      cmd,
					out:      out,
				}
			},
		},

		{
			name: "should return two matching templates in descending order",
			prepare: func() testSubject {
				coreInstance, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				if _, err := seeder.SeedTemplatesN(10); err != nil {
					t.Fatalf("failed to seed chats for templates test - %s\n", err.Error())
				}
				out := &strings.Builder{}
				coreInstance.GetConfig().Stdoout = out
				cmd := createViewCommand(coreInstance)
				if err := cmd.Flags().Set("name", "1"); err != nil {
					t.Fatalf("failed to set flag for test - %s\n", err.Error())
				}
				return testSubject{
					expected: `List of templates:

[Name]    10
[Content] --{{define "10"}}10--{{end}}
--------------------
[Name]    1
[Content] --{{define "1"}}1--{{end}}
--------------------
`,
					cmd: cmd,
					out: out,
				}
			},
		},

		{
			name: "should return all templates in descending order",
			prepare: func() testSubject {
				coreInstance, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				if _, err := seeder.SeedTemplatesN(3); err != nil {
					t.Fatalf("failed to seed chats for templates test - %s\n", err.Error())
				}
				out := &strings.Builder{}
				coreInstance.GetConfig().Stdoout = out
				cmd := createViewCommand(coreInstance)
				return testSubject{
					expected: `List of templates:

[Name]    3
[Content] --{{define "3"}}3--{{end}}
--------------------
[Name]    2
[Content] --{{define "2"}}2--{{end}}
--------------------
[Name]    1
[Content] --{{define "1"}}1--{{end}}
--------------------
`,
					cmd: cmd,
					out: out,
				}
			},
		},

		{
			name: "should notify that no results found on unknown template",
			prepare: func() testSubject {
				coreInstance, db := test_helpers.CreateCore()
				seeder := db_helpers.NewSeeder(db, context.Background())
				if _, err := seeder.SeedTemplatesN(3); err != nil {
					t.Fatalf("failed to seed chats for templates test - %s\n", err.Error())
				}
				out := &strings.Builder{}
				coreInstance.GetConfig().Stdoout = out
				cmd := createViewCommand(coreInstance)
				if err := cmd.Flags().Set("name", "does not exist"); err != nil {
					t.Fatalf("failed to set flag for test - %s\n", err.Error())
				}
				return testSubject{
					expected: `No templates found`,
					cmd:      cmd,
					out:      out,
				}
			},
		},
	}

	for _, test := range table {
		subject := test.prepare()
		if err := subject.cmd.Execute(); err != nil {
			t.Fatalf(
				"failed to execute cmd in %q. error - %s",
				test.name,
				err.Error(),
			)
		}
		str := subject.out.String()
		if str != subject.expected {
			t.Fatalf(
				"failed to get expected result in %q\nexpected: %s\nactual:   %s",
				test.name,
				subject.expected,
				str,
			)
		}
	}
}
