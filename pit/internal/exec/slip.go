package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func RefuseSlippage(p policy.Policy, bps int) error {
	if p.MaxSlippageBps > 0 && bps > p.MaxSlippageBps {
		return fmt.Errorf("slippage")
	}
	return nil
}
