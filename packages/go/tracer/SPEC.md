similar to https://github.com/joshka/tui-tracing but very feature reduced for now.

# Terms
`needle` a string which is used for filtering: if starts with `/` it should be used as a regex otherwise if any of the fields in an object/row include the value then that object/row is understood to be "included". examples:
value | value | included
"jimbo loves to thunks in bed" | "unk", "^jimbo*", " bed", "jim" | TRUE
"jimbo loves to junkies in bed" | "asdf, "^junk",, " bedri", "john" | FALSE

# TUI - full screen utility application run in bash
## Formatting:
keep it really simple %t = time, %l = loglevel, %m = message so:
"%t %l %m" = "20:23:33.199 warn this is a warning"
"%l | %t | %m" = "warn | 20:23:33.199 | this is a warning"
each %option can have a trailing char they are:
"U" (so %lU makes the message appear in full caps) and:
`%t %l %mU" = "20:23:33.199 warn THIS IS A WARNING"

## Sections:

### <log>:
a single log line:

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
