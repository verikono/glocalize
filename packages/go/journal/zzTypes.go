package journal

import (
	KT "github.com/verikono/glocalize/types/go"
)

type (
	JournalName            = KT.JournalName
	JournalHandlerKind     = KT.JournalHandlerKind
	OptionerFn             = KT.OptionerFn
	JournalHandlerProvider = KT.JournalHandlerProvider
	EpochNanoSeconds       = KT.EpochNanoSeconds
	SeverityLevel          = KT.SeverityLevel
)

// -- constants -- //

const (
	SeverityLevelDebug    = KT.SeverityLevelDebug
	SeverityLevelTrace    = KT.SeverityLevelTrace
	SeverityLevelProlix   = KT.SeverityLevelProlix
	SeverityLevelVerbose  = KT.SeverityLevelVerbose
	SeverityLevelInfo     = KT.SeverityLevelInfo
	SeverityLevelNotice   = KT.SeverityLevelNotice
	SeverityLevelWarn     = KT.SeverityLevelWarn
	SeverityLevelError    = KT.SeverityLevelError
	SeverityLevelCritical = KT.SeverityLevelCritical
	SeverityLevelFatal    = KT.SeverityLevelFatal
)

type (
	Journalist            = KT.Journalist
	JournalHandler        = KT.JournalHandler
	JournalEntryReader    = KT.JournalEntryReader
	JournalEntryProcessor = KT.JournalEntryProcessor
)

// -- packages -- //

type (
	CoreDependencies    = KT.CoreDependencies
	JournalDependencies = KT.JournalDependencies
	JournalPackage      = KT.JournalPackage
	ErrantlyPackage     = KT.ErrantlyPackage
)
