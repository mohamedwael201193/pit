// Package venue holds Hyperliquid perp constraints that the host must not invent.
// The documented perp floor is $10 notional. A market's size tick can raise that floor.
package venue

import (
	"fmt"
	"math"
	"strings"
)

// HyperliquidPerpFloorUSD is the documented Hyperliquid perpetual minimum order value.
// It is not invented. Individual markets may require more because of szDecimals.
const HyperliquidPerpFloorUSD = 10.0

func MinSize(szDecimals int) float64 {
	if szDecimals < 0 {
		szDecimals = 0
	}
	if szDecimals > 8 {
		szDecimals = 8
	}
	return math.Pow(10, -float64(szDecimals))
}

// PerpMinNotionalUSD is the smallest sz*mark that meets the $10 Hyperliquid floor
// after rounding size up to szDecimals. One tick below $10 is not a legal order.
func PerpMinNotionalUSD(markPx float64, szDecimals int) float64 {
	floor := HyperliquidPerpFloorUSD
	if markPx <= 0 || math.IsNaN(markPx) || math.IsInf(markPx, 0) {
		return floor
	}
	if szDecimals < 0 {
		szDecimals = 0
	}
	if szDecimals > 8 {
		szDecimals = 8
	}
	pow := math.Pow(10, float64(szDecimals))
	sz := math.Ceil((floor/markPx)*pow) / pow
	if sz <= 0 {
		return floor
	}
	got := sz * markPx
	if got < floor {
		return floor
	}
	return got
}

func WhyThisMarket(coin string, available, minNotional float64) string {
	coin = strings.ToUpper(strings.TrimSpace(coin))
	if coin == "" {
		coin = "This market"
	}
	short := minNotional - available
	if short < 0 {
		short = 0
	}
	return fmt.Sprintf(
		"Your account has $%.2f available. %s requires at least $%.2f notional on Hyperliquid. This candidate is blocked. PIT will not invent size.",
		available, coin, minNotional,
	)
}
