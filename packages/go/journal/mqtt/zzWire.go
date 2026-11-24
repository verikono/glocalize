package mqtt

import (
	"github.com/verikono/glocalize/packages/go/journal"
)

var (
	L *JournalPackage
	E *ErrantlyPackage
)

func Provide() *MqttJournalHandlerPackage {
	inst := MqttJournalHandlerPackage{
		NewMqttJournalHandler: NewMqttJournalHandler,
	}
	return &inst
}

func RegisterDeps(coreDeps *CoreDependencies) error {
	L = coreDeps.JournalPackage
	E = coreDeps.ErrantlyPackage
	return nil
}

func init() {
	journal.RegisterHandler("mqtt", NewMqttJournalHandler)
}
