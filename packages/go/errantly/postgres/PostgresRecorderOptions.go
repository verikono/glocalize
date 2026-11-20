package PostgresRecorder

// WithConnectionString sets the DSN. It defaults to PGSQL_URL.
func WithConnectionString(connectionString string) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*PostgresRecorder); ok {
				inst.ConnectionString = connectionString
			}
		}
		return nil
	}
}

// WithTable sets the table name. Only ErrorRecord is accepted, because the insert names its columns. @see .agents.md "### PostgresRecorder"
func WithTable(table string) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*PostgresRecorder); ok {
				inst.Table = table
			}
		}
		return nil
	}
}

// WithSystem sets the `system` column, which is how a reader tells one service's errors from another's.
func WithSystem(system string) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*PostgresRecorder); ok {
				inst.System = system
			}
		}
		return nil
	}
}
