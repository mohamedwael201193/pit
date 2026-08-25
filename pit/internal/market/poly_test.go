package market

import (
	"testing"
	"time"
)

func TestOptionalPolymarket(t *testing.T) {
	if err := OptionalPolymarket(Quote{Source: "hyperliquid", AsOf: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err := OptionalPolymarket(Quote{Source: "polymarket"}); err == nil {
		t.Fatal("ts")
	}
}
