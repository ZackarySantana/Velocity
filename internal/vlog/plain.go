package vlog

import (
	"context"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type PlainHandler struct {
	opts Options
	out  io.Writer
}

type Options struct {
	Level slog.Leveler
}

func NewPlainHandler(out io.Writer, opts *Options) *PlainHandler {
	h := &PlainHandler{out: out}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

func (h *PlainHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PlainHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	var err error
	switch r.Level {
	case slog.LevelInfo:
		buf = append(buf, r.Message...)
	case slog.LevelDebug:
		buf = append(buf, color.CyanString(r.Message)...)
	case slog.LevelWarn:
		buf = append(buf, color.YellowString(r.Message)...)
	case slog.LevelError:
		buf = append(buf, color.RedString(r.Message)...)
	}

	buf = append(buf, '\n')
	_, err = h.out.Write(buf)
	return err
}

func (h *PlainHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *PlainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}
