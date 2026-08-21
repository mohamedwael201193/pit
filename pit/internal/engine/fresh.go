package engine

import (
	"fmt"
	"time"

	"github.com/mohamedwael201193/pit/internal/market"
)

func RequireFreshQuote(q market.Quote, now time.Time) error {
	return market.RejectStale(q, now, 30*time.Second)
}

func RejectWrongMarket(hostMarket, previewMarket string) error {
	if hostMarket == "" || previewMarket == "" || hostMarket != previewMarket {
		return fmt.Errorf("market_changed")
	}
	return nil
}
