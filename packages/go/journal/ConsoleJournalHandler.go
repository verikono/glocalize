package journal

import (
	"fmt"
	"io"
	"os"
)

const ConsoleJournalHandlerKind JournalHandlerKind = "console"

// ConsoleJournalHandler writes each journal entry as a formatted line to an
// io.Writer (os.Stdout by default). Inject a writer via WithWriter to capture
// output, e.g. in tests.
type ConsoleJournalHandler struct {
	kind   JournalHandlerKind
	writer io.Writer
}

func NewConsoleJournalHandler(options ...OptionerFn) (JournalHandler, error) {
	inst := ConsoleJournalHandler{
		kind:   JournalHandlerKind(ConsoleJournalHandlerKind),
		writer: os.Stdout,
	}
	for _, option := range options {
		if err := option(&inst); err != nil {
			return nil, err
		}
	}
	return &inst, nil
}

func (c *ConsoleJournalHandler) GetKind() JournalHandlerKind {
	return ConsoleJournalHandlerKind
}

// ProcessEntry writes the entry as `TIMESTAMP LEVEL CALLER MESSAGE`.
func (c *ConsoleJournalHandler) ProcessEntry(entry JournalEntryProcessor) error {
	if c.writer == nil {
		c.writer = os.Stdout
	}
	_, err := fmt.Fprintf(c.writer, "%d %s %s %s\n",
		entry.GetTimestamp(), entry.GetLevel(), entry.GetCaller(), entry.GetMessage())
	return err
}
