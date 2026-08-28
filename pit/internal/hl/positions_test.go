package hl

import "testing"

func TestParsePositionsETHLong(t *testing.T) {
	raw := []byte(`{"assetPositions":[{"position":{"coin":"ETH","szi":"0.0041","entryPx":"2489.7","unrealizedPnl":"0.01"}}]}`)
	got := ParsePositions(raw)
	if len(got) != 1 || got[0].Coin != "ETH" || got[0].Sz != "0.0041" {
		t.Fatalf("%+v", got)
	}
}

func TestParsePositionsEmpty(t *testing.T) {
	if ParsePositions([]byte(`{"assetPositions":[]}`)) == nil {
		t.Fatal("empty slice")
	}
	if len(ParsePositions([]byte(`{"assetPositions":[]}`))) != 0 {
		t.Fatal("len")
	}
}

func TestParsePositionsLeverage(t *testing.T) {
	raw := []byte(`{"assetPositions":[{"position":{"coin":"ETH","szi":"0.0041","entryPx":"2489.7","marginUsed":"10","leverage":{"type":"cross","value":"1"}}}]}`)
	got := ParsePositions(raw)
	if len(got) != 1 || got[0].Leverage != "1" || got[0].MarginUsed != "10" {
		t.Fatalf("%+v", got)
	}
}
