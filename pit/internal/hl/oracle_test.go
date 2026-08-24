package hl

import (
	"math"
	"testing"
)

func TestOracleFinite(t *testing.T) {
	if err := OracleFinite(BookSnapshot{OraclePx: 2500}); err != nil {
		t.Fatal(err)
	}
	if err := OracleFinite(BookSnapshot{OraclePx: math.NaN()}); err == nil {
		t.Fatal("nan")
	}
}
