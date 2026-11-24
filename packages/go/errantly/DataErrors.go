package errantly

import "fmt"

// WithBytes attaches the request body bytes to an ErrantlyError's message.
func WithBytes(data []byte) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			issue := findEmbeddedIssue[BytesIssue](target)
			if issue == nil {
				continue
			}
			issue.Bytes = data
			if e, ok := target.(ErrantlyError); ok {
				e.SetMessage(
					fmt.Sprintf("%s — body (%d bytes)", e.GetMessage(), len(data)),
				)
			}
		}
		return nil
	}
}
