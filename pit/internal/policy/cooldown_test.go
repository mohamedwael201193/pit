package policy

import "testing"

func TestInCooldown(t *testing.T) {
	if InCooldown(100, 170, 60) {
		t.Fatal("elapsed")
	}
	if !InCooldown(100, 120, 60) {
		t.Fatal("inside")
	}
	if InCooldown(0, 120, 60) {
		t.Fatal("no last fill")
	}
}
