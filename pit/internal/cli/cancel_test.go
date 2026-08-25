package cli

import "testing"

func TestCancelOnly(t *testing.T) {
	if err := CancelOnly("cancel"); err != nil {
		t.Fatal(err)
	}
	if err := CancelOnly("withdraw3"); err == nil {
		t.Fatal("withdraw")
	}
}
