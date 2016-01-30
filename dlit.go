/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package dlit

import (
	"fmt"
	"reflect"
	"strconv"
)

type Literal struct {
	i          int64
	f          float64
	s          string
	b          bool
	e          error
	canBeInt   canBeKind
	canBeFloat canBeKind
	canBeBool  canBeKind
	canBeError canBeKind
}

type canBeKind int

const (
	unknown canBeKind = iota
	yes
	no
)

func New(v interface{}) (*Literal, error) {
	switch e := v.(type) {
	case int:
		return &Literal{i: int64(e), canBeInt: yes}, nil
	case int64:
		return &Literal{i: int64(e), canBeInt: yes}, nil
	case float64:
		return &Literal{f: e, canBeFloat: yes}, nil
	case string:
		return &Literal{s: e}, nil
	case bool:
		return &Literal{b: e, canBeBool: yes}, nil
	case error:
		return &Literal{e: e, canBeInt: no, canBeFloat: no, canBeBool: no,
			canBeError: yes}, nil
	}
	return nil, ErrInvalidKind(reflect.TypeOf(v).String())
}

func (l *Literal) Int() (int64, error) {
	switch l.canBeInt {
	case yes:
		return l.i, nil
	case unknown:
		if l.canBeFloat == yes {
			possibleInt := int64(l.f)
			if l.f == float64(possibleInt) {
				l.canBeInt = yes
				l.i = possibleInt
				return possibleInt, nil
			}
		} else {
			i, err := strconv.ParseInt(l.String(), 10, 64)
			if err == nil {
				l.canBeInt = yes
				l.i = i
				return i, nil
			}
		}
	}
	l.canBeInt = no
	return 0, ErrInvalidCast{l, "int"}
}

func (l *Literal) Float() (float64, error) {
	switch l.canBeFloat {
	case yes:
		return l.f, nil
	case unknown:
		if l.canBeInt == yes {
			possibleFloat := float64(l.i)
			if l.i == int64(possibleFloat) {
				l.canBeFloat = yes
				return possibleFloat, nil
			}
		} else {
			f, err := strconv.ParseFloat(l.String(), 64)
			if err == nil {
				l.canBeFloat = yes
				l.f = f
				return f, nil
			}
		}
	}
	l.canBeFloat = no
	return 0, ErrInvalidCast{l, "float"}
}

func (l *Literal) Bool() (bool, error) {
	switch l.canBeBool {
	case yes:
		return l.b, nil
	case unknown:
		if l.canBeInt == yes {
			l.canBeBool = yes
			l.b = l.i != 0
			return l.b, nil
		} else if l.canBeFloat == yes {
			l.canBeBool = yes
			l.b = l.f != 0.0
			return l.b, nil
		} else {
			b, err := strconv.ParseBool(l.s)
			if err == nil {
				l.canBeBool = yes
				l.b = b
				return b, nil
			}
		}
	}
	l.canBeBool = no
	return false, ErrInvalidCast{l, "bool"}
}

func (l *Literal) String() string {
	if len(l.s) > 0 {
		return l.s
	}
	switch true {
	case l.canBeInt == yes:
		l.s = strconv.FormatInt(l.i, 10)
	case l.canBeFloat == yes:
		l.s = strconv.FormatFloat(l.f, 'f', -1, 64)
	case l.canBeBool == yes:
		if l.b {
			l.s = "true"
		} else {
			l.s = "false"
		}
	case l.canBeError == yes:
		l.s = l.e.Error()
	}
	return l.s
}

type ErrInvalidCast struct {
	fromLiteral *Literal
	toType      string
}

type ErrInvalidKind string

func (e ErrInvalidCast) Error() string {
	return fmt.Sprintf("Can't cast: %s, to type: %s", e.fromLiteral, e.toType)
}

func (e ErrInvalidKind) Error() string {
	return fmt.Sprintf("Can't create Literal from type: %s", string(e))
}
