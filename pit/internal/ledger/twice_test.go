package ledger

import "testing"

func TestRefuseSecondApply(t *testing.T) {
	if err := RefuseSecondApply(false); err == nil {
		t.Fatal("dup")
	}
	if err := RefuseSecondApply(true); err != nil {
		t.Fatal(err)
	}
}
