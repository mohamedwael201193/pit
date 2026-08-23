package sdk

import "testing"

func TestKeyringOnWebDenied(t *testing.T) {
	if err := (Client{}).KeyringOnWeb(); err == nil {
		t.Fatal("web")
	}
}
