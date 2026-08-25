package watch

import (
	"strings"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func CoinAllowed(p policy.Policy, coin string) bool {
	for _, a := range p.AllowedAssets {
		if strings.EqualFold(a, coin) {
			return true
		}
	}
	return false
}
