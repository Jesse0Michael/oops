package oops

import (
	"context"
	"errors"
	"log/slog"
)

// attrField is the field used to log oops error attributes.
var attrField = "oops"

// Handler logs oops errors with additional context.
type Handler struct {
	slog.Handler
}

// NewOopsHandler creates a handler that logs oops errors with additional context.
func NewOopsHandler(handler slog.Handler) slog.Handler {
	return &Handler{
		Handler: handler,
	}
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return &Handler{Handler: h.Handler.WithGroup(name)}
}

// Handle implements slog.Handler.
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	record.Attrs(func(attr slog.Attr) bool {
		if err, ok := attr.Value.Any().(error); ok {
			var oops *Error
			if errors.As(err, &oops) {
				newRecord.AddAttrs(slog.Any(attr.Key, oops.err))
				if attrs := oops.LogAttrs(); len(attrs) > 0 {
					newRecord.AddAttrs(slog.Any(attrField, attrs))
				}
			} else {
				newRecord.AddAttrs(attr)
			}
			return true
		}
		newRecord.AddAttrs(attr)
		return true
	})

	return h.Handler.Handle(ctx, newRecord)
}

// SetAttrField overrides the field used to oops error attributes.
// default is "oops".
func SetAttrField(field string) {
	attrField = field
}
