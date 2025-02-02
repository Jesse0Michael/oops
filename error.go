package oops

import (
	"log/slog"
)

type Error struct {
	err        error
	attributes map[string]any
	sources    []Source
}

// Error returns undecorated error.
func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) LogValue() slog.Value {
	var attrs []slog.Attr
	if err := e.Error(); err != "" {
		attrs = append(attrs, slog.String("err", err))
	}

	for k, v := range e.attributes {
		attrs = append(attrs, slog.Any(k, v))
	}

	sources := make([]slog.Value, len(e.sources))
	for i, s := range e.sources {
		sources[i] = slog.GroupValue(
			slog.String("file", s.File),
			slog.String("function", s.Function),
			slog.Int("line", s.Line),
		)
	}
	if len(sources) == 1 {
		attrs = append(attrs, slog.Any("source", sources[0]))
	} else if len(sources) > 0 {
		attrs = append(attrs, slog.Any("sources", sources))
	}

	return slog.GroupValue(attrs...)
}

// With collects key-value pairs to return with the error.
// supports supports [string, any]... pairs or slog.Attr values.
func (e Error) With(args ...any) Error {
	e.attributes = argsToAttr(e.attributes, args)
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
