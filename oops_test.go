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
	// level=ERROR msg=error err.err=oops err.id=new err.source.file=oops_test.go err.source.function=github.com/jesse0michael/oops.ExampleNew err.source.line=14
}

func ExampleWrap() {
	l := logger()

	err := errors.New("oops")
	err = Wrap(err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.id=wrap err.source.file=oops_test.go err.source.function=github.com/jesse0michael/oops.ExampleWrap err.source.line=28
}

func ExampleWrap_oopsError() {
	l := logger()

	err := New("oops").With("id", "new", "component", "example")
	err = Wrap(err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// oops
	// level=ERROR msg=error err.err=oops err.id=wrap err.component=example err.sources="[[file=oops_test.go function=github.com/jesse0michael/oops.ExampleWrap_oopsError line=42] [file=oops_test.go function=github.com/jesse0michael/oops.ExampleWrap_oopsError line=41]]"
}

func ExampleErrorf() {
	l := logger()

	err := Errorf("failure: %s, %w", "Errorf", errors.New("oops"))

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// failure: Errorf, oops
	// level=ERROR msg=error err.err="failure: Errorf, oops" err.source.file=oops_test.go err.source.function=github.com/jesse0michael/oops.ExampleErrorf err.source.line=55
}

func ExampleErrorf_wrapOopsError() {
	l := logger()

	err := New("oops").With("id", "new", "component", "example")
	err = Errorf("failure: %s, %w", "Errorf", err).With("id", "wrap")

	fmt.Println(err)
	l.Error("error", "err", err)

	// Output:
	// failure: Errorf, oops
	// level=ERROR msg=error err.err="failure: Errorf, oops" err.id=wrap err.component=example err.sources="[[file=oops_test.go function=github.com/jesse0michael/oops.ExampleErrorf_wrapOopsError line=69] [file=oops_test.go function=github.com/jesse0michael/oops.ExampleErrorf_wrapOopsError line=68]]"
}

// logger returns a slog logger to use in these example tests that removes non-deterministic and host-specific values.
func logger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "file" {
				return slog.String("file", filepath.Base(a.Value.String()))
			}
			if a.Key == "sources" {
				if s, ok := a.Value.Any().([]slog.Value); ok {
					sources := make([]slog.Value, len(s))
					for i, v := range s {
						vals := make([]slog.Attr, len(v.Group()))
						for x, y := range v.Group() {
							if y.Key == "file" {
								vals[x] = slog.String("file", filepath.Base(y.Value.String()))
							} else {
								vals[x] = y
							}
						}
						sources[i] = slog.GroupValue(vals...)
					}
					return slog.Any("sources", sources)
				}
			}
			return a
		},
	}))
}
