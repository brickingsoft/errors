package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brickingsoft/errors"
	"testing"
	"time"
)

func TestErr(t *testing.T) {
	e := errors.New("error")
	t.Log(e)

	def := errors.Define("def")
	t.Log(def)
	ee := errors.From(def, errors.WithWrap(errors.Define("wrapped")),
		errors.WithMeta("s", "s"),
		errors.WithMeta("i", 1),
		errors.WithMeta("i32", int32(32)),
		errors.WithMeta("i64", int64(-64)),
		errors.WithMeta("u", uint(1)),
		errors.WithMeta("u64", uint64(64)),
		errors.WithMeta("f32", float32(32.32)),
		errors.WithMeta("f64", 64.640),
		errors.WithMeta("b", true),
		errors.WithMeta("any", struct{}{}),
		errors.WithMeta("byte", 'b'),
		errors.WithMeta("bytes", []byte("hello world")),
		errors.WithMeta("time", time.Now()),
		errors.WithMeta("ss", []string{"a a", "b"}),
		errors.WithDescription("desc"),
		errors.WithOccur(),
	)
	t.Log(ee)

	t.Log(fmt.Sprintf("%s", ee))
}

func TestJson(t *testing.T) {

	e1 := errors.Define("e1")
	e2 := errors.New("e2", errors.WithWrap(e1))

	e3 := errors.New("e3", errors.WithWrap(e2),
		errors.WithMeta("s", "s"),
		errors.WithMeta("i", 1),
		errors.WithMeta("i32", int32(32)),
		errors.WithMeta("i64", int64(-64)),
		errors.WithMeta("u", uint(1)),
		errors.WithMeta("u64", uint64(64)),
		errors.WithMeta("f32", float32(32.32)),
		errors.WithMeta("f64", 64.640),
		errors.WithMeta("b", true),
		errors.WithMeta("any", struct{}{}),
		errors.WithMeta("byte", 'b'),
		errors.WithMeta("bytes", []byte("hello world")),
		errors.WithMeta("time", time.Now()),
		errors.WithMeta("ss", []string{"a a", "b"}),
		errors.WithDescription("desc"),
		errors.WithOccur(),
	)

	b, err := json.Marshal(e3)
	if err != nil {
		t.Fatal(err)
		return
	}
	buf := bytes.NewBuffer(nil)
	_ = json.Indent(buf, b, "", "\t")
	t.Log(buf.String())

}

func TestWithoutStacktrace(t *testing.T) {
	e := errors.New("e", errors.WithoutStacktrace())
	t.Log(e)
}

func TestFrom(t *testing.T) {
	define := errors.Define("err1", errors.WithMeta("meta", "value"), errors.WithWrap(errors.Define("err2")))
	err := errors.From(define, errors.WithWrap(errors.New("err3")))
	t.Log(err)
}
