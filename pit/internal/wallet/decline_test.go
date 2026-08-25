package wallet

import "testing"

func TestRefuseDeclinedBind(t *testing.T) {
	if err := RefuseDeclinedBind(false); err == nil {
		t.Fatal("declined")
	}
	if err := RefuseDeclinedBind(true); err != nil {
		t.Fatal(err)
	}
}
