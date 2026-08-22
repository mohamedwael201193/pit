package mcp

import "testing"

func TestMayMutate(t *testing.T) {
	if err := MayMutate("cancel"); err == nil {
		t.Fatal("cancel")
	}
	if err := MayMutate("opportunities"); err != nil {
		t.Fatal(err)
	}
}
