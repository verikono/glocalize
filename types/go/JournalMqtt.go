package types

type MqttJournalHandlerPackage struct {
	NewMqttJournalHandler func(options ...OptionerFn) (JournalHandler, error)
}
