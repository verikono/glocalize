---
name: tracer
description: Emit and read log entries with the glocalize tracer package — how a package provides a log entry (origin, level, concise message, optional details), how the line is written and parsed back, and how the TUI renders and filters it. Load when adding logging to a Go package, when a log line is not appearing or is attributed to the wrong origin, or when changing the tracer itself.
---

# tracer

`packages/go/tracer` is the log format and the full-screen log viewer for this
monorepo. A package does not print free text: it hands the tracer a **log
entry** with an origin, a level, a concise message and optional details. The
tracer writes that entry as one line (plus indented detail lines), parses it
back, and renders it.

The contract lives in `types/go/tracer.go` (`TracerLogInterface`); the
implementation lives in `packages/go/tracer`.

## The entry

| field | meaning |
| --- | --- |
| `origin` | the **package name** that produced the entry, e.g. `orm`, `mcp`, `utils` |
| `time` | when it happened |
| `level` | `trace`, `info`, `warn`, `error`, `fatal` |
| `message` | the concise, succinct line a reader scans |
| `details` | optional list of extra lines — the waffle, in any quantity |

Two message types, deliberately: **message** is what you scan, **details** is
what you read once you have decided the entry matters. Details are rendered
only when the entry is focused in the viewer, so a package may be as verbose as
it likes without drowning the log pane.

## Providing a log entry from a package

Use `tracer.Logger`. It stamps the origin on every line, so the caller never
repeats it.

```go
package orm

import (
	"os"

	"github.com/verikono/glocalize/packages/go/tracer"
)

// Origin is the package name stamped on every entry this package produces.
const Origin = "orm"

// ormLog writes to stdout, which is where `riot serve` redirects the daemon.
var ormLog = tracer.NewLogger(Origin, os.Stdout)

func (q *Queries) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	ormLog.Trace("sql query", query, fmt.Sprintf("args: %v", args))
	// ...
	if err != nil {
		ormLog.Error("sql error", err.Error(), query, fmt.Sprintf("args: %v", args))
		return nil, err
	}
	return res, nil
}
```

Rules that follow from the format:

- **One origin per package.** Declare `const Origin = "<package>"` and build one
  package-level `Logger`. Do not pass an origin per call — the origin is a
  property of the producer, not of the event.
- **The message is one line.** A newline in the message is flattened to a space
  by the logger, because a newline would break the header. If you have more to
  say, say it in a detail.
- **Details are variadic and may be many.** Each detail is written as an
  indented line. A detail containing newlines is split, so every physical line
  stays a detail.
- **`Fatal` does not exit.** It writes a fatal-level entry; the caller owns the
  decision to stop.
- **A `Logger` is safe for concurrent use.** One package-level value is correct.
- **A nil writer means stdout.** `tracer.NewLogger(origin, nil)` writes to
  `os.Stdout`, which is where the daemon's output is redirected.

### Levels

`Trace`, `Info`, `Warn`, `Error`, `Fatal`. Use `Trace` for the high-volume
mechanical detail (every SQL statement), `Info` for lifecycle events, `Warn`
for something recovered from, `Error` for a failure being returned.

## The wire format

A `Logger` writes:

```
2026/10/20 14:13:38.825 warn [orm] slow query
	SELECT * FROM Agent
	took 900ms
```

- header: `<timestamp> <level> [<origin>] <message>`
- each detail: a tab-indented line beneath the header

The parser accepts the timestamp in any layout riot or the standard library
emits, with or without brackets, and accepts the level and origin in either
order. A bracketed token that names a level (`[ERROR]`) or an operation
(`[WRITE]:`) is **not** read as an origin.

## Reading entries back

```go
log := tracer.ParseLogLine(line, time.Now())   // one line, no details
logs := tracer.ParseLogLines(lines, time.Now()) // a block; indented lines become details
```

Both return `*tracer.TracerLog`, which satisfies `TracerLogInterface`:
`GetOrigin`, `GetTimestamp`, `GetTime`, `GetLevel`, `GetMessage`, `GetDetails`.

To build one by hand (a test, a synthetic entry):

```go
log := tracer.NewTracerLogFrom("orm", ts, tracer.WarnLogLevel, "slow query", "SELECT 1", "took 900ms")
```

## Rendering

`RenderFormat(raw, log)` applies a format string. Options:

| option | renders |
| --- | --- |
| `%t` | time, `15:04:05.000` |
| `%l` | level |
| `%m` | message |
| `%o` | origin |

Any option takes a trailing `U` for upper case (`%oU` → `ORM`). `%%` is a
literal percent. An unknown option is an error, so a bad format fails at
construction rather than at render time.

`DefaultFormat` is `"%t %l %m"`; `DefaultFormatWithOrigin` is `"%t %l [%o] %m"`.
Riot's `trace` command uses the latter.

## Filtering

`TracerLogFilter` narrows by time window, level and needle. The needle is
matched case-insensitively against **every** field: timestamp, level, origin,
message and each detail line. A needle beginning with `/` is compiled as a
regular expression; one that will not compile degrades to a literal substring
match rather than erroring.

## The viewer

`tracer.NewTracer(tracer.WithSource(src), tracer.WithApp("riot"))` builds the
full-screen viewer; `Run(ctx)` blocks until the user quits and owns the context
the source is followed under. A source is anything satisfying
`TracerLogSource` (`Follow(ctx) (<-chan string, error)`); `NewFileSource` and
`NewFileSourceFrom` follow a file like `tail -f`.

Keys: `q` quit, `f` filter, `c` clear, `↑`/`↓` (or `k`/`j`) navigate, `enter`
focus, `backspace`/`x` unfocus.

## Checklist when adding logging to a package

- [ ] `const Origin = "<package>"` declared in the package
- [ ] one package-level `tracer.NewLogger(Origin, os.Stdout)`
- [ ] the message is concise and single-line
- [ ] anything verbose goes in a detail, not the message
- [ ] the level matches the severity, not the mood
- [ ] no `log.Printf` left behind — a stdlib line has no origin and cannot be
      attributed or filtered by package
