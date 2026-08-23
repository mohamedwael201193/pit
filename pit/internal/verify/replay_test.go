package verify

import "testing"

func TestReceiptReplay(t *testing.T) {
	u := NewUsed()
	h := "0x" + "ab"
	if err := u.File(h); err != nil {
		t.Fatal(err)
	}
	if err := u.File(h); err == nil {
		t.Fatal("replay")
	}
}
