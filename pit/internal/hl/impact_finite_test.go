package hl

import (
	"math"
	"testing"
)

func TestFiniteImpact(t *testing.T) {
	if err := FiniteImpact(Impact{Buy: 2501, Sell: 2498}); err != nil {
		t.Fatal(err)
	}
	if err := FiniteImpact(Impact{Buy: math.Inf(1), Sell: 2498}); err == nil {
		t.Fatal("inf")
	}
}
