package siwe

import "testing"

func TestDomainOK(t *testing.T) {
	if err := DomainOK("pit.local"); err != nil {
		t.Fatal(err)
	}
	if err := DomainOK(""); err == nil {
		t.Fatal("empty")
	}
}
