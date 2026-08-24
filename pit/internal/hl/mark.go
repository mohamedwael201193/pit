package hl

import (
	"fmt"
	"math"
)

func MarkFinite(b BookSnapshot) error {
	if math.IsNaN(b.MarkPx) || math.IsInf(b.MarkPx, 0) || b.MarkPx <= 0 {
		return fmt.Errorf("bad_mark")
	}
	return nil
}
