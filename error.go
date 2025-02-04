package errors

import (
	"errors"
	"fmt"
	"unsafe"
)

func Define(message string) error {
	return errors.New(message)
}

type Options struct {
	Meta    Meta
	Wrapped error
	Depth   int
}

type Option func(*Options)

func WithWrapped(err error) Option {
	return func(o *Options) {
		if err != nil {
			o.Wrapped = err
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
		Depth:   2,
	}
	for _, o := range opt {
		o(&opts)
	}

	return &EnhancedError{
		Message:    message,
		Meta:       opts.Meta,
		Stacktrace: newStacktrace(opts.Depth),
		Wrapped:    opts.Wrapped,
	}
}

type EnhancedError struct {
	Message    string
	Stacktrace Stacktrace
	Meta       Meta
	Wrapped    error
}

func (e *EnhancedError) Error() string {
	return e.String()
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
	return fmt.Sprintf("%+s", e)
}

func (e *EnhancedError) write(state fmt.State) {
	buf := acquireByteBuffer()
	_, _ = buf.WriteString("EnhancedError:\n")
	var err error = e
WRITE:
	switch ee := err.(type) {
	case *EnhancedError:
		_, _ = buf.WriteString(">>>>>>>>>>>>>\n")
		_, _ = buf.WriteString(fmt.Sprintf("ERRO      = %s\n", ee.Message))
		if ee.Meta.Len() > 0 {
			_, _ = buf.WriteString("META      =")
			for j := range ee.Meta {
				if key, val := ee.Meta[j].Key, ee.Meta[j].Value; key != "" {
					_, _ = buf.WriteString(fmt.Sprintf(" [%s: %v]", key, val))
				}
			}
			_, _ = buf.WriteString("\n")
		}
		if fn, file, line := ee.Stacktrace.Fn, ee.Stacktrace.File, ee.Stacktrace.Line; file != "" {
			_, _ = buf.WriteString(fmt.Sprintf("FUNC      = %s\n", fn))
			_, _ = buf.WriteString(fmt.Sprintf("SEEK      = %s:%d\n", file, line))
		}
		_, _ = buf.WriteString("<<<<<<<<<<<<<\n")
		if ee.Wrapped != nil {
			err = ee.Wrapped
			goto WRITE
		}
		break
	default:
		_, _ = buf.WriteString(">>>>>>>>>>>>>\n")
		_, _ = buf.WriteString(err.Error())
		_, _ = buf.WriteString("\n<<<<<<<<<<<<<\n")
		break
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
		switch {
		case state.Flag('+'):
			e.write(state)
			break
		default:
			_, _ = fmt.Fprint(state, e.Message)
			break
		}
	case 'v':
		e.write(state)
		break
	default:
		_, _ = fmt.Fprint(state, e.Message)
		break
	}
}
