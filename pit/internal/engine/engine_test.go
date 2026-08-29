package engine

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

func TestSizerClip(t *testing.T) {
	got, err := SizeOrder(SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 10, RequestedUSD: 1_000_000,
		Side: "buy", Coin: "ETH", AllowedCoins: []string{"ETH"}, MaxLeverage: 1, RequestedLev: 1,
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.NotionalUSD > 10+1e-6 {
		t.Fatalf("%v", got)
	}
}

func TestSizerBelowMin(t *testing.T) {
	_, err := SizeOrder(SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 5, RequestedUSD: 5,
		Side: "sell", Coin: "ETH", AllowedCoins: []string{"ETH"},
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err == nil || err.Error() != "below_min_notional" {
		t.Fatalf("got %v", err)
	}
}

func TestEvaluateNoSide(t *testing.T) {
	rs, _ := json.Marshal(RoleJSON{ProposedSide: "none"})
	ch, _ := json.Marshal(map[string]any{"survives": true, "kill": false})
	got := EvaluateCommittee(2500, 4, 15, 15, "", "ETH", []string{"ETH"}, 1, 1, false, rs, ch, ch)
	if got.Eligible || got.Deny != "no_side" {
		t.Fatalf("%+v", got)
	}
}

func TestEvaluateChallengerStructured(t *testing.T) {
	falseV := false
	ch, _ := json.Marshal(RoleJSON{Survives: &falseV})
	rs, _ := json.Marshal(RoleJSON{ProposedSide: "sell"})
	got := Evaluate(2500, 4, 15, 15, "sell", "ETH", []string{"ETH"}, 1, 1, false, rs, ch)
	if got.Eligible || got.Deny != "challenger_killed" {
		t.Fatalf("%+v", got)
	}
}

func TestEvaluateDoesNotTreatSubstring(t *testing.T) {
	rs, _ := json.Marshal(map[string]any{"proposed_side": "buy", "note": "thesis_killed is a quote"})
	ch, _ := json.Marshal(map[string]any{"survives": true, "comment": "not thesis_killed as a flag"})
	got := Evaluate(2500, 4, 15, 15, "buy", "ETH", []string{"ETH"}, 1, 1, false, rs, ch)
	if !got.Eligible {
		t.Fatalf("substring must not kill: %+v", got)
	}
}

func TestSizerBTCMeetsVenueMin(t *testing.T) {
	got, err := SizeOrder(SizerInput{
		MarkPx: 80000, SzDecimals: 5, MaxClipUSD: 10, RequestedUSD: 10,
		Side: "buy", Coin: "BTC", AllowedCoins: []string{"BTC"},
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.NotionalUSD+1e-9 < 10 {
		t.Fatalf("%v", got)
	}
}

func TestEvaluateRiskKill(t *testing.T) {
	trueV := true
	rs, _ := json.Marshal(RoleJSON{ProposedSide: "buy"})
	ch, _ := json.Marshal(RoleJSON{Survives: &trueV})
	rk, _ := json.Marshal(RoleJSON{Kill: &trueV})
	got := EvaluateCommittee(2500, 4, 15, 15, "", "ETH", []string{"ETH"}, 1, 1, false, rs, ch, rk)
	if got.Eligible || got.Deny != "risk_killed" {
		t.Fatalf("%+v", got)
	}
}

func TestPreviewBindIgnoresModel(t *testing.T) {
	host := Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit",
		LimitPx: "1000", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
	}
	got, err := BindPreview(host, map[string]any{"sz": 99.0, "side": "sell"})
	if err != nil || got.Sz != host.Sz {
		t.Fatalf("%v %+v", err, got)
	}
	h, _ := CanonicalHash(host)
	used := map[string]struct{}{}
	if err := Authorize(host, h, time.Now().UnixMilli(), used); err != nil {
		t.Fatal(err)
	}
	used[host.Cloid] = struct{}{}
	if err := Authorize(host, h, time.Now().UnixMilli(), used); err == nil {
		t.Fatal("replay")
	}
	_ = math.NaN()
}

func TestSizerDoesNotPadAboveRequested(t *testing.T) {
	_, err := SizeOrder(SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 10, RequestedUSD: 9.38,
		Side: "buy", Coin: "ETH", AllowedCoins: []string{"ETH"}, MaxLeverage: 1, RequestedLev: 1,
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err == nil || err.Error() != "below_min_notional" {
		t.Fatalf("padded %v", err)
	}
}
