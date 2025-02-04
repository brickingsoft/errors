package errors

import (
	"runtime"
	"strings"
)

type Stacktrace struct {
	Fn   string
	File string
	Line int
}

func newStacktrace(skip int) Stacktrace {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return Stacktrace{
			Fn:   "unknown",
			File: "unknown",
			Line: 0,
		}
	}
	if strings.IndexByte(file, '/') == 0 || strings.IndexByte(file, ':') == 1 {
		idx := strings.Index(file, "/src/")
		if idx > 0 {
			file = file[idx+5:]
		} else {
			idx = strings.Index(file, "/pkg/mod/")
			if idx > 0 {
				file = file[idx+9:]
			}
		}
	}
	fn := runtime.FuncForPC(pc)
	return Stacktrace{
		Fn:   fn.Name(),
		File: file,
		Line: line,
	}
}
