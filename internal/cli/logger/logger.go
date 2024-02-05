package logger

type Level int

const (
	Info Level = iota
	Warning
	Error
)

var (
	infoPrefix    = []byte("[INFO] ")
	warningPrefix = []byte("[WARNING] ")
	errorPrefix   = []byte("[ERROR] ")
	newLine       = []byte("\n")
)

// TODO: change logger in to an interface and implement 'CollectLogger' and 'LiveLogger' types
type Logger interface {
	Write([]byte) (int, error)
	Info([]byte)
	InfoStr(string)
	InfoErr(error)
	Warning([]byte)
	WarningStr(string)
	WarningErr(error)
	Error([]byte)
	ErrorStr(string)
	ErrorErr(error)
}
