package go2errorsinspection

import (
	"errors"
	"fmt"
	"reflect"
)

type Wrapper interface {
	Unwrap() error
}

func Is(err, target error) bool {
	for {
		if err == target {
			return true
		}

		wrapper, ok := err.(Wrapper)
		if !ok {
			return false
		}

		if err = wrapper.Unwrap(); err == nil {
			return false
		}
	}
}

func AsValue(target interface{}, err error) bool {
	v := reflect.ValueOf(target)
	if v.Type().Kind() != reflect.Ptr {
		panic("non-pointer passed to AsValue")
	}

	v = v.Elem()
	vt := v.Type()

	for {
		et := reflect.TypeOf(err)
		if et == vt || vt.Kind() == reflect.Interface && et.Implements(vt) {
			v.Set(reflect.ValueOf(err))
			return true
		}

		wrapper, ok := err.(Wrapper)
		if !ok {
			return false
		}

		if err = wrapper.Unwrap(); err == nil {
			return false
		}
	}
}

type errWrapper struct {
	err     error
	wrapped error
}

func (ew *errWrapper) Error() string {
	if ew.wrapped == nil {
		return ew.err.Error()
	}

	return fmt.Sprintf("%s: %s", ew.err, ew.wrapped)
}

func (ew *errWrapper) Unwrap() error {
	return ew.wrapped
}

func WithError(wrapped error, err error) error {
	if wrapped == nil {
		return nil
	}

	if err == nil {
		return wrapped
	}

	return &errWrapper{
		err:     err,
		wrapped: wrapped,
	}
}

type ErrText struct {
	Text    string
	Wrapped error
}

func (et *ErrText) Error() string {
	if et.Wrapped == nil {
		return et.Text
	}

	return fmt.Sprintf("%s: %s", et.Text, et.Wrapped)
}

func (et *ErrText) Unwrap() error {
	return et.Wrapped
}

func New(text string) error {
	return errors.New(text)
}

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
