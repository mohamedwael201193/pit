package hl

import (
	"fmt"
	"math"
)

func FundingFinite(b BookSnapshot) error {
	if math.IsNaN(b.Funding) || math.IsInf(b.Funding, 0) {
		return fmt.Errorf("bad_funding")
	}
	return nil
}
