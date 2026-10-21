package tracer

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// maxLogs bounds the in-memory log buffer so a long-running trace cannot grow
// without limit. The oldest lines are dropped first.
const maxLogs = 5000

// nowFunc is the clock the parser stamps undated lines with; a test may replace it.
var nowFunc = time.Now

// Tracer is a full-screen log viewer. It consumes lines from a TracerLogSource,
// parses them, and renders them through a format string with an interactive
// filter and focus.
type Tracer struct {
	app      string
	format   logFormat
	source   TracerLogSource
	ctx      context.Context
	lines    <-chan string
	logs     []TracerLogInterface
	last     *TracerLog
	filter   *TracerLogFilter
	cursor   int
	focused  int
	filterUI bool
	needle   string
	width    int
	height   int
	err      error
}

// NewTracer builds a tracer from options. It returns a fully-wired tracer or an
// error — never a half-built one.
func NewTracer(options ...OptionerFn) (*Tracer, error) {
	tracer := &Tracer{
		app:     "tracer",
		format:  mustParseFormat(DefaultFormat),
		filter:  &TracerLogFilter{},
		focused: -1,
		width:   80,
		height:  24,
	}

	for _, option := range options {
		if err := option(tracer); err != nil {
			return nil, err
		}
	}
	if tracer.source == nil {
		return nil, fmt.Errorf("tracer: no log source configured")
	}
	return tracer, nil
}

// App returns the application name shown in the title bar.
func (t *Tracer) App() string {
	return t.app
}

// Format returns the format string in use.
func (t *Tracer) Format() string {
	return t.format.String()
}

// Logs returns the buffered logs, unfiltered.
func (t *Tracer) Logs() []TracerLogInterface {
	return t.logs
}

// Visible returns the logs that pass the current filter.
func (t *Tracer) Visible() []TracerLogInterface {
	if !t.filter.Active() {
		return t.logs
	}
	visible := make([]TracerLogInterface, 0, len(t.logs))
	for _, log := range t.logs {
		if t.filter.Match(log) {
			visible = append(visible, log)
		}
	}
	return visible
}

// SetLogFilter applies options to the filter and returns the visible logs.
func (t *Tracer) SetLogFilter(options ...OptionerFn) ([]TracerLogInterface, error) {
	for _, option := range options {
		if err := option(t.filter); err != nil {
			return nil, err
		}
	}
	t.clampCursor()
	return t.Visible(), nil
}

// Run starts the full-screen viewer and blocks until the user quits. It owns
// the context the source is followed under, so quitting stops the reader.
func (t *Tracer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t.ctx = ctx
	program := tea.NewProgram(t, tea.WithAltScreen(), tea.WithContext(ctx))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("tracer: run: %w", err)
	}
	return t.err
}

// Init opens the source and starts waiting for its first line.
func (t *Tracer) Init() tea.Cmd {
	ctx := t.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	lines, err := t.source.Follow(ctx)
	if err != nil {
		t.err = err
		return tea.Quit
	}
	t.lines = lines
	return waitForLine(lines)
}

// Update handles a message from the terminal or the source.
func (t *Tracer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width, t.height = msg.Width, msg.Height
		t.clampCursor()
		return t, nil
	case logLineMsg:
		t.appendLine(string(msg))
		return t, waitForLine(t.lines)
	case sourceDoneMsg:
		if msg.err != nil {
			t.err = msg.err
		}
		return t, nil
	case tea.KeyMsg:
		return t, t.handleKey(msg)
	}
	return t, nil
}

// View renders the whole screen.
func (t *Tracer) View() string {
	sections := []string{
		t.renderTitle(),
		t.renderLogs(),
	}
	if t.focused >= 0 {
		sections = append(sections, t.renderFocused())
	}
	if t.filterUI {
		sections = append(sections, t.renderFilter())
	}
	sections = append(sections, t.renderCommands())
	return strings.Join(sections, "\n")
}

// handleKey dispatches a key press to the filter editor or the log navigator.
func (t *Tracer) handleKey(msg tea.KeyMsg) tea.Cmd {
	if t.filterUI {
		return t.handleFilterKey(msg)
	}

	switch msg.String() {
	case "q", "Q", "ctrl+c":
		return tea.Quit
	case "f", "F":
		t.filterUI = true
		t.needle = t.filter.Needle()
		return nil
	case "c", "C":
		t.clearFilter()
		return nil
	case "up", "k":
		t.moveCursor(-1)
		return nil
	case "down", "j":
		t.moveCursor(1)
		return nil
	case "enter":
		t.focusCursor()
		return nil
	case "backspace", "x", "X":
		t.focused = -1
		return nil
	}
	return nil
}

// handleFilterKey edits the needle while the filter UI is open.
func (t *Tracer) handleFilterKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlF, tea.KeyEsc:
		t.filterUI = false
		return nil
	case tea.KeyEnter:
		t.filterUI = false
		t.applyNeedle()
		return nil
	case tea.KeyBackspace:
		if t.needle != "" {
			t.needle = t.needle[:len(t.needle)-1]
		}
		t.applyNeedle()
		return nil
	case tea.KeyRunes:
		t.needle += string(msg.Runes)
		t.applyNeedle()
		return nil
	}
	return nil
}

// applyNeedle pushes the edited needle into the filter.
func (t *Tracer) applyNeedle() {
	if err := t.filter.SetFilterString(t.needle); err != nil {
		t.err = err
		t.needle = t.filter.Needle()
	}
	t.clampCursor()
}

// clearFilter resets the needle and the level restriction.
func (t *Tracer) clearFilter() {
	t.needle = ""
	_ = t.filter.SetFilterString("") // an empty needle cannot fail to compile
	_ = t.filter.SetLevel("")        // an empty level is always valid
	t.clampCursor()
}

// moveCursor shifts the selection by delta within the visible logs.
func (t *Tracer) moveCursor(delta int) {
	visible := len(t.Visible())
	if visible == 0 {
		t.cursor = 0
		return
	}
	t.cursor += delta
	if t.cursor < 0 {
		t.cursor = 0
	}
	if t.cursor >= visible {
		t.cursor = visible - 1
	}
}

// focusCursor pins the selected log into the focused pane.
func (t *Tracer) focusCursor() {
	if len(t.Visible()) == 0 {
		return
	}
	t.focused = t.cursor
}

// clampCursor keeps the cursor and the focus inside the visible range after a
// filter change.
func (t *Tracer) clampCursor() {
	visible := len(t.Visible())
	if visible == 0 {
		t.cursor = 0
		t.focused = -1
		return
	}
	if t.cursor >= visible {
		t.cursor = visible - 1
	}
	if t.cursor < 0 {
		t.cursor = 0
	}
	if t.focused >= visible {
		t.focused = -1
	}
}

// appendLine parses one raw line into the buffer, dropping the oldest when full.
// An indented line is a detail of the log above it rather than a new entry.
func (t *Tracer) appendLine(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}
	if isDetailLine(line) && t.last != nil {
		t.last.AddDetail(strings.TrimSpace(line))
		return
	}
	log := parseLogLine(line, nowFunc())
	t.last = log
	t.logs = append(t.logs, log)
	if len(t.logs) > maxLogs {
		t.logs = t.logs[len(t.logs)-maxLogs:]
	}
	t.clampCursor()
}

// waitForLine returns a command that reads the next line from lines.
func waitForLine(lines <-chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-lines
		if !ok {
			return sourceDoneMsg{}
		}
		return logLineMsg(line)
	}
}
