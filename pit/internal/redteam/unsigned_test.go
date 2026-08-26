package redteam

import "testing"

func TestUnsignedExchange(t *testing.T) {
	if err := UnsignedExchange(false); err == nil {
		t.Fatal("unsigned")
	}
	if err := UnsignedExchange(true); err != nil {
		t.Fatal(err)
	}
}
