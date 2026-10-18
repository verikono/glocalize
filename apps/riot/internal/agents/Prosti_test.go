package agents_test

import (
	"testing"

	A "github.com/verikono/glocalize/apps/riot/internal/agents"
)

func TestProsti(t *testing.T) {
	agent, err := A.NewProsti()
	if err != nil {
		t.Fatalf("failed")
	}
	_ = agent
}
