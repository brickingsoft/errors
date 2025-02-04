package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

func Define(message string) error {
	return errors.New(message)
}

type Options struct {
	Meta    Meta
	Wrapped *EnhancedError
	Depth   int
}

type Option func(*Options)

func WithWrapped(err error) Option {
	return func(o *Options) {
		if err != nil {
			if ee, ok := err.(*EnhancedError); ok {
				o.Wrapped = ee
			} else {
				o.Wrapped = &EnhancedError{
					Message: err.Error(),
				}
			}
		}
	}
}

func WithDepth(n int) Option {
	return func(o *Options) {
		o.Depth = n
	}
}

func New(message string, opt ...Option) error {
	opts := Options{
		Meta:    nil,
		Wrapped: nil,
		Depth:   1,
	}
	for _, o := range opt {
		o(&opts)
	}

	st := Stacktrace{}
	pc, file, line, ok := runtime.Caller(opts.Depth)
	if ok {
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

		st.Fn = fn.Name()
		st.File = file
		st.Line = line
	}

	return &EnhancedError{
		Message:    message,
		Meta:       opts.Meta,
		Stacktrace: st,
		Wrapped:    opts.Wrapped,
	}
}

type Stacktrace struct {
	Fn   string
	File string
	Line int
}

type EnhancedError struct {
	Message    string
	Stacktrace Stacktrace
	Meta       Meta
	Wrapped    *EnhancedError
}

func (e *EnhancedError) Error() string {
	return e.Message
}

func (e *EnhancedError) Unwrap() error {
	return e.Wrapped
}

func (e *EnhancedError) Is(err error) bool {
	if err == nil {
		return false
	}
	if ee, ok := err.(*EnhancedError); ok {
		if e.Message == ee.Message {
			return true
		}
	}
	return true
}

func (e *EnhancedError) String() string {
	return fmt.Sprintf("%v", e)
}

func (e *EnhancedError) write(state fmt.State) {
	buf := acquireByteBuffer()
	_, _ = buf.WriteString("EnhancedError:\n")
	var err = e
WRITE:
	_, _ = buf.WriteString(">>>>>>>>>>>>>\n")
	_, _ = buf.WriteString(fmt.Sprintf("ERRO      = %s\n", err.Message))
	if err.Meta.Len() > 0 {
		_, _ = buf.WriteString("META      =")
		for j := range err.Meta {
			if key, val := err.Meta[j].Key, err.Meta[j].Value; key != "" {
				_, _ = buf.WriteString(fmt.Sprintf(" [%s: %v]", key, val))
			}
		}
		_, _ = buf.WriteString("\n")
	}
	if fn, file, line := err.Stacktrace.Fn, err.Stacktrace.File, err.Stacktrace.Line; file != "" {
		_, _ = buf.WriteString(fmt.Sprintf("FUNC      = %s\n", fn))
		_, _ = buf.WriteString(fmt.Sprintf("SEEK      = %s:%d\n", file, line))
	}
	_, _ = buf.WriteString("<<<<<<<<<<<<<\n")
	if err.Wrapped != nil {
		err = err.Wrapped
		goto WRITE
	}
	b := buf.Bytes()
	bLen := len(b)
	content := unsafe.String(&b[0], bLen)
	_, _ = fmt.Fprint(state, content)
	releaseByteBuffer(buf)
	return
}

func (e *EnhancedError) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = fmt.Fprint(state, e.Message)
		break
	case 'v':
		e.write(state)
		break
	default:
		_, _ = fmt.Fprint(state, e.Message)
		break
	}
}
