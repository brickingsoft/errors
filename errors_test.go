package errors_test

import (
	"fmt"
	"github.com/brickingsoft/errors"
	"testing"
)

func TestJoin(t *testing.T) {
	e1 := errors.New("err1")
	e2 := errors.New("err2")
	e3 := errors.New("err3")
	e := errors.Join(nil, e1, nil, e2, e3)
	t.Log(e)
}

func TestJoin_Std(t *testing.T) {
	e1 := fmt.Errorf("err1")
	e2 := errors.New("err2")

	e := errors.Join(e1, e2)
	t.Log(e)
}

func TestIs(t *testing.T) {
	t.Log(errors.Is(
		errors.New("err"),
		errors.New("err"),
	))
	t.Log(errors.Is(
		errors.New("err"),
		errors.New("err1"),
	))
	t.Log(errors.Is(
		errors.New("err", errors.WithWrap(errors.New("err1"))),
		errors.New("err1"),
	))
	t.Log(errors.Is(
		errors.New("err", errors.WithWrap(errors.New("err1"))),
		fmt.Errorf("err1"),
	))

	t.Log(errors.Is(
		fmt.Errorf("err"),
		errors.New("err"),
	))

	t.Log(errors.Is(
		fmt.Errorf("err"),
		errors.New("err1"),
	))
}

func TestAs(t *testing.T) {
	err := errors.New("err")
	var ee *errors.EnhancedError
	ok := errors.As(err, &ee)
	if ok {
		t.Log(ee)
	} else {
		t.Error("err should be enhanced")
	}
}
