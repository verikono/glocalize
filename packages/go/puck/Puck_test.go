package puck_test

import (
	"fmt"
	"testing"

	"github.com/verikono/glocalize/packages/go/puck"
)

func TestPuck(t *testing.T) {

	t.Run("resolves URL from absolute path", func(t *testing.T) {
		p := puck.From("/home/")
		pk, ok := p.(*puck.Puck)
		if !ok {
			t.Fatal("not a puck")
		}
		path := pk.Path()
		if path == "/dev/null" {
			t.Fatal("path not set")
		}
		fmt.Printf("path: %s", path)
	})
}
