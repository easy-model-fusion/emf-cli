package test

import (
	"reflect"
	"strings"
	"testing"
)

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}, message ...string) {
	if a == b {
		return
	}

	var errorMessage string
	if len(message) != 0 {
		errorMessage = strings.Join(message, " ") + "\n"
	}

	t.Helper()
	t.Errorf("%sReceived %v (type %v), expected %v (type %v)", errorMessage, a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	t.FailNow()
}

// AssertNotEqual checks if values are not equal
func AssertNotEqual(t *testing.T, a interface{}, b interface{}, message ...string) {
	if a != b {
		return
	}

	var errorMessage string
	if len(message) != 0 {
		errorMessage = strings.Join(message, " ") + "\n"
	}

	t.Helper()
	t.Errorf("%sReceived %v (type %v), expected != %v (type %v)", errorMessage, a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	t.FailNow()
}
