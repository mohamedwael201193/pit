package exec

import "testing"

func TestRefuseWithdraw(t *testing.T) {
	if err := RefuseWithdraw("withdraw3"); err == nil {
		t.Fatal("withdraw")
	}
	if err := RefuseWithdraw("order"); err != nil {
		t.Fatal(err)
	}
}
