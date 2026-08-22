package engine

import (
	"fmt"
	"math"
)

func RejectNonFinite(v float64, name string) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("non_finite_%s", name)
	}
	return nil
}

func SizeOrderGuarded(in SizerInput) (SizedOrder, error) {
	if err := RejectNonFinite(in.MarkPx, "mark"); err != nil {
		return SizedOrder{}, err
	}
	if err := RejectNonFinite(in.RequestedUSD, "size"); err != nil {
		return SizedOrder{}, err
	}
	if err := RejectNonFinite(in.MaxClipUSD, "clip"); err != nil {
		return SizedOrder{}, err
	}
	return SizeOrder(in)
}
