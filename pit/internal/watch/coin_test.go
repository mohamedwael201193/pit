package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestCoinAllowed(t *testing.T) {
	p := policy.Default()
	if CoinAllowed(p, "XYZ") {
		t.Fatal("xyz")
	}
	if !CoinAllowed(p, "ETH") {
		t.Fatal("eth")
	}
}
