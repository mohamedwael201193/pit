package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestBlockedAsset(t *testing.T) {
	p := policy.Default()
	c := Blocked("XYZ", p)
	if c.Eligible {
		t.Fatal("xyz")
	}
	eth := Blocked("ETH", p)
	if !eth.Eligible {
		t.Fatal("eth")
	}
}
