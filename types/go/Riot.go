package types

type CoreDependencies struct {
	JournalPackage  *JournalPackage
	ErrantlyPackage *ErrantlyPackage
	Must            func(value any, err error) any
}
