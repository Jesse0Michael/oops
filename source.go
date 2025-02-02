package oops

import "runtime"

type Source struct {
	Function string
	File     string
	Line     int
}

func source() *Source {
	if WithSource {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc := pcs[0]
		fs := runtime.CallersFrames([]uintptr{pc})
		f, _ := fs.Next()
		return &Source{
			Function: f.Function,
			File:     f.File,
			Line:     f.Line,
		}
	}
	return nil
}
