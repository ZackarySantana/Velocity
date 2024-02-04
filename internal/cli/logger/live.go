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

func (l *Live) Write(p []byte) (int, error) {
	for _, writer := range l.info {
		writer.Write(p)
	}
	return len(p), nil
}

func (l *Live) SubscribeInfo(writer io.Writer) {
	l.info = append(l.info, writer)
}

func (l *Live) SubscribeWarning(writer io.Writer) {
	l.info = append(l.info, writer)
	l.warning = append(l.warning, writer)
}

func (l *Live) SubscribeError(writer io.Writer) {
	l.info = append(l.info, writer)
	l.warning = append(l.warning, writer)
	l.err = append(l.err, writer)
}

func (l *Live) WrapInfo(info string) {
	l.sendInfo(info)
}

func (l *Live) Info(info error) {
	if info != nil {
		l.sendInfo(info.Error())
	}
}

func (l *Live) sendInfo(info string) {
	if info == "" {
		return
	}
	for _, writer := range l.info {
		writer.Write([]byte("[INFO] " + info + "\n"))
	}
}

func (l *Live) WrapWarning(warning string) {
	l.sendWarning(warning)
}

func (l *Live) Warning(warning error) {
	if warning != nil {
		l.sendWarning(warning.Error())
	}
}

func (l *Live) sendWarning(warning string) {
	if warning == "" {
		return
	}
	for _, writer := range l.warning {
		writer.Write([]byte("[WARNING] " + warning + "\n"))
	}
}

func (l *Live) WrapError(err string) {
	l.sendError(err)
}

func (l *Live) Error(err error) {
	if err != nil {
		l.sendError(err.Error())
	}
}

func (l *Live) sendError(err string) {
	if err == "" {
		return
	}
	for _, writer := range l.err {
		writer.Write([]byte("[ERROR] " + err + "\n"))
	}
}
