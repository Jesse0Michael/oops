package oops

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

// Custom error type for testing errors.As
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

// Wrapped error type for testing error unwrapping
type wrappedError struct {
	underlying error
}

func (w *wrappedError) Error() string {
	return w.underlying.Error()
}

func (w *wrappedError) Unwrap() error {
	return w.underlying
}

func TestError_Unwrap(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		isOopsError   bool
		asOopsError   bool
		isCustomError bool
		asCustomError bool
	}{
		{
			name:          "standard error",
			err:           errors.New("test-error"),
			asOopsError:   false,
			asCustomError: false,
		},
		{
			name:          "custom error",
			err:           &customError{msg: "wrapped-error"},
			asOopsError:   false,
			asCustomError: true,
		},
		{
			name:          "oops wrapped standard error",
			err:           Wrap(errors.New("test-error")),
			asOopsError:   true,
			asCustomError: false,
		},
		{
			name:          "oops wrapped custom error",
			err:           Wrap(&customError{msg: "wrapped-error"}),
			asOopsError:   true,
			asCustomError: true,
		},
		{
			name:          "Nested wrapped errors",
			err:           Wrap(Wrap(errors.New("test-error"))),
			asOopsError:   true,
			asCustomError: false,
		},
		{
			name:          "Not matching error",
			err:           Wrap(Wrap(&customError{msg: "wrapped-error"})),
			asOopsError:   true,
			asCustomError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// as oops error
			var oops *Error
			if got := errors.As(tt.err, &oops); got != tt.asOopsError {
				t.Errorf("errors.As(*Error) = %v, want %v", got, tt.asOopsError)
			}

			// as custom error
			var custom *customError
			if got := errors.As(tt.err, &custom); got != tt.asCustomError {
				t.Errorf("errors.As(*customError) = %v, want %v", got, tt.asCustomError)
			}
		})
	}
}

func Test_argsToAttr(t *testing.T) {
	tests := []struct {
		name  string
		attrs map[string]any
		args  []any
		want  map[string]any
	}{
		{
			name:  "nil args",
			attrs: nil,
			args:  nil,
			want:  nil,
		},
		{
			name:  "one arg",
			attrs: nil,
			args:  []any{"key"},
			want:  map[string]any{},
		},
		{
			name:  "one slog attr",
			attrs: nil,
			args:  []any{slog.String("key", "value")},
			want:  map[string]any{"key": "value"},
		},
		{
			name:  "arg pair",
			attrs: nil,
			args:  []any{"key", "value"},
			want:  map[string]any{"key": "value"},
		},
		{
			name:  "arg pairs",
			attrs: nil,
			args:  []any{"key", "value", "foo", true, "bar", 42},
			want:  map[string]any{"key": "value", "foo": true, "bar": 42},
		},
		{
			name:  "non string key",
			attrs: nil,
			args:  []any{1, "value"},
			want:  map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := argsToAttr(tt.attrs, tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("argsToAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
