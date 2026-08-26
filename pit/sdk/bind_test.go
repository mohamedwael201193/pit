package sdk

import "testing"

func TestBindWorkspace(t *testing.T) {
	a := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	b := "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	if err := BindWorkspace(a, a); err != nil {
		t.Fatal(err)
	}
	if err := BindWorkspace(a, b); err == nil {
		t.Fatal("cross")
	}
}
