package oops

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestOopsHandler_Handle(t *testing.T) {
	tests := []struct {
		name  string
		attrs []slog.Attr
		log   string
	}{
		{
			name:  "no attributes",
			attrs: []slog.Attr{},
			log:   `{"level":"ERROR","msg":"message"}`,
		},
		{
			name: "basic attribute types",
			attrs: []slog.Attr{
				slog.String("string", "value"),
				slog.Int("int", 42),
				slog.Int64("int64", 1234567890),
				slog.Bool("bool", true),
				slog.Float64("float", 3.14),
				slog.Any("error", errors.New("oops")),
			},
			log: `{"level":"ERROR","msg":"message","string":"value","int":42,"int64":1234567890,"bool":true,"float":3.14,"error":"oops"}`,
		},
		{
			name: "empty oops",
			attrs: []slog.Attr{
				slog.Any("error", &Error{}),
			},
			log: `{"level":"ERROR","msg":"message","error":null}`,
		},
		{
			name: "oops error without attributes",
			attrs: []slog.Attr{
				slog.Any("error", New("oops")),
			},
			log: `{"level":"ERROR","msg":"message","error":"oops","oops":{"source":[{"function":"github.com/jesse0michael/oops.TestOopsHandler_Handle","file":"handler_test.go","line":48}]}}`,
		},
		{
			name: "oops error with attributes",
			attrs: []slog.Attr{
				slog.Any("error", New("oops").With("id", "new", "component", "example")),
			},
			log: `{"level":"ERROR","msg":"message","error":"oops","oops":{"component":"example","id":"new","source":[{"function":"github.com/jesse0michael/oops.TestOopsHandler_Handle","file":"handler_test.go","line":55}]}}`,
		},
		{
			name: "oops wrapped in non-oops error",
			attrs: []slog.Attr{
				slog.Any("error", &wrappedError{underlying: New("oops").With("id", "new", "component", "example")}),
			},
			log: `{"level":"ERROR","msg":"message","error":"oops","oops":{"component":"example","id":"new","source":[{"function":"github.com/jesse0michael/oops.TestOopsHandler_Handle","file":"handler_test.go","line":62}]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for output
			f, _ := os.CreateTemp(t.TempDir(), "out")
			os.Stderr = f

			h := NewOopsHandler(
				slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
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
				}),
			)

			// Create and process the record
			record := slog.NewRecord(time.Time{}, slog.LevelError, "message", 0)
			record.AddAttrs(tt.attrs...)
			_ = h.Handle(t.Context(), record)

			// Read the output and compare
			_, _ = f.Seek(0, 0)
			b, _ := io.ReadAll(f)
			output := strings.TrimSpace(string(b))

			// Compare the output to expected log
			if !reflect.DeepEqual(tt.log, output) {
				t.Errorf("OopsHandler.Handle =\n%v\nwant\n%v", output, tt.log)
			}

			f.Close()
		})
	}
}
