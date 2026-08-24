package hl

import (
	"math"
	"testing"
)

func TestOpenInterestFinite(t *testing.T) {
	if err := OpenInterestFinite(BookSnapshot{OpenInterest: 1}); err != nil {
		t.Fatal(err)
	}
	if err := OpenInterestFinite(BookSnapshot{OpenInterest: math.Inf(1)}); err == nil {
		t.Fatal("inf")
	}
}
