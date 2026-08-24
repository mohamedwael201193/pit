package compute

import "testing"

func TestRefuseIndependentClaim(t *testing.T) {
	if err := RefuseIndependentClaim(EnvelopeOnly, "independent_providers"); err == nil {
		t.Fatal("honesty")
	}
	if err := RefuseIndependentClaim(Providers, "independent_providers"); err != nil {
		t.Fatal(err)
	}
}
