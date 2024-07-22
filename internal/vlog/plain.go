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

	buf = append(buf, r.Message...)
	r.Attrs(func(a slog.Attr) bool {
		buf = append(buf, ": "+a.Value.String()...)
		return true
	})
	buf = append(buf, '\n')

	var err error
	switch r.Level {
	case slog.LevelDebug:
		buf = []byte(color.CyanString(string(buf)))
	case slog.LevelWarn:
		buf = []byte(color.YellowString(string(buf)))
	case slog.LevelError:
		buf = []byte(color.RedString(string(buf)))
	}

	_, err = h.out.Write(buf)
	return err
}

func (h *PlainHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *PlainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}
