package engine

import "testing"

func TestRefuseEmptySide(t *testing.T) {
	if err := RefuseEmptySide(""); err == nil {
		t.Fatal("empty")
	}
	if err := RefuseEmptySide("none"); err == nil {
		t.Fatal("none")
	}
	if err := RefuseEmptySide("buy"); err != nil {
		t.Fatal(err)
	}
}
