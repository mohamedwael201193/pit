package hl

import "testing"

func TestValidCloid(t *testing.T) {
	if err := ValidCloid("0x11111111111111111111111111111111"); err != nil {
		t.Fatal(err)
	}
	if err := ValidCloid("0x1"); err == nil {
		t.Fatal("short")
	}
	if err := ValidCloid("11111111111111111111111111111111"); err == nil {
		t.Fatal("prefix")
	}
}
