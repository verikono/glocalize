package errantly

import "fmt"

// WithURL attaches the target URL to an ErrantlyError's message.
func WithURL(url string) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			issue := findEmbeddedIssue[URLIssue](target)
			if issue == nil {
				continue
			}
			issue.URL = url
			if e, ok := target.(ErrantlyError); ok {
				e.SetMessage(
					fmt.Sprintf("%s — url: %s", e.GetMessage(), url),
				)
			}
		}
		return nil
	}
}
