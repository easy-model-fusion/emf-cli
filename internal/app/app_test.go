package app

import "testing"

func TestInit(t *testing.T) {
	Init("1.0.0", "2021-01-01")
	if Version != "1.0.0" {
		t.Errorf("Version should be 1.0.0, got %s", Version)
	}
	if BuildDate != "2021-01-01" {
		t.Errorf("BuildDate should be 2021-01-01, got %s", BuildDate)
	}
}
