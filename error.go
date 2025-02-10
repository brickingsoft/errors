package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

// Define
// 定义一个错误，不包含栈和时间。
func Define(message string, opt ...Option) error {
	opt = append(opt, WithoutStacktrace(), WithoutOccur())
	return New(message, opt...)
}

type Options struct {
	Description       string
	Occur             time.Time
	Meta              Meta
	Wrap              *EnhancedError
	Depth             int
	DisableStacktrace bool
}

type Option func(*Options)

// WithDescription
// 设置描述
func WithDescription(desc string) Option {
	return func(o *Options) {
		o.Description = desc
	}
}

// WithOccur
// 设置发生时间为当前
func WithOccur() Option {
	return func(o *Options) {
		o.Occur = time.Now()
	}
}

// WithOccurAt
// 设置发生时间
func WithOccurAt(t time.Time) Option {
	return func(o *Options) {
		o.Occur = t
	}
}

// WithoutOccur
// 移除发生时间
func WithoutOccur() Option {
	return func(o *Options) {
		o.Occur = time.Time{}
	}
}

// WithWrap
// 添加包裹
func WithWrap(err error) Option {
	return func(o *Options) {
		if err != nil {
			if ee, ok := err.(*EnhancedError); ok {
				o.Wrap = ee
			} else {
				o.Wrap = &EnhancedError{
					Message: err.Error(),
				}
			}
		}
	}
}

// WithDepth
// 设置跟踪深度
func WithDepth(n int) Option {
	return func(o *Options) {
		o.Depth = n
	}
}

// WithoutStacktrace
// 取消跟踪
func WithoutStacktrace() Option {
	return func(o *Options) {
		o.DisableStacktrace = true
	}
}

// From
// 从一个错误中创建一个增强错误。
func From(err error, opt ...Option) error {
	if err == nil {
		return nil
	}
	// options
	opts := Options{
		Description:       "",
		Occur:             time.Time{},
		Meta:              nil,
		Wrap:              nil,
		Depth:             2,
		DisableStacktrace: false,
	}
	for _, o := range opt {
		o(&opts)
	}
	ee, ok := err.(*EnhancedError)
	if ok {
		copied := &EnhancedError{}
		Copy(copied, ee)
		ee = copied
	} else {
		ee = &EnhancedError{
			Message: err.Error(),
		}
	}
	if opts.Description != "" {
		ee.Description = opts.Description
	}
	if len(opts.Meta) > 0 {
		ee.Meta = append(ee.Meta, opts.Meta...)
	}
	if !opts.DisableStacktrace {
		// stack
		st := stack(opts.Depth)
		ee.Stacktrace = st
	}
	if !opts.Occur.IsZero() {
		ee.Occur = opts.Occur
	}
	if opts.Wrap != nil {
		ee.Wrap(opts.Wrap)
	}
	return ee
}

// New
// 创建一个增强错误。
func New(message string, opt ...Option) error {
	// options
	opts := Options{
		Description:       "",
		Occur:             time.Time{},
		Meta:              nil,
		Wrap:              nil,
		Depth:             2,
		DisableStacktrace: false,
	}
	for _, o := range opt {
		o(&opts)
	}
	// stack
	ee := &EnhancedError{
		Message:     message,
		Description: opts.Description,
		Meta:        opts.Meta,
		Occur:       opts.Occur,
		Wrapped:     opts.Wrap,
	}
	if opts.DisableStacktrace {
		return ee
	}
	ee.Stacktrace = stack(opts.Depth)
	return ee
}

func stack(depth int) Stacktrace {
	pc, file, line, ok := runtime.Caller(depth)
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
		fn := runtime.FuncForPC(pc).Name()
		return Stacktrace{
			Fn:   fn,
			File: file,
			Line: line,
		}
	}
	return Stacktrace{}
}

type Stacktrace struct {
	Fn   string
	File string
	Line int
}

type EnhancedError struct {
	Message     string
	Description string
	Stacktrace  Stacktrace
	Occur       time.Time
	Meta        Meta
	Wrapped     *EnhancedError
}

func (e *EnhancedError) Stack() (string, string, int) {
	return e.Stacktrace.Fn, e.Stacktrace.File, e.Stacktrace.Line
}

func (e *EnhancedError) Error() string {
	return e.Message
}

func (e *EnhancedError) Unwrap() error {
	if e.Wrapped != nil {
		return e.Wrapped
	}
	return nil
}

func (e *EnhancedError) Wrap(err error) {
	if e.Wrapped == nil {
		if ee, ok := err.(*EnhancedError); ok {
			e.Wrapped = ee
		} else {
			e.Wrapped = &EnhancedError{
				Message: err.Error(),
			}
		}
		return
	}
	e.Wrapped.Wrap(err)
}

func (e *EnhancedError) Is(err error) bool {
	if err == nil {
		return e == err
	}
	return e.Message == err.Error()
}

func (e *EnhancedError) String() string {
	return fmt.Sprintf("%v", e)
}

func (e *EnhancedError) write(state fmt.State) {
	buf := acquireByteBuffer()
	_, _ = buf.WriteString("[ENHANCED ERROR]\n")
	var err = e
WRITE:
	_, _ = buf.WriteString(">>>>>>>>>>>>>\n")
	_, _ = buf.WriteString(fmt.Sprintf("ERRO      = %s\n", err.Message))
	if err.Description != "" {
		_, _ = buf.WriteString(fmt.Sprintf("DESC      = %s\n", err.Description))
	}
	if err.Meta.Len() > 0 {
		_, _ = buf.WriteString("META      =")
		for j := range err.Meta {
			if key, val := err.Meta[j].Key, err.Meta[j].Value; key != "" {
				_, _ = buf.WriteString(fmt.Sprintf(" [%s: %v]", key, val))
			}
		}
		_, _ = buf.WriteString("\n")
	}
	if !err.Occur.IsZero() {
		_, _ = buf.WriteString(fmt.Sprintf("OCCU      = %s\n", err.Occur.Format(time.RFC3339)))
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

// Copy
// 复制
func Copy(dst *EnhancedError, src *EnhancedError) bool {
	if dst == nil || src == nil {
		return false
	}
	dst.Message = src.Message
	dst.Description = src.Description
	dst.Stacktrace = src.Stacktrace
	dst.Occur = src.Occur
	if metaLen := len(dst.Meta); metaLen > 0 {
		dst.Meta = make(Meta, metaLen)
		copy(dst.Meta, src.Meta)
	}
	if wrapped := src.Wrapped; wrapped != nil {
		dst.Wrapped = &EnhancedError{}
		Copy(dst.Wrapped, wrapped)
	}
	return true
}
