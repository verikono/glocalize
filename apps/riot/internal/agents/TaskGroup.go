package agents

import "fmt"

type TaskGroup struct {
	agents     []Prosti
	supervisor Prosti
	// resources MCP services
	context string
}

func (tg *TaskGroup) AddMember() error {
	return fmt.Errorf("not implemented")
}

func (tg *TaskGroup) DropMember() error {
	return fmt.Errorf("not implemented")
}
