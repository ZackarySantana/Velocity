package vlog

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/fatih/color"
)

// Options defines configuration for the PlainHandler.
type Options struct {
	// Level is the minimum level of log to output.
	Level slog.Level
	// TimeFormat is the format for timestamps. If empty, RFC3339 is used.
	TimeFormat string
}

// PlainHandler is a user-friendly CLI slog handler.
type PlainHandler struct {
	out   io.Writer
	opts  *Options
	mu    *sync.Mutex
	group []string
	attrs []slog.Attr
}

// NewPlainHandler constructs a new PlainHandler.
func NewPlainHandler(out io.Writer, opts *Options) *PlainHandler {
	if opts == nil {
		t := Options{Level: slog.LevelInfo, TimeFormat: time.RFC3339}
		opts = &t
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}
	return &PlainHandler{
		out:  out,
		opts: opts,
		mu:   &sync.Mutex{},
	}
}

// Enabled reports whether the Handler handles records at the given level.
func (h *PlainHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

// Handle formats and writes the Record to the output.
func (h *PlainHandler) Handle(ctx context.Context, rec slog.Record) error {
	// Ensure only one write at a time
	h.mu.Lock()
	defer h.mu.Unlock()

	// Timestamp
	timestamp := rec.Time.Format(h.opts.TimeFormat)

	// Level
	levelText := rec.Level.String()
	switch rec.Level {
	case slog.LevelDebug:
		levelText = color.New(color.FgCyan).Sprint(levelText)
	case slog.LevelInfo:
		levelText = color.New(color.FgGreen).Sprint(levelText)
	case slog.LevelWarn:
		levelText = color.New(color.FgYellow).Sprint(levelText)
	case slog.LevelError:
		levelText = color.New(color.FgRed).Sprint(levelText)
	}

	// Group prefix (joined by ".")
	groupPrefix := ""
	if len(h.group) > 0 {
		groupPrefix = "[" + strings.Join(h.group, ".") + "] "
	}

	// Message
	msg := rec.Message

	// Collect and format attributes
	var parts []string
	for _, a := range h.attrs {
		parts = append(parts, formatAttr(a))
	}
	rec.Attrs(func(a slog.Attr) bool {
		parts = append(parts, formatAttr(a))
		return true
	})

	// Build line
	line := timestamp + " " + levelText + " " + groupPrefix + msg
	if len(parts) > 0 {
		line += " " + strings.Join(parts, " ")
	}

	// Write out
	_, err := io.WriteString(h.out, line+"\n")
	return err
}

// WithAttrs returns a new handler whose records will include the given attributes.
func (h *PlainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	r := &PlainHandler{
		out:   h.out,
		opts:  h.opts,
		mu:    h.mu,
		group: append([]string{}, h.group...),
		attrs: append([]slog.Attr{}, h.attrs...),
	}
	r.attrs = append(r.attrs, attrs...)
	return r
}

// WithGroup returns a new handler whose records will include the given group.
func (h *PlainHandler) WithGroup(name string) slog.Handler {
	r := &PlainHandler{
		out:   h.out,
		opts:  h.opts,
		mu:    h.mu,
		group: append([]string{}, h.group...),
		attrs: append([]slog.Attr{}, h.attrs...),
	}
	r.group = append(r.group, name)
	return r
}

// formatAttr formats an Attr as key=value
func formatAttr(a slog.Attr) string {
	return a.Key + "=" + a.Value.String()
}
