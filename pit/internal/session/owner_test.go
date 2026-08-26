package session

import "testing"

func TestBindWorkspace(t *testing.T) {
	if err := BindWorkspace("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"); err != nil {
		t.Fatal(err)
	}
	if err := BindWorkspace("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"); err == nil {
		t.Fatal("cross")
	}
}
