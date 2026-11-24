package mqtt

const MqttJournalHandlerKind JournalHandlerKind = "mqtt"

type MqttJournalHandler struct {
	kind JournalHandlerKind
}

func NewMqttJournalHandler(options ...OptionerFn) (JournalHandler, error) {
	inst := MqttJournalHandler{
		kind: JournalHandlerKind(MqttJournalHandlerKind),
	}
	for _, option := range options {
		if err := option(inst); err != nil {
			return nil, err
		}
	}
	return &inst, nil
}

func (mq *MqttJournalHandler) GetKind() JournalHandlerKind {
	return MqttJournalHandlerKind
}

func (mq *MqttJournalHandler) ProcessEntry(entry JournalEntryProcessor) error {
	return E.NotImplemented()
}
