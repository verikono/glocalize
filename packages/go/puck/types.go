package puck

import (
	GT "github.com/verikono/glocalize/types/go"
)

type (
	OptionerFn = GT.OptionerFn
)

type (
	IPuck interface {
		Protocol() UrlProtocol
		Path() string
		Format() DataFormat
	}
)
