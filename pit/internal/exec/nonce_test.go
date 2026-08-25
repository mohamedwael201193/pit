package exec

import "testing"

func TestBindNonce(t *testing.T) {
	if err := BindNonce(7, 8); err == nil {
		t.Fatal("nonce")
	}
	if err := BindNonce(7, 7); err != nil {
		t.Fatal(err)
	}
}
