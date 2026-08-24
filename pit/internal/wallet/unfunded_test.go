package wallet

import "testing"

func TestMapUnfunded(t *testing.T) {
	if err := MapUnfunded(true); err != nil {
		t.Fatal(err)
	}
	if err := MapUnfunded(false); err == nil {
		t.Fatal("unfunded")
	}
}
