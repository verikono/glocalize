package tracer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logLineMsg carries one raw line read from the source.
type logLineMsg string

// defaultWidth and defaultHeight are the fallbacks used before the terminal
// reports its size, or when it reports a zero size (a pty with no window size).
const (
	defaultWidth  = 80
	defaultHeight = 24
)

// sourceDoneMsg reports that the source stopped, with an error when it failed.
type sourceDoneMsg struct {
	err error
}

// SourceDoneMsg is the exported form of sourceDoneMsg, so a test driving the
// model directly can recognise the end of the stream.
type SourceDoneMsg = sourceDoneMsg

// Styles for the viewer. They are package-level so every render path agrees.
var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	ruleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	selectedStyle = lipgloss.NewStyle().Reverse(true)
	focusedStyle  = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	hintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	needleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

// mustParseFormat compiles a format string that is known good at build time.
func mustParseFormat(raw string) logFormat {
	format, err := parseFormat(raw)
	if err != nil {
		panic(err) // programmer error: DefaultFormat is a constant
	}
	return format
}

// renderTitle draws the title bar.
func (t *Tracer) renderTitle() string {
	title := fmt.Sprintf("%s — %d logs", t.app, len(t.Visible()))
	return titleStyle.Render(title) + "\n" + t.rule()
}

// renderLogs draws the scrolling log pane, sized to the remaining height.
func (t *Tracer) renderLogs() string {
	visible := t.Visible()
	if len(visible) == 0 {
		return hintStyle.Render("no logs yet")
	}

	rows := t.logRows()
	start := t.windowStart(len(visible), rows)
	end := start + rows
	if end > len(visible) {
		end = len(visible)
	}

	lines := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		line := t.format.render(visible[i])
		if i == t.cursor {
			line = selectedStyle.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// renderFocused draws the focused log in full, beneath the log pane.
func (t *Tracer) renderFocused() string {
	visible := t.Visible()
	if t.focused < 0 || t.focused >= len(visible) {
		return ""
	}
	body := t.format.render(visible[t.focused])
	return t.rule() + "\n" + focusedStyle.Render(body)
}

// renderFilter draws the one-line needle editor.
func (t *Tracer) renderFilter() string {
	cursor := "_"
	if t.needle == "" {
		cursor = hintStyle.Render("type a needle, / for regex")
	}
	return t.rule() + "\n" + needleStyle.Render("filter: ") + t.needle + cursor
}

// renderCommands draws the key hints.
func (t *Tracer) renderCommands() string {
	hints := []string{"[Q]uit", "[F]ilter", "[C]lear"}
	if t.filterUI {
		hints[1] = "[Ctrl-F]ilter"
	}
	return t.rule() + "\n" + hintStyle.Render(strings.Join(hints, "  "))
}

// rule draws a horizontal separator across the terminal width.
func (t *Tracer) rule() string {
	width := t.width
	if width < 1 {
		width = defaultWidth
	}
	return ruleStyle.Render(strings.Repeat("─", width))
}

// logRows returns how many log lines fit between the title and the footer.
func (t *Tracer) logRows() int {
	height := t.height
	if height < 1 {
		height = defaultHeight
	}
	rows := height - 4 // title, rule, rule, commands
	if t.focused >= 0 {
		rows -= 3 // rule, focused log, border
	}
	if t.filterUI {
		rows -= 2 // rule, needle
	}
	if rows < 1 {
		return 1
	}
	return rows
}

// windowStart returns the first visible index so the cursor stays on screen.
func (t *Tracer) windowStart(total, rows int) int {
	if total <= rows {
		return 0
	}
	start := t.cursor - rows + 1
	if start < 0 {
		return 0
	}
	if start > total-rows {
		return total - rows
	}
	return start
}
