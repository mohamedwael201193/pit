package hl

import (
	"math"
	"testing"
)

func TestMarkFinite(t *testing.T) {
	if err := MarkFinite(BookSnapshot{MarkPx: 2500}); err != nil {
		t.Fatal(err)
	}
	if err := MarkFinite(BookSnapshot{MarkPx: math.Inf(1)}); err == nil {
		t.Fatal("inf")
	}
}
