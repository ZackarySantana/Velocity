package logger

import "errors"

type Collect struct {
	infos    []error
	warnings []error
	errs     []error
}

func NewCollectLogger() Logger {
	return &Collect{}
}

func (l *Collect) Write(p []byte) (int, error) {
	l.infos = append(l.infos, errors.New(string(p)))
	return len(p), nil
}

func (l *Collect) WrapInfo(info string) {
	if info == "" {
		return
	}
	l.infos = append(l.infos, errors.New(info))
}

func (l *Collect) Info(info error) {
	if info == nil {
		return
	}
	l.infos = append(l.infos, info)
}

func (l *Collect) WrapWarning(warning string) {
	if warning == "" {
		return
	}
	l.warnings = append(l.warnings, errors.New(warning))
}

func (l *Collect) Warning(warning error) {
	if warning == nil {
		return
	}
	l.warnings = append(l.warnings, warning)
}

func (l *Collect) WrapError(err string) {
	if err == "" {
		return
	}
	l.errs = append(l.errs, errors.New(err))
}

func (l *Collect) Error(err error) {
	if err == nil {
		return
	}
	l.errs = append(l.errs, err)
}

func (l *Collect) Output(level string) error {
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
