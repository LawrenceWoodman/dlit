package dlit

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		in        interface{}
		want      *Literal
		wantError error
	}{
		{6, MustNew(6), nil},
		{6.0, MustNew(6.0), nil},
		{6.6, MustNew(6.6), nil},
		{float32(6.6), MustNew(float32(6.6)), nil},
		{int64(922336854775807), MustNew(922336854775807), nil},
		{int64(9223372036854775807), MustNew(9223372036854775807), nil},
		{"98292223372036854775807", MustNew("98292223372036854775807"), nil},
		{complex64(1), MustNew(InvalidKindError("complex64")),
			InvalidKindError("complex64")},
		{complex128(1), MustNew(InvalidKindError("complex128")),
			InvalidKindError("complex128")},
		{"6", MustNew("6"), nil},
		{"6.6", MustNew("6.6"), nil},
		{"abc", MustNew("abc"), nil},
		{true, MustNew(true), nil},
		{false, MustNew(false), nil},
		{errors.New("This is an error"), MustNew(errors.New("This is an error")),
			nil},
	}

	for _, c := range cases {
		got, err := New(c.in)
		if !errorMatch(err, c.wantError) {
			t.Errorf("New(%q) - err == %q, wantError == %q", c.in, err, c.wantError)
		}

		if got.String() != c.want.String() {
			t.Errorf("New(%q) - got == %s, want == %s", c.in, got, c.want)
		}
	}
}

func TestNewString(t *testing.T) {
	cases := []struct {
		in   string
		want *Literal
	}{
		{"", MustNew("")},
		{"6", MustNew("6")},
		{"6.27", MustNew("6.27")},
		{"Hello how are you today", MustNew("Hello how are you today")},
	}

	for _, c := range cases {
		got := NewString(c.in)
		if got.String() != c.want.String() {
			t.Errorf("New(%q) - got == %s, want == %s", c.in, got, c.want)
		}
	}
}

