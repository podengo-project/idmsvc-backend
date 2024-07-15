package logger

import (
	"context"
	"log/slog"
)

// type SlogMetaHandler slog.Handler
type SlogMetaHandler interface {
	slog.Handler
	Add(handler slog.Handler)
}

type slogMetaHandler struct {
	handlers []slog.Handler
}

func NewSlogMetaHandler() SlogMetaHandler {
	return &slogMetaHandler{
		handlers: []slog.Handler{},
	}
}

// Add a new handler to the collection
func (h *slogMetaHandler) Add(handler slog.Handler) {
	h.handlers = append(h.handlers, handler)
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *slogMetaHandler) Enabled(ctx context.Context, level slog.Level) bool {
	ret := true
	for i := range h.handlers {
		if ok := h.handlers[i].Enabled(ctx, level); !ok {
			ret = ok
		}
	}
	return ret
}

// WithAttrs returns a new MetaHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *slogMetaHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newMeta := &slogMetaHandler{handlers: make([]slog.Handler, len(h.handlers))}
	for i := range h.handlers {
		newMeta.handlers[i] = h.handlers[i].WithAttrs(attrs)
	}
	return newMeta
}

// WithGroup add a new group to the structured log.
func (h *slogMetaHandler) WithGroup(name string) slog.Handler {
	newMeta := &slogMetaHandler{handlers: make([]slog.Handler, len(h.handlers))}
	for i := range h.handlers {
		newMeta.handlers[i] = h.handlers[i].WithGroup(name)
	}
	return newMeta
}

// Handle formats its argument Record as a single line of space-separated
// key=value items.
func (h *slogMetaHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for i := range h.handlers {
		if err2 := h.handlers[i].Handle(ctx, r); err2 != nil {
			err = err2
		}
	}
	return err
}
