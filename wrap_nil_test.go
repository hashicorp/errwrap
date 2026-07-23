package errwrap

import (
	"errors"
	"testing"
)

func TestWrapNil(t *testing.T) {
	if WrapNil(nil, nil) != nil {
		t.Fatal("expected nil")
	}
	inner := errors.New("inner")
	if WrapNil(nil, inner) != inner {
		t.Fatal("expected inner")
	}
	outer := errors.New("outer")
	if WrapNil(outer, nil) != outer {
		t.Fatal("expected outer")
	}
	w := WrapNil(outer, inner)
	if w == nil || !Contains(w, "inner") || !Contains(w, "outer") {
		t.Fatalf("wrap failed: %v", w)
	}
	if WrapfNil("x: {{err}}", nil) != nil {
		t.Fatal("WrapfNil nil")
	}
	if WrapfNil("x: {{err}}", inner) == nil {
		t.Fatal("WrapfNil expected error")
	}
}
