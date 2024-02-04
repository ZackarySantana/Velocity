package logger

type Level int

const (
	Info Level = iota
	Warning
	Error
)

// TODO: change logger in to an interface and implement 'CollectLogger' and 'LiveLogger' types
type Logger interface {
	WrapInfo(string)
	Info(error)
	WrapWarning(string)
	Warning(error)
	WrapError(string)
	Error(error)
}
