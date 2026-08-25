package wallet

import "testing"

func TestRefuseStaleAddress(t *testing.T) {
	if err := RefuseStaleAddress("0xaaa", "0xbbb"); err == nil {
		t.Fatal("switch")
	}
	if err := RefuseStaleAddress("0xaaa", "0xAAA"); err != nil {
		t.Fatal(err)
	}
}
