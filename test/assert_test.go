package test

import "testing"

func TestAssertEqual(t *testing.T) {
	// Test with integers
	AssertEqual(t, 1, 1)

	// Test with strings
	AssertEqual(t, "test", "test")

	// Test with booleans
	AssertEqual(t, true, true)

	// Test with nil
	AssertEqual(t, nil, nil)
}

func TestAssertNotEqual(t *testing.T) {
	// Test with integers
	AssertNotEqual(t, 1, 2)

	// Test with strings
	AssertNotEqual(t, "test", "test2")

	// Test with booleans
	AssertNotEqual(t, true, false)

	// Test with nil
	AssertNotEqual(t, nil, 1)

	// Test with different types
	AssertNotEqual(t, 1, "test")
}
