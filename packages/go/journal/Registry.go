package journal

import "fmt"

var JournalHandlers = make(map[JournalHandlerKind]JournalHandlerProvider)

func RegisterHandler(name JournalHandlerKind, handler JournalHandlerProvider) error {
	handlerName := JournalHandlerKind(name)
	if _, exists := JournalHandlers[handlerName]; exists {
		return fmt.Errorf("Handler with name %s already exists", name)
	}
	JournalHandlers[handlerName] = handler
	return nil
}

func HandlerExists(name string) bool {
	handlerName := JournalHandlerKind(name)
	_, exists := JournalHandlers[handlerName]
	return exists
}

// newHandler instantiates the registered handler of the given kind.
func newHandler(kind JournalHandlerKind) (JournalHandler, error) {
	provider, ok := JournalHandlers[kind]
	if !ok {
		return nil, fmt.Errorf("Handler '%s' is not registered", kind)
	}
	return provider()
}

func init() {
	RegisterHandler(ConsoleJournalHandlerKind, NewConsoleJournalHandler)
}
