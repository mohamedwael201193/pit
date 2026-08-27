package hl

import (
	"encoding/hex"
	"testing"
)

func TestActionHashMatchesPythonSDK(t *testing.T) {
	raw, err := BuildOrder(1, true, "2500", "0.004", "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	h, err := ActionHash(raw, 1700000000000)
	if err != nil {
		t.Fatal(err)
	}
	want := "29a3e02331f9dbf15b8886f40532766662ad4f8fe889ab3574595c1e5e7eb790"
	if hex.EncodeToString(h[:]) != want {
		t.Fatalf("got %s want %s packed=%s", hex.EncodeToString(h[:]), want, mustPack(t, raw))
	}
}

func mustPack(t *testing.T, raw []byte) string {
	t.Helper()
	b, err := packAction(raw)
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(b)
}
