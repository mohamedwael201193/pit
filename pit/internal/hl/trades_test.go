package hl

import (
	"encoding/json"
	"testing"
)

func TestParseTrades(t *testing.T) {
	raw := json.RawMessage(`[{"coin":"ETH","px":"3501.5","sz":"0.02","time":1710000000000,"side":"B"}]`)
	ts, err := ParseTrades(raw)
	if err != nil || len(ts) != 1 || ts[0].Px != 3501.5 {
		t.Fatalf("%v %v", ts, err)
	}
}

func TestParseTradesEmpty(t *testing.T) {
	if _, err := ParseTrades(json.RawMessage(`[]`)); err == nil {
		t.Fatal("empty")
	}
}
