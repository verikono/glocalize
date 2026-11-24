package journal

import (
	KT "github.com/verikono/glocalize/types/go"
)

var E *ErrantlyPackage

func Provider() *JournalPackage {
	return &JournalPackage{
		NewJournal: func(name string, options ...OptionerFn) (Journalist, error) {
			return NewJournal(JournalName(name), options...)
		},
		NewNoopJournal: func() Journalist {
			return KT.NewNoopJournal()
		},
	}
}

func RegisterDeps(deps *JournalDependencies) error {
	E = deps.ErrantlyPackage
	return nil
}
