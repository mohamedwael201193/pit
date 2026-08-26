package storage

import "testing"

func TestRefuseGlobalMemoryKey(t *testing.T) {
	if err := RefuseGlobalMemoryKey("0xabc"); err == nil {
		t.Fatal("global")
	}
	if err := RefuseGlobalMemoryKey(""); err != nil {
		t.Fatal(err)
	}
}
