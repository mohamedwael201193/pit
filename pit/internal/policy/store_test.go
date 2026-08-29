package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClampLocksLeverageAndVenue(t *testing.T) {
	p := Default()
	p.MaxLeverage = 20
	p.MaxClipUSD = 1000
	p.AllowedVenues = []string{"somewhere"}
	p.AllowedMarketTypes = []string{"spot"}
	p.AllowedAssets = []string{"XYZ", "eth", "ETH"}
	p.MaxOpenPositions = 99
	got := Clamp(p)
	if got.MaxLeverage != 1 {
		t.Fatalf("leverage %d", got.MaxLeverage)
	}
	if got.MaxClipUSD != ClipCeilUSD {
		t.Fatalf("clip %v", got.MaxClipUSD)
	}
	if got.AllowedVenues[0] != "hyperliquid" {
		t.Fatal(got.AllowedVenues)
	}
	if got.AllowedMarketTypes[0] != "perp" {
		t.Fatal(got.AllowedMarketTypes)
	}
	if len(got.AllowedAssets) != 1 || got.AllowedAssets[0] != "ETH" {
		t.Fatalf("assets %v", got.AllowedAssets)
	}
	if got.MaxOpenPositions != OpenCeil {
		t.Fatalf("open %d", got.MaxOpenPositions)
	}
}

func TestSaveLoadRoundTripAndTamper(t *testing.T) {
	dir := t.TempDir()
	ws := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	p := Default()
	p.MaxClipUSD = 12
	p.MaxOpenPositions = 2
	saved, hash, err := Save(dir, ws, p)
	if err != nil {
		t.Fatal(err)
	}
	if hash == "" || saved.MaxClipUSD != 12 || saved.MaxOpenPositions != 2 {
		t.Fatalf("%+v %s", saved, hash)
	}
	got := Load(dir, ws)
	if got.MaxClipUSD != 12 || got.MaxOpenPositions != 2 {
		t.Fatalf("load %+v", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "policy.json"), []byte(`{"maxClipUsd":40,"allowedAssets":["ETH"]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	tampered := Load(dir, ws)
	if tampered.MaxClipUSD != Default().MaxClipUSD {
		t.Fatalf("tamper must fail closed, got %+v", tampered)
	}
}

func TestChatCannotRaiseClipViaClamp(t *testing.T) {
	p := Default()
	p.MaxClipUSD = 1e9
	p.MaxLeverage = 50
	got := Clamp(p)
	if got.MaxClipUSD > ClipCeilUSD || got.MaxLeverage != 1 {
		t.Fatalf("%+v", got)
	}
}

func TestExecWhyOpenCeiling(t *testing.T) {
	p := Default()
	block, why := ExecWhy(1, 50, p)
	if block != "max_open_positions" || !strings.Contains(why, "ceiling") {
		t.Fatalf("%s %s", block, why)
	}
	block, _ = ExecWhy(0, 1, p)
	if block != "insufficient_margin" {
		t.Fatalf("%s", block)
	}
	block, why = ExecWhy(0, 50, p)
	if block != "" || why != "" {
		t.Fatalf("%s %s", block, why)
	}
	block, why = ExecWhy(0, 9.38, p)
	if block != "insufficient_margin" || !strings.Contains(why, "$9.38") || !strings.Contains(why, "$0.62") {
		t.Fatalf("shortfall %s %s", block, why)
	}
}

func TestDiffMentionsOpenCeiling(t *testing.T) {
	a := Default()
	b := Default()
	b.MaxOpenPositions = 2
	lines := Diff(a, b)
	joined := strings.Join(lines, " ")
	if !strings.Contains(joined, "Open position ceiling") {
		t.Fatal(joined)
	}
	if !strings.Contains(joined, "Leverage stays 1x") {
		t.Fatal(joined)
	}
}

func TestHashIgnoresOpenCeiling(t *testing.T) {
	a := Default()
	b := Default()
	b.MaxOpenPositions = 2
	ha, err := a.Hash()
	if err != nil {
		t.Fatal(err)
	}
	hb, err := b.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if ha != hb {
		t.Fatal("open ceiling must not change pin identity")
	}
}
