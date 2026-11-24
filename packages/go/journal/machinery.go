package journal

var JournalDict = make(map[JournalName]Journal)

func journalExists(name JournalName) bool {
	_, exists := JournalDict[name]
	return exists
}
