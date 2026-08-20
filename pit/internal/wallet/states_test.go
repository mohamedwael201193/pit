package wallet

import "testing"

func TestNamedStates(t *testing.T) {
	if err := RejectSeedField(true); err == nil || err.Error() != SeedForbidden {
		t.Fatal("seed")
	}
	if err := RejectSeedField(false); err != nil {
		t.Fatal(err)
	}
	if err := MapChain(1, 16661); err == nil || err.Error() != WrongNetwork {
		t.Fatal("chain")
	}
	if err := MapSignature(false); err == nil || err.Error() != SignatureDeclined {
		t.Fatal("sig")
	}
}
