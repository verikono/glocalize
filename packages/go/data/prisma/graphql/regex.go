package graphql

import "regexp"

var (
	modelStartRe = regexp.MustCompile(`^\s*model.*`)
	enumStartRe  = regexp.MustCompile(`^\s*enum.*`)
	modelEndRe   = regexp.MustCompile(`^\s*}.*`)
	commentLine  = regexp.MustCompile(`^\s*\/\/.+$`)
	attributeRe  = regexp.MustCompile(`^\s*@@.*`)
	// a "//@@ property" line declares a field as if it were a normal property.
	propertyDirectiveRe = regexp.MustCompile(`^\s.\/\/@@ property`)
	// skip a field or enum option whose name starts with whitespace or contains "//".
	skipNameRe = regexp.MustCompile(`^\s|\/\/`)
)
