package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func RefuseCooldown(p policy.Policy, lastFillUnix, nowUnix int64) error {
	if p.CooldownSeconds > 0 && lastFillUnix > 0 && nowUnix-lastFillUnix < p.CooldownSeconds {
		return fmt.Errorf("cooldown")
	}
	return nil
}
