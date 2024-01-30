package logger

import "errors"

type Level int

const (
	Info Level = iota
	Warning
	Error
)

type logger struct {
	infos    []error
	warnings []error
	errs     []error
}

func NewLogger() logger {
	return logger{}
}

func (l *logger) WrapInfo(info string) {
	if info == "" {
		return
	}
	l.infos = append(l.infos, errors.New(info))
}

func (l *logger) Info(info error) {
	if info == nil {
		return
	}
	l.infos = append(l.infos, info)
}

func (l *logger) WrapWarning(warning string) {
	if warning == "" {
		return
	}
	l.warnings = append(l.warnings, errors.New(warning))
}

func (l *logger) Warning(warning error) {
	if warning == nil {
		return
	}
	l.warnings = append(l.warnings, warning)
}

func (l *logger) WrapError(err string) {
	if err == "" {
		return
	}
	l.errs = append(l.errs, errors.New(err))
}

func (l *logger) Error(err error) {
	if err == nil {
		return
	}
	l.errs = append(l.errs, err)
}

func (l *logger) Output(level string) error {
	length := len(l.errs)
	if level == "warning" || level == "info" {
		length += len(l.warnings)
	}
	if level == "info" {
		length += len(l.infos)
	}

	all := make([]error, length)
	for _, err := range l.errs {
		all = append(all, err)
	}
	if level == "warning" || level == "info" {
		for _, warning := range l.warnings {
			all = append(all, warning)
		}
	}
	if level == "info" {
		for _, info := range l.infos {
			all = append(all, info)
		}
	}

	return errors.Join(all...)

}
