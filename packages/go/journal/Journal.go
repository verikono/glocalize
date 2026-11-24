package journal

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Journal struct {
	Name     JournalName
	lock     *sync.Mutex
	handlers []JournalHandler
	console  bool
}

func NewJournal(name JournalName, options ...OptionerFn) (*Journal, error) {
	if err := (&name).Validate(); err != nil {
		return nil, err
	}
	if journalExists(name) {
		return nil, fmt.Errorf("Journal with name %s already exists", name)
	}
	inst := Journal{
		Name:     name,
		console:  true,
		handlers: []JournalHandler{},
		lock:     &sync.Mutex{},
	}
	for _, option := range options {
		if err := option(&inst); err != nil {
			return nil, err
		}
	}
	if inst.console {
		consHdlr, err := NewConsoleJournalHandler(options...)
		if err != nil {
			return nil, err
		}
		inst.handlers = append(inst.handlers, consHdlr)
	}
	JournalDict[name] = inst
	return &inst, nil
}

func (e *Journal) Error() string {
	return "Journal: An error occurred in the Journal package."
}

// -- Journalist --

func (j *Journal) GetKind() JournalHandlerKind { return ConsoleJournalHandlerKind }
func (j *Journal) GetName() JournalName        { return j.Name }

// SetHandler replaces the journal's handlers with a single handler of the
// given kind (which must have been registered via RegisterHandler).
func (j *Journal) SetHandler(kind JournalHandlerKind) error {
	h, err := newHandler(kind)
	if err != nil {
		return err
	}
	j.lock.Lock()
	defer j.lock.Unlock()
	j.handlers = []JournalHandler{h}
	return nil
}

// Severity-ordered output methods (highest to lowest). Each formats its
// arguments into the message before dispatching to the journal's handlers.

func (j *Journal) Critical(format string, args ...any) {
	j.emit(SeverityLevelCritical, captureCaller(2), format, args...)
}
func (j *Journal) Error_(format string, args ...any) {
	j.emit(SeverityLevelError, captureCaller(2), format, args...)
}
func (j *Journal) Warn(format string, args ...any) {
	j.emit(SeverityLevelWarn, captureCaller(2), format, args...)
}
func (j *Journal) Notice(format string, args ...any) {
	j.emit(SeverityLevelNotice, captureCaller(2), format, args...)
}
func (j *Journal) Info(format string, args ...any) {
	j.emit(SeverityLevelInfo, captureCaller(2), format, args...)
}
func (j *Journal) Debug(format string, args ...any) {
	j.emit(SeverityLevelDebug, captureCaller(2), format, args...)
}
func (j *Journal) Trace(format string, args ...any) {
	j.emit(SeverityLevelTrace, captureCaller(2), format, args...)
}
func (j *Journal) Verbose(format string, args ...any) {
	j.emit(SeverityLevelVerbose, captureCaller(2), format, args...)
}
func (j *Journal) Prolix(format string, args ...any) {
	j.emit(SeverityLevelProlix, captureCaller(2), format, args...)
}
func (j *Journal) Log(format string, args ...any) { j.Info(format, args...) }

// emit formats the message into a JournalEntry and dispatches it to every
// handler. A failing handler is skipped so one broken sink cannot take down
// the journal.
func (j *Journal) emit(level SeverityLevel, caller, format string, args ...any) {
	entry := &JournalEntry{
		timestamp:  EpochNanoSeconds(time.Now().UnixNano()),
		level:      level,
		message:    fmt.Sprintf(format, args...),
		caller:     caller,
		stacktrace: "not implemented for now.",
	}
	j.lock.Lock()
	defer j.lock.Unlock()
	for _, h := range j.handlers {
		_ = h.ProcessEntry(entry)
	}
}

func (j *Journal) GetChildren() []*Journal {
	children := []*Journal{}
	prefix := string(j.Name) + "."
	for k, v := range JournalDict {
		if strings.HasPrefix(string(k), prefix) {
			// Take a copy to avoid referencing the loop variable.
			child := v
			children = append(children, &child)
		}
	}
	return children
}
