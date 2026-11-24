package mqtt

import (
	KT "github.com/verikono/glocalize/types/go"
)

type (
	JournalHandlerKind = KT.JournalHandlerKind
	OptionerFn         = KT.OptionerFn
)

// -- interfaces -- //
type (
	JournalHandler        = KT.JournalHandler
	JournalEntryProcessor = KT.JournalEntryProcessor
)

// -- packages -- //
type (
	MqttJournalHandlerPackage = KT.MqttJournalHandlerPackage
	CoreDependencies          = KT.CoreDependencies
	JournalPackage            = KT.JournalPackage
	ErrantlyPackage           = KT.ErrantlyPackage
)
