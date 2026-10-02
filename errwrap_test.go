// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package errwrap

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrappedError_impl(t *testing.T) {
	var _ error = new(wrappedError)
}

func TestGetAll(t *testing.T) {
	cases := []struct {
		Err error
		Msg string
		Len int
	}{
		{},
		{
			fmt.Errorf("foo"),
			"foo",
			1,
		},
		{
			fmt.Errorf("bar"),
			"foo",
			0,
		},
		{
			Wrapf("bar", fmt.Errorf("foo")),
			"foo",
			1,
		},
		{
			Wrapf("{{err}}", fmt.Errorf("foo")),
			"foo",
			2,
		},
		{
			Wrapf("bar", Wrapf("baz", fmt.Errorf("foo"))),
			"foo",
			1,
		},
		{
			fmt.Errorf("foo: %w", fmt.Errorf("bar")),
			"foo: bar",
			1,
		},
		{
			fmt.Errorf("foo: %w", fmt.Errorf("bar")),
			"bar",
			1,
		},
	}

	for i, tc := range cases {
		actual := GetAll(tc.Err, tc.Msg)
		if len(actual) != tc.Len {
			t.Fatalf("%d: bad: %#v", i, actual)
		}
		for _, v := range actual {
			if v.Error() != tc.Msg {
				t.Fatalf("%d: bad: %#v", i, actual)
			}
		}
	}
}

func TestGetAllType(t *testing.T) {
	cases := []struct {
		Err  error
		Type interface{}
		Len  int
	}{
		{},
		{
			fmt.Errorf("foo"),
			"foo",
			0,
		},
		{
			fmt.Errorf("bar"),
			fmt.Errorf("foo"),
			1,
		},
		{
			Wrapf("bar", fmt.Errorf("foo")),
			fmt.Errorf("baz"),
			2,
		},
		{
			Wrapf("bar", Wrapf("baz", fmt.Errorf("foo"))),
			Wrapf("", nil),
			0,
		},
		{
			fmt.Errorf("one: %w", fmt.Errorf("two: %w", fmt.Errorf("three"))),
			fmt.Errorf("%w", errors.New("")),
			2,
		},
	}

	for i, tc := range cases {
		actual := GetAllType(tc.Err, tc.Type)
		if len(actual) != tc.Len {
			t.Fatalf("%d: bad: %#v", i, actual)
		}
	}
}

func TestWrappedError_IsCompatibleWithErrorsUnwrap(t *testing.T) {
	inner := errors.New("inner error")
	err := Wrap(errors.New("outer"), inner)
	actual := errors.Unwrap(err)
	if actual != inner {
		t.Fatal("wrappedError did not unwrap to inner")
	}
}

func TestWrapNil(t *testing.T) {
	if err := WrapNil(errors.New("outer"), nil); err != nil {
		t.Fatalf("expected nil when inner is nil, got: %v", err)
	}

	if err := WrapNil(nil, nil); err != nil {
		t.Fatalf("expected nil when both are nil, got: %v", err)
	}

	inner := errors.New("inner")
	outer := errors.New("outer")
	err := WrapNil(outer, inner)
	if err == nil {
		t.Fatal("expected non-nil error when inner is non-nil")
	}
	if err.Error() != "outer" {
		t.Fatalf("expected error %q, got %q", "outer", err.Error())
	}
	if !Contains(err, "inner") {
		t.Fatal("expected error to contain inner")
	}
	if !Contains(err, "outer") {
		t.Fatal("expected error to contain outer")
	}
	if actual := errors.Unwrap(err); actual != inner {
		t.Fatalf("expected unwrapped error to be inner, got %v", actual)
	}
}

func TestWrapfNil(t *testing.T) {
	if err := WrapfNil("failed: {{err}}", nil); err != nil {
		t.Fatalf("expected nil when err is nil, got: %v", err)
	}

	inner := errors.New("inner")
	err := WrapfNil("failed: {{err}}", inner)
	if err == nil {
		t.Fatal("expected non-nil error when err is non-nil")
	}
	if err.Error() != "failed: inner" {
		t.Fatalf("expected error %q, got %q", "failed: inner", err.Error())
	}
	if !Contains(err, "inner") {
		t.Fatal("expected error to contain inner")
	}
	if actual := errors.Unwrap(err); actual != inner {
		t.Fatalf("expected unwrapped error to be inner, got %v", actual)
	}
}
