package oops

import (
	"log/slog"
	"reflect"
	"testing"
)

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
