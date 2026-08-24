package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func RefuseLeverage(p policy.Policy, requested int) error {
	if requested > 0 && p.MaxLeverage > 0 && requested > p.MaxLeverage {
		return fmt.Errorf("leverage")
	}
	return nil
}
