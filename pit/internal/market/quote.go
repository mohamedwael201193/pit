package market

import (
	"fmt"
	"strings"
	"time"
)

type Quote struct {
	Source    string    `json:"source"`
	Network   string    `json:"network"`
	Symbol    string    `json:"symbol"`
	MarkPx    float64   `json:"markPx,omitempty"`
	OraclePx  float64   `json:"oraclePx,omitempty"`
	Funding   float64   `json:"funding,omitempty"`
	OpenInterest float64 `json:"openInterest,omitempty"`
	AsOf      time.Time `json:"asOf"`
	Ref       string    `json:"ref,omitempty"`
}

func Validate(q Quote) error {
	if strings.TrimSpace(q.Source) == "" {
		return fmt.Errorf("missing_source")
	}
	if q.AsOf.IsZero() {
		return fmt.Errorf("missing_timestamp")
	}
	if q.MarkPx < 0 || q.OraclePx < 0 {
		return fmt.Errorf("negative_px")
	}
	return nil
}

func Hyperliquid(network, coin string, mark, oracle, funding, oi float64, asOf time.Time) Quote {
	return Quote{
		Source:       "hyperliquid",
		Network:      network,
		Symbol:       coin,
		MarkPx:       mark,
		OraclePx:     oracle,
		Funding:      funding,
		OpenInterest: oi,
		AsOf:         asOf,
		Ref:          "info/metaAndAssetCtxs",
	}
}
