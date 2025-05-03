package vlog

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type YAMLHandler struct {
	w      io.Writer
	opts   *Options
	mu     *sync.Mutex
	groups []string
	attrs  []slog.Attr
}

func NewYAMLHandler(w io.Writer, opts *Options) *YAMLHandler {
	if opts == nil {
		t := Options{Level: slog.LevelInfo, TimeFormat: time.RFC3339}
		opts = &t
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}
	return &YAMLHandler{
		w:    w,
		opts: opts,
		mu:   &sync.Mutex{},
	}
}

func (h *YAMLHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

func (h *YAMLHandler) Handle(_ context.Context, r slog.Record) error {
	attrsMap := collectAttrs(h.attrs, r)
	entry := map[string]any{r.Message: attrsMap}

	output, err := yaml.Marshal(entry)
	if err != nil {
		return err
	}

	var coloredOutput string
	switch r.Level {
	case slog.LevelDebug:
		coloredOutput = color.New(color.FgBlue).Sprint(string(output))
	case slog.LevelWarn:
		coloredOutput = color.New(color.FgYellow).Sprint(string(output))
	case slog.LevelError:
		coloredOutput = color.New(color.FgRed).Sprint(string(output))
	default:
		coloredOutput = string(output)
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.w.Write([]byte(coloredOutput))
	return err
}

func (h *YAMLHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := *h
	newH.attrs = append(append([]slog.Attr(nil), h.attrs...), attrs...)
	return &newH
}

func (h *YAMLHandler) WithGroup(name string) slog.Handler {
	newH := *h
	newH.groups = append(append([]string(nil), h.groups...), name)
	return &newH
}

func collectAttrs(baseAttrs []slog.Attr, r slog.Record) map[string]any {
	attrsMap := make(map[string]any)

	// Inline attribute collection and mapping
	addAttr := func(a slog.Attr) {
		if a.Equal(slog.Attr{}) {
			return
		}
		attrsMap[a.Key] = resolveValue(a.Value)
	}

	for _, a := range baseAttrs {
		addAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		addAttr(a)
		return true
	})

	return attrsMap
}

func resolveValue(v slog.Value) any {
	switch v.Kind() {
	case slog.KindGroup:
		m := make(map[string]any)
		for _, nested := range v.Group() {
			if !nested.Equal(slog.Attr{}) {
				m[nested.Key] = resolveValue(nested.Value)
			}
		}
		return m
	case slog.KindLogValuer:
		return resolveValue(v.LogValuer().LogValue())
	case slog.KindString:
		return v.String()
	case slog.KindTime:
		return v.Time()
	default:
		return v.Any()
	}
}
