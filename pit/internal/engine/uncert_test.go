package engine

import "testing"

func TestRefuseUncertainty(t *testing.T) {
	if err := RefuseUncertainty(0.2, 1); err != nil {
		t.Fatal(err)
	}
	if err := RefuseUncertainty(0.9, 0.3); err == nil {
		t.Fatal("unc")
	}
}
