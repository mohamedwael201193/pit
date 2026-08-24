package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestRefuseLiquidity(t *testing.T) {
	p := policy.Default()
	p.MinLiquidityUSD = 1000
	if err := RefuseLiquidity(p, 500); err == nil {
		t.Fatal("thin")
	}
	if err := RefuseLiquidity(p, 2000); err != nil {
		t.Fatal(err)
	}
}
