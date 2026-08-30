package companion

import (
	"strings"
	"testing"
)

func TestHexHashIsNotAKey(t *testing.T) {
	digest := "0x" + strings.Repeat("ab", 32)
	if secretful(digest) {
		t.Fatal("0x+64 is a digest, not a session key")
	}
	if !looksLikeHexKey(strings.Repeat("ab", 32)) {
		t.Fatal("bare 64 hex still refused")
	}
	if !secretful("app-sk-demo") {
		t.Fatal("router key")
	}
}
