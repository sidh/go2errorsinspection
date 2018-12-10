package go2errorsinspection_test

import (
	"testing"

	errors "github.com/sidh/go2errorsinspection"
	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	sentinel := errors.New("sentinel")

	// Sentinel error wrapped by another error, Is can find sentinel error
	{
		err := errors.WithError(sentinel, errors.New("wrapper"))
		assert.True(t, errors.Is(err, sentinel))
	}

	// Sentinel error wraps another error, Is cannot find sentinel error
	// This check is impossible to implement
	{
		err := errors.WithError(errors.New("wrapped"), sentinel)
		assert.True(t, errors.Is(err, sentinel))
	}
}

func TestAs(t *testing.T) {
	type myError struct {
		errors.ErrText
	}

	// Typed error wrapped by another error, As can find typed error
	{
		var me *myError
		err := errors.WithError(
			&myError{ErrText: errors.ErrText{
				Text: "my error message",
			}},
			errors.New("wrapper"),
		)
		assert.True(t, errors.AsValue(&me, err))
	}

	// Typed error wraps another error, As can find typed error
	{
		var me *myError
		err := &myError{ErrText: errors.ErrText{
			Text:    "my error message",
			Wrapped: errors.New("wrapped"),
		}}
		assert.True(t, errors.AsValue(&me, err))
	}
}
