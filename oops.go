package oops

import (
	"errors"
	"fmt"
)

// WithSource is a package variable to enable/disable tracking the source location.
var WithSource = true

// New returns a new oops error.
func New(msg string) Error {
	var sources []Source
	if s := source(); s != nil {
		sources = append(sources, *s)
	}

	return Error{err: errors.New(msg), sources: sources}
}

// Wrap wraps an error as an oops error.
func Wrap(err error) Error {
	var sources []Source
	if s := source(); s != nil {
		sources = append(sources, *s)
	}

	var oops Error
	if errors.As(err, &oops) {
		return Error{
			err:        err,
			attributes: oops.attributes,
			sources:    append(sources, oops.sources...),
		}
	}

	return Error{err: err, sources: sources}
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
