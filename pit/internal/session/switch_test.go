package session

import "testing"

func TestSwitchAccount(t *testing.T) {
	if err := SwitchAccount("0xabc", "0xabc"); err != nil {
		t.Fatal(err)
	}
	if err := SwitchAccount("0xabc", "0xdef"); err == nil {
		t.Fatal("switch")
	}
}
