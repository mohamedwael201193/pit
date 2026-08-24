package hl

import (
	"fmt"
	"math"
)

func OracleFinite(b BookSnapshot) error {
	if math.IsNaN(b.OraclePx) || math.IsInf(b.OraclePx, 0) || b.OraclePx < 0 {
		return fmt.Errorf("bad_oracle")
	}
	return nil
}
