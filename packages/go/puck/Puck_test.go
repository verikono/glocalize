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

	t.Run("resolves reads a valid glob source (wildcard) to a string", func(t *testing.T) {
		validPath := "/home/bren/mcpikon/internal/data/schema/provider.prisma"
		p := puck.From(validPath)
		pk, ok := p.(*puck.Puck)
		if !ok {
			t.Fatal("not a puck")
		}
		if pk.Path() != "/home/bren/mcpikon/internal/data/schema/provider.prisma" {
			t.Fatal("path should not be set (source is a wildcard/glob)")
		}
		if pk.Url() != "file:///home/bren/mcpikon/internal/data/schema/provider.prisma" {
			t.Fatal("path should not be set (source is a wildcard/glob)")
		}
		if pk.Format() != "prisma" {
			t.Fatal("format should be 'prisma'")
		}
		if pk.Protocol() != "file" {
			t.Fatalf("protocol should be 'file' for this glob: %s", validPath)
		}
		if pk.ResourceType() != "file" {
			t.Fatalf("protocol should be 'file' for this glob: %s", validPath)
		}
		contents, err := pk.Read()
		if err != nil {
			t.Fatalf("read of a valid glob failed: %s", validPath)
		}
		if contents == "" {
			t.Fatalf("read of valid glob returned no content: %s", validPath)
		}
	})

	t.Run("resolves reads a valid glob source (wildcard) to a string", func(t *testing.T) {
		validglob := "/home/bren/mcpikon/internal/data/schema/*.prisma"
		p := puck.From(validglob)
		pk, ok := p.(*puck.Puck)
		if !ok {
			t.Fatal("not a puck")
		}
		if pk.Path() != "" {
			t.Fatal("path should not be set (source is a wildcard/glob)")
		}
		if pk.Format() != "prisma" {
			t.Fatal("format should be 'prisma'")
		}
		if pk.Protocol() != "file" {
			t.Fatalf("protocol should be 'file' for this glob: %s", validglob)
		}
		if pk.ResourceType() != "file" {
			t.Fatalf("protocol should be 'file' for this glob: %s", validglob)
		}
		contents, err := pk.Read()
		if err != nil {
			t.Fatalf("read of a valid glob failed: %s", validglob)
		}
		if contents == "" {
			t.Fatalf("read of valid glob returned no content: %s", validglob)
		}
	})

}
