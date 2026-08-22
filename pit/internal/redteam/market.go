package redteam

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func MarketMutation(host engine.Preview, boundHash string, nowMs int64) error {
	mut := host
	if mut.Market == "ETH" {
		mut.Market = "BTC"
	} else {
		mut.Market = "ETH"
	}
	used := map[string]struct{}{}
	if err := engine.Authorize(mut, boundHash, nowMs, used); err == nil {
		return fmt.Errorf("market_mutation_accepted")
	}
	return nil
}
