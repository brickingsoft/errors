package errors

import "errors"

// Join
// 组合错误
func Join(errs ...error) error {
	errsLen := len(errs)
	if errsLen == 0 {
		return nil
	}
	var prev *EnhancedError
	for i := errsLen - 1; i > -1; i-- {
		err := errs[i]
		if err == nil {
			continue
		}
		ee, ok := err.(*EnhancedError)
		if !ok {
			err = New(err.Error(), WithDepth(2))
			ee = err.(*EnhancedError)
		}
		if prev == nil {
			prev = ee
			continue
		}
		ee.Wrapped = prev
		prev = ee
	}
	return prev
}

// Is
// 判断 err 是否为或包含 target。
// err 可以是一组错误，target 建议是一个错误，即便 target 是一组，也是取顶层进行判断。
func Is(err error, target error) bool {
	if err == nil || target == nil {
		return errors.Is(err, target)
	}
	if eq := err.Error() == target.Error(); eq {
		return true
	}
	return errors.Is(err, target)
}

// As
// 转化为目标。
// target 必须是指针。当类型是指针，target 为指针的指针。
func As(err error, target any) bool {
	return errors.As(err, target)
}

// AsEnhancedError
// 转化为 EnhancedError
func AsEnhancedError(err error) (*EnhancedError, bool) {
	if err == nil {
		return nil, false
	}

	ee, ok := err.(*EnhancedError)
	return ee, ok
}

// Unwrap
// 取包裹
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// StackOf
// 取跟踪栈
func StackOf(err error) (fn string, file string, line int, ok bool) {
	ee, isEE := AsEnhancedError(err)
	if isEE {
		fn, file, line = ee.Stack()
		ok = true
	}
	return
}