func TestMustNew(t *testing.T) {
	cases := []struct {
		in   interface{}
		want *Literal
	}{
		{6, MustNew(6)},
		{6.0, MustNew(6.0)},
		{6.6, MustNew(6.6)},
		{float32(6.6), MustNew(float32(6.6))},
		{int64(922336854775807), MustNew(922336854775807)},
		{int64(9223372036854775807), MustNew(9223372036854775807)},
		{"98292223372036854775807", MustNew("98292223372036854775807")},
		{"6", MustNew("6")},
		{"6.6", MustNew("6.6")},
		{"abc", MustNew("abc")},
		{true, MustNew(true)},
		{false, MustNew(false)},
		{errors.New("This is an error"), MustNew(errors.New("This is an error"))},
	}

	for _, c := range cases {
		got := MustNew(c.in)
		if got.String() != c.want.String() {
			t.Errorf("MustNew(%q) - got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestMustNew_panic(t *testing.T) {
	cases := []struct {
		in        interface{}
		wantPanic string
	}{
		{6, ""},
		{complex64(1), InvalidKindError("complex64").Error()},
	}

	for _, c := range cases {
		paniced := false
		defer func() {
			if r := recover(); r != nil {
				if r.(string) == c.wantPanic {
					paniced = true
				} else {
					t.Errorf("MustNew(%q) - got panic: %s, wanted: %s",
						c.in, r, c.wantPanic)
				}
			}
		}()
		MustNew(c.in)
		if c.wantPanic != "" && !paniced {
			t.Errorf("MustNew(%q) - failed to panic with: %s", c.in, c.wantPanic)
		}
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      int64
		wantIsInt bool
	}{
		{MustNew(6), 6, true},
		{MustNew(6.0), 6, true},
		{MustNew(float32(6.0)), 6, true},
		{MustNew("6"), 6, true},
		{MustNew("6.0"), 6, true},
		{MustNew(fmt.Sprintf("%d", int64(math.MinInt64))),
			int64(math.MinInt64), true},
		{MustNew(fmt.Sprintf("%d", int64(math.MaxInt64))),
			int64(math.MaxInt64), true},
		{MustNew(fmt.Sprintf("-1%d", int64(math.MinInt64))), 0, false},
		{MustNew(fmt.Sprintf("1%d", int64(math.MaxInt64))), 0, false},
		{MustNew("-9223372036854775809"), 0, false},
		{MustNew("9223372036854775808"), 0, false},
		{MustNew(6.6), 0, false},
		{MustNew("6.6"), 0, false},
		{MustNew("abc"), 0, false},
		{MustNew(true), 0, false},
		{MustNew(false), 0, false},
		{MustNew(errors.New("This is an error")), 0, false},
	}

	for _, c := range cases {
		got, gotIsInt := c.in.Int()
		if got != c.want || gotIsInt != c.wantIsInt {
			t.Errorf("Int() with Literal: %s - return: %d, %t - want: %d, %t",
				c.in, got, gotIsInt, c.want, c.wantIsInt)
		}
	}
}

func TestFloat(t *testing.T) {
	cases := []struct {
		in          *Literal
		want        float64
		wantIsFloat bool
	}{
		{MustNew(6), 6.0, true},
		{MustNew(int64(922336854775807)), 922336854775807.0, true},
		{MustNew(fmt.Sprintf("%G", float64(math.SmallestNonzeroFloat64))),
			math.SmallestNonzeroFloat64, true},
		{MustNew(fmt.Sprintf("%f", float64(math.MaxFloat64))),
			float64(math.MaxFloat64), true},
		{MustNew(6.0), 6.0, true},
		{MustNew("6"), 6.0, true},
		{MustNew(6.678934), 6.678934, true},
		{MustNew("6.678394"), 6.678394, true},
		{MustNew("abc"), 0, false},
		{MustNew(true), 0, false},
		{MustNew(false), 0, false},
		{MustNew(errors.New("This is an error")), 0, false},
	}

	for _, c := range cases {
		got, gotIsFloat := c.in.Float()
		if got != c.want || gotIsFloat != c.wantIsFloat {
			t.Errorf("Float() with Literal: %s - return: %f, %t - want: %f, %t",
				c.in, got, gotIsFloat, c.want, c.wantIsFloat)
		}
	}
}

func TestBool(t *testing.T) {
	cases := []struct {
		in         *Literal
		want       bool
		wantIsBool bool
	}{
		{MustNew(1), true, true},
		{MustNew(2), false, false},
		{MustNew(0), false, true},
		{MustNew(1.0), true, true},
		{MustNew(2.0), false, false},
		{MustNew(2.25), false, false},
		{MustNew(0.0), false, true},
		{MustNew(true), true, true},
		{MustNew(false), false, true},
		{MustNew("true"), true, true},
		{MustNew("false"), false, true},
		{MustNew("True"), true, true},
		{MustNew("False"), false, true},
		{MustNew("TRUE"), true, true},
		{MustNew("FALSE"), false, true},
		{MustNew("t"), true, true},
		{MustNew("f"), false, true},
		{MustNew("T"), true, true},
		{MustNew("F"), false, true},
		{MustNew("1"), true, true},
		{MustNew("0"), false, true},
		{MustNew("bob"), false, false},
		{MustNew(errors.New("This is an error")), false, false},
	}

	for _, c := range cases {
		got, gotIsBool := c.in.Bool()
		if got != c.want || gotIsBool != c.wantIsBool {
			t.Errorf("Bool() with Literal: %s - return: %t, %t - want: %t, %t",
				c.in, got, gotIsBool, c.want, c.wantIsBool)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in   *Literal
		want string
	}{
		{MustNew(124), "124"},
		{MustNew(int64(922336854775807)), "922336854775807"},
		{MustNew(int64(9223372036854775807)), "9223372036854775807"},
		{MustNew("98292223372036854775807"), "98292223372036854775807"},
		{MustNew("Hello my name is fred"), "Hello my name is fred"},
		{MustNew(124.0), "124"},
		{MustNew(124.56728482274629), "124.56728482274629"},
		{MustNew(true), "true"},
		{MustNew(false), "false"},
		{MustNew(errors.New("This is an error")), "This is an error"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("String() with Literal: %q - return: %q, want: %q",
				c.in, got, c.want)
		}
	}
}

func TestErr(t *testing.T) {
	cases := []struct {
		in   *Literal
		want error
	}{
		{MustNew(1), nil},
		{MustNew(2), nil},
		{MustNew("true"), nil},
		{MustNew(2.25), nil},
		{MustNew("hello"), nil},
		{MustNew(errors.New("This is an error")), errors.New("This is an error")},
	}

	for _, c := range cases {
		got := c.in.Err()
		if !errorMatch(c.want, got) {
			t.Errorf("Err() with Literal: %s - got: %s, want: %s", c.in, got, c.want)
		}
	}
}

/***********************
   Helper functions
************************/
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
