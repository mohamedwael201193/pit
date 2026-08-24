package hl

import (
	"math"
	"testing"
)

func TestFundingFinite(t *testing.T) {
	if err := FundingFinite(BookSnapshot{Funding: 0.0001}); err != nil {
		t.Fatal(err)
	}
	if err := FundingFinite(BookSnapshot{Funding: math.NaN()}); err == nil {
		t.Fatal("nan")
	}
}
