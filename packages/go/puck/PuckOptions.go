package puck

import (
	"strings"
)

// "strings"

func From(url string) IPuck {
	return NewPuck(from(url))
}

func from(url string) OptionerFn {
	return func(args ...any) error {
		for _, a := range args {
			puck, ok := a.(*Puck)
			if ok {
				if strings.HasPrefix(url, "/") {
					url = "file://" + url
				}
				puck.url = url
				puck.protocol = puck.Protocol()
				puck.format = puck.Format()
			}
		}
		return nil
	}
}
