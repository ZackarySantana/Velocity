package logger

import (
	"io"
)

type Live struct {
	info    []io.Writer
	warning []io.Writer
	err     []io.Writer
}

func NewLiveLogger() *Live {
	return &Live{}
}

func (l *Live) SubscribeInfo(writer io.Writer) {
	l.info = append(l.info, writer)
}

func (l *Live) SubscribeWarning(writer io.Writer) {
	l.warning = append(l.warning, writer)
}

func (l *Live) SubscribeError(writer io.Writer) {
	l.err = append(l.err, writer)
}

func (l *Live) Write(p []byte) (int, error) {
	for _, writer := range l.info {
		writer.Write(p)
	}
	for _, writer := range l.warning {
		writer.Write(p)
	}
	for _, writer := range l.err {
		writer.Write(p)
	}
	return len(p), nil
}

func (l *Live) Info(info []byte) {
	writeTo(l.info, infoPrefix, info, newLine)
	writeTo(l.warning, infoPrefix, info, newLine)
	writeTo(l.err, infoPrefix, info, newLine)
}

func (l *Live) InfoStr(info string) {
	if info == "" {
		return
	}
	l.Info([]byte(info))
}

func (l *Live) InfoErr(info error) {
	if info == nil {
		return
	}
	l.Info([]byte(info.Error()))
}

func (l *Live) Warning(warning []byte) {
	writeTo(l.warning, warningPrefix, warning, newLine)
	writeTo(l.err, warningPrefix, warning, newLine)
}

func (l *Live) WarningStr(warning string) {
	if warning == "" {
		return
	}
	l.Warning([]byte(warning))
}

func (l *Live) WarningErr(warning error) {
	if warning == nil {
		return
	}
	l.Warning([]byte(warning.Error()))
}

func (l *Live) Error(err []byte) {
	writeTo(l.err, errorPrefix, err, newLine)
}

func (l *Live) ErrorStr(err string) {
	if err == "" {
		return
	}
	l.Error([]byte(err))
}

func (l *Live) ErrorErr(err error) {
	if err == nil {
		return
	}
	l.Error([]byte(err.Error()))
}

func writeTo(writers []io.Writer, p ...[]byte) {
	for _, writer := range writers {
		for _, data := range p {
			writer.Write(data)
		}
	}
}
