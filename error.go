package oops

import (
	"log/slog"
	"maps"
	"slices"
)

type Error struct {
	err        error
	attributes map[string]any
	sources    []Source
	code       int
}

// Error returns undecorated error.
func (e *Error) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// Unwrap allows errors.As and errors.Is to find the inner error.
func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) LogValue() slog.Value {
	var attrs []slog.Attr
	if err := e.Error(); err != "" {
		attrs = append(attrs, slog.String("err", err))
	}

	attrs = append(attrs, e.LogAttrs()...)

	return slog.GroupValue(attrs...)
}

func (e *Error) LogAttrs() []slog.Attr {
	var attrs []slog.Attr
	for _, k := range slices.Sorted(maps.Keys(e.attributes)) {
		attrs = append(attrs, slog.Any(k, e.attributes[k]))
	}

	if len(e.sources) > 0 {
		attrs = append(attrs, slog.Any("source", e.sources))
	}
	return attrs
}

// With collects key-value pairs to return with the error.
// supports supports [string, any]... pairs or slog.Attr values.
func (e *Error) With(args ...any) *Error {
	e.attributes = argsToAttr(e.attributes, args)
	return e
}

// Code sets the error code.
func (e *Error) Code(code int) *Error {
	e.code = code
	return e
}

// argsToAttr recursively builds the attributes map from the args slice
// supports supports [string, any]... pairs or slog.Attr values.
// bad pairs are skipped.
func argsToAttr(attrs map[string]any, args []any) map[string]any {
	if len(args) == 0 {
		return attrs
	}
	if attrs == nil {
		attrs = make(map[string]any)
	}

	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			// skipping bad pair
			return attrs
		}
		attrs[x] = args[1]
		return argsToAttr(attrs, args[2:])

	case slog.Attr:
		attrs[x.Key] = x.Value.Any()
		return argsToAttr(attrs, args[1:])

	default:
		// skipping bad pair
		return argsToAttr(attrs, args[1:])
	}
}
