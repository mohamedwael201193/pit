package mcp

import "testing"

func TestKeyNever(t *testing.T) {
	if !KeyNever() {
		t.Fatal("key")
	}
}
