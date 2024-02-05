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

func (l *Collect) Info(info []byte) {
	l.infos = append(l.infos, errors.New(string(infoPrefix)+string(info)))
	l.warnings = append(l.warnings, errors.New(string(infoPrefix)+string(info)))
	l.errs = append(l.errs, errors.New(string(infoPrefix)+string(info)))
}

func (l *Collect) InfoStr(info string) {
	if info == "" {
		return
	}
	l.Info([]byte(info))
}

func (l *Collect) InfoErr(info error) {
	if info == nil {
		return
	}
	l.Info([]byte(info.Error()))
}

func (l *Collect) Warning(warning []byte) {
	l.warnings = append(l.warnings, errors.New(string(warningPrefix)+string(warning)))
	l.errs = append(l.errs, errors.New(string(warningPrefix)+string(warning)))
}

func (l *Collect) WarningStr(warning string) {
	if warning == "" {
		return
	}
	l.Warning([]byte(warning))
}

func (l *Collect) WarningErr(warning error) {
	if warning == nil {
		return
	}
	l.Warning([]byte(warning.Error()))
}

func (l *Collect) Error(err []byte) {
	l.errs = append(l.errs, errors.New(string(errorPrefix)+string(err)))
}

func (l *Collect) ErrorStr(err string) {
	if err == "" {
		return
	}
	l.Error([]byte(err))
}

func (l *Collect) ErrorErr(err error) {
	if err == nil {
		return
	}
	l.Error([]byte(err.Error()))
}

func (l *Collect) Output(level Level) error {
	if level == Info {
		return errors.Join(l.infos...)
	}
	if level == Warning {
		return errors.Join(l.warnings...)
	}
	if level == Error {
		return errors.Join(l.errs...)
	}
	return nil
}
