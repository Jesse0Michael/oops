package oops

import (
	"errors"
	"fmt"
)

// sourceEnabled is a package variable to enable/disable tracking the source location.
var sourceEnabled = true

// New returns a new oops error.
func New(msg string) Error {
	var sources []Source
	if s := source(); s != nil {
		sources = append(sources, *s)
	}

	return Error{err: errors.New(msg), sources: sources}
}

// Wrap wraps an error as an oops error.
// supports supports [string, any]... pairs or slog.Attr values.
func Wrap(err error, args ...any) error {
	if err == nil {
		return nil
	}
	var sources []Source
	if s := source(); s != nil {
		sources = append(sources, *s)
	}

	var oops Error
	if errors.As(err, &oops) {
		return Error{
			err:        err,
			attributes: argsToAttr(oops.attributes, args),
			sources:    append(sources, oops.sources...),
			code:       oops.code,
		}
	}

	return Error{err: err, attributes: argsToAttr(nil, args), sources: sources}
}

// Errorf formats a string and returns a new oops error.
func Errorf(format string, args ...any) Error {
	var sources []Source
	if s := source(); s != nil {
		sources = append(sources, *s)
	}

	e := Error{err: fmt.Errorf(format, args...), sources: sources}
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			var oops Error
			if errors.As(err, &oops) {
				e.attributes = oops.attributes
				e.sources = append(e.sources, oops.sources...)
			}
		}
	}

	return e
}

// Code returns the oops error code, if the error is an oops error.
// if the error code is not set, it looks for a "code" attribute set to an int.
func Code(err error) int {
	var e Error
	if errors.As(err, &e) {
		if e.code > 0 {
			return e.code
		} else if code, ok := e.attributes["code"].(int); ok {
			return code
		}
	}
	return 0
}

// EnableSource enables or disables source tracking.
// default is true.
func EnableSource(enabled bool) {
	sourceEnabled = enabled
}
