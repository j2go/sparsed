package sparsed

import (
	"flag"
	"os"
	"testing"
)

// TestMain main for tests
func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
