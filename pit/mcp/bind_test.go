package mcp

import "testing"

func TestBindRejectsOtherWorkspace(t *testing.T) {
	a := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	b := "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	if err := Bind(a, b); err == nil {
		t.Fatal("cross")
	}
}
