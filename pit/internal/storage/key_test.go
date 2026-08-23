package storage

import (
	"strings"
	"testing"
)

func TestKeysMatch(t *testing.T) {
	a := "0x" + strings.Repeat("aa", 32)
	b := "0x" + strings.Repeat("bb", 32)
	if err := KeysMatch(a, a); err != nil {
		t.Fatal(err)
	}
	if err := KeysMatch(a, b); err == nil {
		t.Fatal("wrong")
	}
}
