package hl

import (
	"fmt"
	"math"
)

func FiniteImpact(im Impact) error {
	for _, v := range []float64{im.Buy, im.Sell} {
		if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
			return fmt.Errorf("bad_impact")
		}
	}
	return nil
}
