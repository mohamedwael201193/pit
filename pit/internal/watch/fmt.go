package watch

import (
	"fmt"
	"math"
)

func Compact(n float64) string {
	a := math.Abs(n)
	switch {
	case a >= 1e9:
		return fmt.Sprintf("%.2fB", n/1e9)
	case a >= 1e6:
		return fmt.Sprintf("%.2fM", n/1e6)
	case a >= 1000:
		return fmt.Sprintf("%.0f", n)
	case a >= 1:
		return fmt.Sprintf("%.4f", n)
	case a == 0:
		return "0"
	default:
		return fmt.Sprintf("%.6f", n)
	}
}

func CompactUSD(n float64) string {
	if n == 0 {
		return "0"
	}
	return "$" + Compact(n)
}

func Price(n float64) string {
	a := math.Abs(n)
	switch {
	case a >= 1000:
		return fmt.Sprintf("%.0f", n)
	case a >= 1:
		return fmt.Sprintf("%.2f", n)
	default:
		return fmt.Sprintf("%.6f", n)
	}
}

func FundingPct(n float64) string {
	return fmt.Sprintf("%.4f%%", n*100)
}
