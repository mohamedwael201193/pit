package market

import "testing"

func TestRefuseMock(t *testing.T) {
	if err := RefuseMock("hyperliquid"); err != nil {
		t.Fatal(err)
	}
	if err := RefuseMock("mock"); err == nil {
		t.Fatal("mock")
	}
}
