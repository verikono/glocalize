similar to https://github.com/joshka/tui-tracing but very feature reduced for now.

# Terms
`needle` a string which is used for filtering: if starts with `/` it should be used as a regex otherwise if any of the fields in an object/row include the value then that object/row is understood to be "included". examples:
value | value | included
"jimbo loves to thunks in bed" | "unk", "^jimbo*", " bed", "jim" | TRUE
"jimbo loves to junkies in bed" | "asdf, "^junk",, " bedri", "john" | FALSE

`origin` an identifier of the package which produced the log entry, rendered as a bracketed token: `[orm] sql query`. A bracketed token that names a level (`[ERROR]`) or an operation (`[WRITE]:`) is not an origin.

`message` the concise, succinct line a reader scans.

`details` an optional list of extra lines a package may produce in any quantity. A detail is written as an indented line beneath its log and is shown only when the log is focused.

> **Implementation note.** The prose above and the table disagree: `^jimbo*` is unprefixed, yet it only matches under regex semantics. The table is the acceptance criterion, so **every** needle is compiled as a case-insensitive regex; a leading `/` is stripped first, and a needle that will not compile falls back to a literal substring match. The needle is matched against every field of a log: timestamp, level, origin, message and each detail line.

# Logger
`tracer.NewLogger(origin, out)` writes tracer-native lines for one package. Each entry is a header line followed by one indented line per detail:

```
2026/10/21 14:13:40.709 trace [orm] sql query
	SELECT * FROM "Project" WHERE "name" = ?
	args: [all]
```

A message containing newlines is flattened onto one line, because a message is by definition concise and a newline would break the header. A detail containing newlines is split, so every physical line stays indented and parses back as a detail. A `Logger` is safe for concurrent use.

# TUI - full screen utility application run in bash
## Formatting:
keep it really simple %t = time, %l = loglevel, %m = message, %o = origin so:
"%t %l %m" = "20:23:33.199 warn this is a warning"
"%l | %t | %m" = "warn | 20:23:33.199 | this is a warning"
"%t %l [%o] %m" = "20:23:33.199 warn [orm] this is a warning"
each %option can have a trailing char they are:
"U" (so %lU makes the message appear in full caps) and:
`%t %l %mU" = "20:23:33.199 warn THIS IS A WARNING"

## Sections:

### <log>:
a single log line:

### <focused-log>:
the selected log rendered in full, followed by one line per detail it carries.

### <commands>
[Q]uit - quits the app
[F]ilter - see <filtering-ui> , changes to [Ctrl-F]ilter when active <filtering-ui>
[C]lear - clears the filter

### <filtering-ui>
not visible until user his [F]
1 line high, lets the user provide a needle
Logs Section:
up and down arrows navigate to log items
enter: focuses a selected log
backspace: (if focused log) clears the focus 
x: an alias of backspace

Filter

Interface:
---
    ---
        {title}
    ---
        {<log>} 
        {<log>}
        {<log>}
        ... more logs ... (flexible to full screen)
    ---
        {<focused-log> (condition: if a focused log is set)}
    ---
        {filter (condition: if filtering filtering)}
    ---
        commands
    ---
---
