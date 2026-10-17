package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/verikono/glocalize/apps/glo/internal"

	"github.com/verikono/glocalize/packages/go/puck"
)

func TestGloParse(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not determine working directory: %s", err)
	}

	source := filepath.Join(wd, "testdata", "sample.prisma")
	target := filepath.Join(wd, "testdata", "out.graphql")

	t.Run("parses a known file", func(t *testing.T) {

		result, err := internal.ParseSchema(source, target)
		if err != nil {
			t.Fatalf("received error: %s", err)
		}
		if result == nil {
			t.Fatalf("expected result, got an empty string")
		}

		tpuck := puck.From(target)
		if !tpuck.Exists() {
			t.Fatalf("expected target graphql file which does not exist")
		}

		if tpuck.Size() == 0 {
			t.Fatalf("target graphql file empty")
		}

	})

	t.Cleanup(func() {
		_ = os.Remove(target)
	})

}
