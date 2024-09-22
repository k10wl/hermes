package test_helpers

import (
	"os"
	"testing"
)

func Skip(t *testing.T) {
	if os.Getenv("HERMES_TEST_HELPERS_SKIP") == "true" {
		t.SkipNow()
	}
}
