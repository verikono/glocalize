package puck

import "fmt"

func From(url string) OptionerFn {
	return func(args ...any) error {
		puck, ok := args[0].(*Puck)
		if !ok {
			return fmt.Errorf("expected *Puck, got %T", args[0])
		}
		puck.url = url
		return nil
	}
}
