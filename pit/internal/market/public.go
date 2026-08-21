package market

import (
	"fmt"
	"strings"
	"time"
)

func CoinGecko(symbol string, price float64, asOf time.Time, ref string) (Quote, error) {
	q := Quote{Source: "coingecko", Network: "public", Symbol: symbol, MarkPx: price, AsOf: asOf, Ref: ref}
	return q, Validate(q)
}

func DefiLlama(symbol string, tvl float64, asOf time.Time, ref string) (Quote, error) {
	if tvl < 0 {
		return Quote{}, fmt.Errorf("negative_tvl")
	}
	q := Quote{Source: "defillama", Network: "public", Symbol: symbol, OpenInterest: tvl, AsOf: asOf, Ref: ref}
	return q, Validate(q)
}

func RejectStale(q Quote, now time.Time, maxAge time.Duration) error {
	if err := Validate(q); err != nil {
		return err
	}
	if now.Sub(q.AsOf) > maxAge {
		return fmt.Errorf("stale_quote")
	}
	if strings.EqualFold(q.Source, "mock") {
		return fmt.Errorf("mock_quote_denied")
	}
	return nil
}
