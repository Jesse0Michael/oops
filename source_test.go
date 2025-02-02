package oops

import (
	"strings"
	"testing"
)

// wrap is a helper function to test the correct call depth.
func wrap() *Source {
	return source()
}

func Test_source(t *testing.T) {
	defer func() { WithSource = true }()
	tests := []struct {
		name    string
		enabled bool
		want    *Source
	}{
		{
			name:    "enabled",
			enabled: true,
			want: &Source{
				Function: "github.com/jesse0michael/oops.Test_source.func2",
				File:     "source_test.go",
				Line:     38,
			},
		},
		{
			name:    "disabled",
			enabled: false,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithSource = tt.enabled
			got := wrap()

			if tt.want == nil {
				if got != nil {
					t.Errorf("source() = %v, want %v", got, tt.want)
				}
				return
			}

			if got == nil {
				t.Errorf("source() = %v, want %v", got, tt.want)
				return
			}

			if got.Function != tt.want.Function {
				t.Errorf("source().Function = %v, want %v", got.Function, tt.want.Function)
			}
			if !strings.HasSuffix(got.File, tt.want.File) {
				t.Errorf("source().File = %v, want suffix %v", got.File, tt.want.File)
			}
			if got.Line != tt.want.Line {
				t.Errorf("source().Line = %v, want %v", got.Line, tt.want.Line)
			}
		})
	}
}
