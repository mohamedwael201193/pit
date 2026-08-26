package mcp

import "testing"

func TestPostNever(t *testing.T) {
	if !PostNever() {
		t.Fatal("post")
	}
}
