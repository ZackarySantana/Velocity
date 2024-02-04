package logger

import "errors"

type Level int

const (
	Info Level = iota
	Warning
	Error
)

// TODO: change logger in to an interface and implement 'CollectLogger' and 'LiveLogger' types
type Logger struct {
	infos    []error
	warnings []error
	errs     []error
}

func NewLogger() Logger {
	return Logger{}
}

func (l *Logger) WrapInfo(info string) {
	if info == "" {
		return
	}
	l.infos = append(l.infos, errors.New(info))
}

func (l *Logger) Info(info error) {
	if info == nil {
		return
	}
	l.infos = append(l.infos, info)
}

func (l *Logger) WrapWarning(warning string) {
	if warning == "" {
		return
	}
	l.warnings = append(l.warnings, errors.New(warning))
}

func (l *Logger) Warning(warning error) {
	if warning == nil {
		return
	}
	l.warnings = append(l.warnings, warning)
}

func (l *Logger) WrapError(err string) {
	if err == "" {
		return
	}
	l.errs = append(l.errs, errors.New(err))
}

func (l *Logger) Error(err error) {
	if err == nil {
		return
	}
	l.errs = append(l.errs, err)
}

func (l *Logger) Output(level string) error {
	length := len(l.errs)
	if level == "warning" || level == "info" {
		length += len(l.warnings)
	}
	if level == "info" {
		length += len(l.infos)
	}

	all := make([]error, length)
	all = append(all, l.errs...)
	if level == "warning" || level == "info" {
		all = append(all, l.warnings...)
	}
	if level == "info" {
		all = append(all, l.infos...)
	}

	return errors.Join(all...)

}
