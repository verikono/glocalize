package types

type NotatedCall struct {
	Line     int
	FilePath string
	Function string
}

type CallStack struct {
	Calls []NotatedCall
}

type (
	EnvMode string
)

const (
	ProdEnvMode    EnvMode = "production"
	DevEnvMode     EnvMode = "development"
	StagingEnvMode EnvMode = "staging"
	TestingEnvMode EnvMode = "testing"
)

type Service interface {
	Start() error
}

type TransactionId string

type ErrorMap = map[TransactionId]error

type SmtpTransactionId = TransactionId

type SmtpErrorMap = map[SmtpTransactionId]error
