package hl

import (
	"fmt"
	"math"
)

func OpenInterestFinite(b BookSnapshot) error {
	if math.IsNaN(b.OpenInterest) || math.IsInf(b.OpenInterest, 0) || b.OpenInterest < 0 {
		return fmt.Errorf("bad_oi")
	}
	return nil
}
