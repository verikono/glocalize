package errantly

type ErrantlyWriterName string
type ErrantlyWriter interface {
	Write(error ErrantlyError) error
}

var writers map[ErrantlyWriterName][]ErrantlyWriter

// registration handler for errantly writers
func registerWriter(name ErrantlyWriterName, writer ErrantlyWriter) error {
	return NotImplemented()
}

// send the error to the Writer battery for their purposes.
func SendToWriters(err ErrantlyError) error {
	return NotImplemented()
	// for _, w := range writers {
	// w.Write()
	// }
}
