package oops

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func ExampleNew() {
	l := logger()

	err := New("oops").With("id", "new")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.id=new err.source="[{Function:github.com/jesse0michael/oops.ExampleNew File:oops_test.go Line:14}]"
}

func ExampleWrap() {
	l := logger()

	err := errors.New("oops")
	err = Wrap(err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.id=wrap err.source="[{Function:github.com/jesse0michael/oops.ExampleWrap File:oops_test.go Line:28}]"
}

func ExampleWrap_oopsError() {
	l := logger()

	err := New("oops").With("id", "new", "component", "example")
	err = Wrap(err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.component=example err.id=wrap err.source="[{Function:github.com/jesse0michael/oops.ExampleWrap_oopsError File:oops_test.go Line:42} {Function:github.com/jesse0michael/oops.ExampleWrap_oopsError File:oops_test.go Line:41}]"
}

func ExampleErrorf() {
	l := logger()

	err := Errorf("failure: %s, %w", "Errorf", errors.New("oops"))

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// failure: Errorf, oops
	// level=ERROR msg=error err.err="failure: Errorf, oops" err.source="[{Function:github.com/jesse0michael/oops.ExampleErrorf File:oops_test.go Line:55}]"
}

func ExampleErrorf_wrapOopsError() {
	l := logger()

	err := New("oops").With("id", "new", "component", "example")
	err = Errorf("failure: %s, %w", "Errorf", err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// failure: Errorf, oops
	// level=ERROR msg=error err.err="failure: Errorf, oops" err.component=example err.id=wrap err.source="[{Function:github.com/jesse0michael/oops.ExampleErrorf_wrapOopsError File:oops_test.go Line:69} {Function:github.com/jesse0michael/oops.ExampleErrorf_wrapOopsError File:oops_test.go Line:68}]"
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
