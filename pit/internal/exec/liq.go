package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func RefuseLiquidity(p policy.Policy, impactUSD float64) error {
	if p.MinLiquidityUSD > 0 && impactUSD+1e-9 < p.MinLiquidityUSD {
		return fmt.Errorf("liquidity")
	}
	return nil
}
