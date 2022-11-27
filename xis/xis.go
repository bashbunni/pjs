package xis

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/matryer/is"
)

type I interface {
	Contains(interface{}, interface{})
	Subset(interface{}, interface{})
}

type x struct {
	is *is.I
}

func New(i *is.I) I {
	return &x{is: i}
}

func (i *x) Contains(s, contains interface{}) {
	i.is.Helper()
	ok, found := includeElement(s, contains)
	i.is.True(ok && found)
}

func (i *x) Subset(list, subset interface{}) {
	eval := true

	defer func() {
		i.is.True(eval)
	}()

	if subset == nil {
		eval = true // we consider nil to be equal to the nil set
		return
	}

	subsetValue := reflect.ValueOf(subset)
	if e := recover(); e != nil {
		eval = false
		return
	}

	listKind := reflect.TypeOf(list).Kind()
	subsetKind := reflect.TypeOf(subset).Kind()

	if listKind != reflect.Array && listKind != reflect.Slice {
		eval = false
		return
	}

	if subsetKind != reflect.Array && subsetKind != reflect.Slice {
		eval = false
		return
	}

	for i := 0; i < subsetValue.Len(); i++ {
		element := subsetValue.Index(i).Interface()
		ok, found := includeElement(list, element)
		if !ok {
			eval = false
			return
		}
		if !found {
			eval = false
			return
		}
	}
}

func includeElement(list interface{}, element interface{}) (ok, found bool) {

	listValue := reflect.ValueOf(list)
	listKind := reflect.TypeOf(list).Kind()
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if listKind == reflect.String {
		elementValue := reflect.ValueOf(element)
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if listKind == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false

}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
