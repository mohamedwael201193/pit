package ledger

import "testing"

func TestAfterTimeout(t *testing.T) {
	if err := AfterTimeout(false); err == nil {
		t.Fatal("blind")
	}
	if err := AfterTimeout(true); err != nil {
		t.Fatal(err)
	}
}
