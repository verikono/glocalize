package errantly

import (
	KT "github.com/verikono/glocalize/types/go"
)

type (
	OptionerFn = KT.OptionerFn
	URLString  = KT.URLString
)

// -- interfaces --
type (
	ErrantlyError   = KT.ErrantlyError
	Journalist      = KT.Journalist
	ErrantlyHandler = KT.ErrantlyHandler
)

// -- structs --
type (
	NotatedCall         = KT.NotatedCall
	ErrantlyMessageMeta = KT.ErrantlyMessageMeta
)

type (
	ErrantlyPackage      = KT.ErrantlyPackage
	JournalPackage       = KT.JournalPackage
	CoreDependencies     = KT.CoreDependencies
	ErrantlyDependencies = KT.ErrantlyDependencies
)
