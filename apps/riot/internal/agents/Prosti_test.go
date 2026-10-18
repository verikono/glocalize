package agents_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudwego/eino-ext/components/model/ark"
	A "github.com/verikono/glocalize/apps/riot/internal/agents"
	U "github.com/verikono/glocalize/packages/go/utils"
)

// TestMain loads every `.env` between the repository root and this package before any test runs, so tests see the same environment the binary does.
func TestMain(m *testing.M) {
	if err := U.WalkDotEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "load dotenv: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestNewProsti(t *testing.T) {
	agent, err := A.NewProsti(A.WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("NewProsti() returned an error: %v", err)
	}
	if agent == nil {
		t.Fatal("NewProsti() returned a nil agent")
	}
	if !agent.Active() {
		t.Fatal("NewProsti() returned an inactive agent")
	}
	if agent.GetModel() == nil {
		t.Error("GetModel() is nil")
	}
	if agent.GetAgent() == nil {
		t.Error("GetAgent() is nil")
	}
	if agent.GetRunner() == nil {
		t.Error("GetRunner() is nil")
	}
	if agent.GetContext() == nil {
		t.Error("GetContext() is nil")
	}
	if agent.GetName() == "" {
		t.Error("GetName() is empty")
	}
}

func TestNewProstiWithOptions(t *testing.T) {
	agent, err := A.NewProsti(
		A.WithName("test-agent"),
		A.WithDescription("a test agent"),
		A.WithInstruction("answer in one word"),
		A.WithModel("test-model"),
		A.WithAPIKey("test-key"),
		A.WithStreaming(false),
	)
	if err != nil {
		t.Fatalf("NewProsti() with options returned an error: %v", err)
	}
	if got := agent.GetName(); got != "test-agent" {
		t.Errorf("GetName() = %q, want %q", got, "test-agent")
	}
	if got := agent.GetAgent().Name(context.Background()); got != "test-agent" {
		t.Errorf("agent Name() = %q, want %q", got, "test-agent")
	}
	if got := agent.GetAgent().Description(context.Background()); got != "a test agent" {
		t.Errorf("agent Description() = %q, want %q", got, "a test agent")
	}
}

func TestNewProstiWithModelConfig(t *testing.T) {
	agent, err := A.NewProsti(A.WithModelConfig(ark.ChatModelConfig{
		Model:  "configured-model",
		APIKey: "configured-key",
	}))
	if err != nil {
		t.Fatalf("NewProsti() with model config returned an error: %v", err)
	}
	if !agent.Active() {
		t.Fatal("agent built from a model config is not active")
	}
}

func TestNewProstiRejectsFailingOption(t *testing.T) {
	failing := func(args ...any) error {
		return context.Canceled
	}
	if _, err := A.NewProsti(failing); err == nil {
		t.Fatal("NewProsti() accepted a failing option")
	}
}

func TestBrothelLifecycle(t *testing.T) {
	brothel, err := A.NewAgenticBrothel()
	if err != nil {
		t.Fatalf("NewAgenticBrothel() returned an error: %v", err)
	}

	agent, err := brothel.NewAgent("alpha", A.WithName("alpha"), A.WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("NewAgent() returned an error: %v", err)
	}
	if agent.GetName() != "alpha" {
		t.Errorf("GetName() = %q, want %q", agent.GetName(), "alpha")
	}

	got, err := brothel.GetAgent("alpha")
	if err != nil {
		t.Fatalf("GetAgent() returned an error: %v", err)
	}
	if got != agent {
		t.Error("GetAgent() returned a different agent")
	}

	if ids := brothel.ListAgents(); len(ids) != 1 || ids[0] != "alpha" {
		t.Errorf("ListAgents() = %v, want [alpha]", ids)
	}

	if err := brothel.RegisterAgent("alpha", agent); err == nil {
		t.Error("RegisterAgent() accepted a duplicate identifier")
	}

	if err := brothel.RemoveAgent("alpha"); err != nil {
		t.Fatalf("RemoveAgent() returned an error: %v", err)
	}
	if _, err := brothel.GetAgent("alpha"); err == nil {
		t.Error("GetAgent() found a removed agent")
	}
	if err := brothel.RemoveAgent("alpha"); err == nil {
		t.Error("RemoveAgent() accepted an unknown identifier")
	}
}

func TestBrothelRejectsEmptyIdentifier(t *testing.T) {
	brothel, err := A.NewAgenticBrothel()
	if err != nil {
		t.Fatalf("NewAgenticBrothel() returned an error: %v", err)
	}
	if _, err := brothel.NewAgent(""); err == nil {
		t.Error("NewAgent() accepted an empty identifier")
	}
	if err := brothel.RegisterAgent("", nil); err == nil {
		t.Error("RegisterAgent() accepted an empty identifier")
	}
}

func TestAskQuestion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live model call in short mode")
	}
	for _, key := range []string{"RIOT_DEFAULT_MODEL", "RIOT_MODEL_KEY", "RIOT_MODEL_BASEURL"} {
		if os.Getenv(key) == "" {
			t.Skipf("%s is not set; skipping live model call", key)
		}
	}

	answer, err := A.AskAQuestion("how long does it take to fly to pluto?")
	if err != nil {
		t.Fatalf("AskAQuestion() returned an error: %v", err)
	}
	if answer == "" {
		t.Fatal("AskAQuestion() returned an empty answer")
	}
	t.Logf("answer: %s", answer)
	fmt.Printf("\n\nAnswer:\n%s\n\n", answer)
}
