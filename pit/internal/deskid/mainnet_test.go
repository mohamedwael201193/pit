package deskid

import "testing"

func TestRefuseAristotleTransfer(t *testing.T) {
	if err := RefuseAristotleTransfer(); err == nil {
		t.Fatal("xfer")
	}
}
