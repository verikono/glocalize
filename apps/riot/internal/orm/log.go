package riot

import (
	"os"

	"github.com/verikono/glocalize/packages/go/tracer"
)

// Origin is the package name stamped on every log entry the orm produces, so a
// tracer can attribute a SQL line to the data layer.
const Origin = "orm"

// ormLog is the orm's tracer-native logger. It writes to stdout, which is where
// `riot serve` redirects the daemon's output.
var ormLog = tracer.NewLogger(Origin, os.Stdout)
