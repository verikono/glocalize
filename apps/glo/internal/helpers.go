package internal

import (
	"strings"
)

func IsGlobPattern(subject string) bool {
	return strings.Contains(subject, "*")
}
