package agents

import (
	"context"
	"fmt"
)

type Answer = string

// AskAQuestion builds a Prosti from the environment, sends question, waits for the complete answer and returns it.
//
// Usage:
//
//	answer, err := agents.AskAQuestion("how long does it take to fly to pluto?")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(answer)
func AskAQuestion(question string) (Answer, error) {
	if question == "" {
		return "", fmt.Errorf("agents: a question is required")
	}
	agent, err := NewProsti()
	if err != nil {
		return "", fmt.Errorf("agents: build prosti: %w", err)
	}
	events, err := agent.Query(context.Background(), question)
	if err != nil {
		return "", fmt.Errorf("agents: query: %w", err)
	}
	answer, err := Collect(events)
	if err != nil {
		return "", fmt.Errorf("agents: collect answer: %w", err)
	}
	return answer, nil
}
