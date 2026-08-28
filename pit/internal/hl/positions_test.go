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

func TestParseClearinghouse(t *testing.T) {
	raw := []byte(`{"marginSummary":{"accountValue":"54.23","totalNtlPos":"10.2","totalMarginUsed":"10.2"},"withdrawable":"44"}`)
	got := ParseClearinghouse(raw)
	if got.AccountValue != "54.23" || got.Withdrawable != "44" {
		t.Fatalf("%+v", got)
	}
}

func TestParsePositionsLeverage(t *testing.T) {
	raw := []byte(`{"assetPositions":[{"position":{"coin":"ETH","szi":"0.0041","entryPx":"2489.7","marginUsed":"10","leverage":{"type":"cross","value":"1"}}}]}`)
	got := ParsePositions(raw)
	if len(got) != 1 || got[0].Leverage != "1" || got[0].MarginUsed != "10" {
		t.Fatalf("%+v", got)
	}
}
