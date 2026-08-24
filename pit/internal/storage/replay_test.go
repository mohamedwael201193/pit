package storage

import "testing"

func TestRefuseRootReplay(t *testing.T) {
	if err := RefuseRootReplay("0xaaa", "0xbbb"); err == nil {
		t.Fatal("replay")
	}
	if err := RefuseRootReplay("0xaaa", "0xaaa"); err != nil {
		t.Fatal(err)
	}
}
