package oops

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func ExampleNew() {
	l := logger()

	err := New("oops").With("id", "new")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.id=new err.source="[{Function:github.com/jesse0michael/oops.ExampleNew File:oops_test.go Line:15}]"
}

func ExampleWrap() {
	l := logger()

	err := errors.New("oops")
	err = Wrap(err, "id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.id=wrap err.source="[{Function:github.com/jesse0michael/oops.ExampleWrap File:oops_test.go Line:29}]"
}

func ExampleWrap_oopsError() {
	l := logger()

	oopsErr := New("oops").With("id", "new", "component", "example")
	err := Wrap(oopsErr, "id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.component=example err.id=wrap err.source="[{Function:github.com/jesse0michael/oops.ExampleWrap_oopsError File:oops_test.go Line:42} {Function:github.com/jesse0michael/oops.ExampleWrap_oopsError File:oops_test.go Line:43}]"
}

func ExampleErrorf() {
	l := logger()

	err := Errorf("failure: %s, %w", "Errorf", errors.New("oops"))

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// failure: Errorf, oops
	// level=ERROR msg=error err.err="failure: Errorf, oops" err.source="[{Function:github.com/jesse0michael/oops.ExampleErrorf File:oops_test.go Line:56}]"
}

func ExampleErrorf_wrapOopsError() {
	l := logger()

	err := New("oops").With("id", "new", "component", "example")
	err = Errorf("failure: %s, %w", "Errorf", err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// failure: Errorf, oops
	// level=ERROR msg=error err.err="failure: Errorf, oops" err.component=example err.id=wrap err.source="[{Function:github.com/jesse0michael/oops.ExampleErrorf_wrapOopsError File:oops_test.go Line:70} {Function:github.com/jesse0michael/oops.ExampleErrorf_wrapOopsError File:oops_test.go Line:69}]"
}

func ExampleWrap_nil() {
	l := logger()

	err := Wrap(nil)

	if err != nil {
		fmt.Println(err)
	}
	l.Error("error", "err", err)

	// Output:
	// level=ERROR msg=error err=<nil>
}

func ExampleError_Code() {
	l := logger()

	err := New("oops").Code(404)

	fmt.Println(err)
	l.Error("error", "err", err)
	fmt.Println(Code(err))

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.source="[{Function:github.com/jesse0michael/oops.ExampleError_Code File:oops_test.go Line:97}]"
	// 404
}

func ExampleNew_nonOopsWrappedError() {
	// This example shows how the LogValue won't be hit for a wrapped error :(
	l := logger()

	var err error
	err = New("oops").With("id", "new", "component", "example")
	err = &wrappedError{underlying: err}

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err=oops
}

// logger returns a slog logger to use in these example tests that removes non-deterministic and host-specific values.
func logger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "source" {
				if s, ok := a.Value.Any().([]Source); ok {
					for i, v := range s {
						s[i].File = filepath.Base(v.File)
					}
					return slog.Any("source", s)
				}
			}
			return a
		},
	}))
}

func TestCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "nil error",
			err:  nil,
			want: 0,
		},
		{
			name: "regular error",
			err:  errors.New("test-error"),
			want: 0,
		},
		{
			name: "oops error without code",
			err:  New("oops"),
			want: 0,
		},
		{
			name: "oops error with code",
			err:  New("oops").Code(404),
			want: 404,
		},
		{
			name: "Wrapped error",
			err:  errors.Join(errors.New("test-error"), Wrap(New("oops").Code(404))),
			want: 404,
		},
		{
			name: "oops error code attribute",
			err:  New("oops").With("code", 404),
			want: 404,
		},
		{
			name: "oops error non int code attribute",
			err:  New("oops").With("code", "NOT_FOUND"),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Code(tt.err); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}
