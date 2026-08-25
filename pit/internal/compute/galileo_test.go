package compute

import "testing"

func TestRefuseUnprovenGalileo(t *testing.T) {
	if err := RefuseUnprovenGalileo(); err == nil {
		t.Fatal("galileo")
	}
}
