package mcp

import "testing"

func TestOrderNever(t *testing.T) {
	if !OrderNever() {
		t.Fatal("order")
	}
}
