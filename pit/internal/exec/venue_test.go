package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestRefuseVenue(t *testing.T) {
	p := policy.Default()
	if err := RefuseVenue(p, "binance"); err == nil {
		t.Fatal("venue")
	}
	if err := RefuseVenue(p, "hyperliquid"); err != nil {
		t.Fatal(err)
	}
}
