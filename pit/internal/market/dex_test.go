package market

import (
	"testing"
	"time"
)

func TestOptionalDexScreener(t *testing.T) {
	if err := OptionalDexScreener(Quote{Source: "hyperliquid", AsOf: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err := OptionalDexScreener(Quote{Source: "dexscreener"}); err == nil {
		t.Fatal("ts")
	}
	if err := OptionalDexScreener(Quote{Source: "dexscreener", AsOf: time.Now()}); err != nil {
		t.Fatal(err)
	}
}
