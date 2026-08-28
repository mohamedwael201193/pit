package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestHomeEmptyIsReal(t *testing.T) {
	if Attention(0) != "No opportunities match your policy." {
		t.Fatal("empty")
	}
	cards := Home(nil, policy.Default())
	if len(cards) != 0 {
		t.Fatal("ghost")
	}
	b := BlockedCard("XYZ", "asset_not_allowed")
	if b.Eligible {
		t.Fatal("blocked")
	}
}
