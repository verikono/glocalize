package journal

func Must[T any](value T, err error) T {
	if err != nil {
		panic(E.MustFailure())
	}
	return value
}
