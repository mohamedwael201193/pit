package venue

import (
	"strings"
	"testing"
)

func TestETHTickBelowDocumentedFloor(t *testing.T) {
	got := PerpMinNotionalUSD(2500, 4)
	if got != HyperliquidPerpFloorUSD {
		t.Fatalf("%v", got)
	}
}

func TestCoarseTickRaisesFloor(t *testing.T) {
	got := PerpMinNotionalUSD(50, 0)
	if got < 50-1e-9 {
		t.Fatalf("want market min 50 got %v", got)
	}
}

func TestBTCRoundedMinExceedsDocumentedFloor(t *testing.T) {
	got := PerpMinNotionalUSD(80000, 5)
	if got < 10.39 || got > 10.41 {
		t.Fatalf("BTC szDecimals=5 at 80000 must round the $10 floor up, got %v", got)
	}
}

func TestHighMarkBTCNeedsEleven(t *testing.T) {
	got := PerpMinNotionalUSD(110000, 5)
	if got < 10.99 || got > 11.01 {
		t.Fatalf("want ~11 got %v", got)
	}
}

func TestWhyThisMarketNamesTheAsset(t *testing.T) {
	s := WhyThisMarket("BTC", 9.38, 10)
	if !strings.Contains(s, "BTC") || !strings.Contains(s, "9.38") || !strings.Contains(s, "10.00") {
		t.Fatal(s)
	}
}

func TestWhyPolicyTightDoesNotAskToFund(t *testing.T) {
	s := WhyPolicyTight("ETH", 10, 10.08, 16.18)
	if !strings.Contains(s, "ETH") || !strings.Contains(s, "0.08") || !strings.Contains(s, "16.18") {
		t.Fatal(s)
	}
	if strings.Contains(strings.ToLower(s), "fund") {
		t.Fatal(s)
	}
}
