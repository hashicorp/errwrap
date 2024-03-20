//go:build go1.13
// +build go1.13

package errwrap_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/errwrap"
)

func TestWrappedError_IsCompatibleWithErrorsUnwrap(t *testing.T) {
	inner := errors.New("inner error")
	err := errwrap.Wrap(errors.New("outer"), inner)
	actual := errors.Unwrap(err)
	if actual != inner {
		t.Fatal("wrappedError did not unwrap to inner")
	}
}
