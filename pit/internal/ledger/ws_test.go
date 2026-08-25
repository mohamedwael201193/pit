package ledger

import "testing"

func TestRefuseEmptyWorkspace(t *testing.T) {
	if err := RefuseEmptyWorkspace(""); err == nil {
		t.Fatal("ws")
	}
	if err := RefuseEmptyWorkspace("abc"); err != nil {
		t.Fatal(err)
	}
}
