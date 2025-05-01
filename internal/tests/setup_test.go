package tests

import (
	"os"
	"testing"

	"github.com/kootic/koogo/internal/tests/testutils"
)

// TestMain is the entry point for the integration tests.
func TestMain(m *testing.M) {
	os.Exit(testutils.RunIntegrationTests(m))
}
