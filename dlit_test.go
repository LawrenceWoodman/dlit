package dlit

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		in        interface{}
		wantError error
	}{
		{6, nil},
		{6.0, nil},
		{6.6, nil},
		{float32(6.6), nil},
		{int64(922336854775807), nil},
		{int64(9223372036854775807), nil},
		{int64(9223372036854775807), nil},
		{"98292223372036854775807", nil},
		{complex64(1), ErrInvalidKind("complex64")},
		{complex128(1), ErrInvalidKind("complex128")},
		{"6", nil},
		{"6.6", nil},
		{"abc", nil},
		{true, nil},
		{false, nil},
		{errors.New("This is an error"), nil},
	}

	for _, c := range cases {
		_, err := New(c.in)
		if !errorMatch(err, c.wantError) {
			t.Errorf("New(%q) - err == %q, wantError == %q",
				c.in, err, c.wantError)
		}
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      int64
		wantError error
	}{
		{makeLit(6), 6, nil},
		{makeLit(6.0), 6, nil},
		{makeLit("6"), 6, nil},
		{makeLit("98292223372036854775807"), 0,
			ErrInvalidCast{makeLit("98292223372036854775807"), "int"}},
		{makeLit(6.6), 0, ErrInvalidCast{makeLit(6.6), "int"}},
		{makeLit("6.6"), 0, ErrInvalidCast{makeLit("6.6"), "int"}},
		{makeLit("abc"), 0, ErrInvalidCast{makeLit("abc"), "int"}},
		{makeLit(true), 0, ErrInvalidCast{makeLit(true), "int"}},
		{makeLit(false), 0, ErrInvalidCast{makeLit(false), "int"}},
		{makeLit(errors.New("This is an error")), 0,
			ErrInvalidCast{makeLit(errors.New("This is an error")), "int"}},
	}

	for _, c := range cases {
		got, err := c.in.Int()
		if !errorMatch(err, c.wantError) {
			t.Errorf("Int() with Literal: %q - err == %q, wantError == %q",
				c.in, err, c.wantError)
		}
		if got != c.want {
			t.Errorf("Int() with Literal: %q - return: %q, want: %q",
				c.in, got, c.want)
		}
	}
}

func TestFloat(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      float64
		wantError error
	}{
		{makeLit(6), 6.0, nil},
		{makeLit(int64(922336854775807)), 922336854775807.0, nil},
		{makeLit(int64(9223372036854775807)), 0.0,
			ErrInvalidCast{makeLit(int64(9223372036854775807)), "float"}},
		{makeLit(6.0), 6.0, nil},
		{makeLit("6"), 6.0, nil},
		{makeLit(6.678934), 6.678934, nil},
		{makeLit("6.678394"), 6.678394, nil},
		{makeLit("abc"), 0, ErrInvalidCast{makeLit("abc"), "float"}},
		{makeLit(true), 0, ErrInvalidCast{makeLit(true), "float"}},
		{makeLit(false), 0, ErrInvalidCast{makeLit(false), "float"}},
		{makeLit(errors.New("This is an error")), 0,
			ErrInvalidCast{makeLit(errors.New("This is an error")), "float"}},
	}

	for _, c := range cases {
		got, err := c.in.Float()
		if !errorMatch(err, c.wantError) {
			t.Errorf("Float() with Literal: %q - err == %q, wantError == %q",
				c.in, err, c.wantError)
		}
		if got != c.want {
			t.Errorf("Float() with Literal: %q - return: %q, want: %q",
				c.in, got, c.want)
		}
	}
}

func TestBool(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      bool
		wantError error
	}{
		{makeLit(1), true, nil},
		{makeLit(2), true, nil},
		{makeLit(0), false, nil},
		{makeLit(1.0), true, nil},
		{makeLit(2.0), true, nil},
		{makeLit(2.25), true, nil},
		{makeLit(0.0), false, nil},
		{makeLit(true), true, nil},
		{makeLit(false), false, nil},
		{makeLit("true"), true, nil},
		{makeLit("false"), false, nil},
		{makeLit("True"), true, nil},
		{makeLit("False"), false, nil},
		{makeLit("TRUE"), true, nil},
		{makeLit("FALSE"), false, nil},
		{makeLit("t"), true, nil},
		{makeLit("f"), false, nil},
		{makeLit("T"), true, nil},
		{makeLit("F"), false, nil},
		{makeLit("1"), true, nil},
		{makeLit("0"), false, nil},
		{makeLit("bob"), false, ErrInvalidCast{makeLit("bob"), "bool"}},
		{makeLit(errors.New("This is an error")), false,
			ErrInvalidCast{makeLit(errors.New("This is an error")), "bool"}},
	}

	for _, c := range cases {
		got, err := c.in.Bool()
		if !errorMatch(err, c.wantError) {
			t.Errorf("Bool() with Literal: %q - err == %q, wantError == %q",
				c.in, err, c.wantError)
		}
		if got != c.want {
			t.Errorf("Bool() with Literal: %q - return: %q, want: %q",
				c.in, got, c.want)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in   *Literal
		want string
	}{
		{makeLit(124), "124"},
		{makeLit(int64(922336854775807)), "922336854775807"},
		{makeLit(int64(9223372036854775807)), "9223372036854775807"},
		{makeLit("98292223372036854775807"), "98292223372036854775807"},
		{makeLit("Hello my name is fred"), "Hello my name is fred"},
		{makeLit(124.0), "124"},
		{makeLit(124.56728482274629), "124.56728482274629"},
		{makeLit(true), "true"},
		{makeLit(false), "false"},
		{makeLit(errors.New("This is an error")), "This is an error"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("String() with Literal: %q - return: %q, want: %q",
				c.in, got, c.want)
		}
	}
}

/***********************
   Helper functions
************************/
func makeLit(v interface{}) *Literal {
	l, err := New(v)
	if err != nil {
		panic(fmt.Sprintf("MakeLit(%q) gave err: %q", v, err))
	}
	return l
}

func errorMatch(e1 error, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Error() == e2.Error() {
		return true
	}
	return false
}
