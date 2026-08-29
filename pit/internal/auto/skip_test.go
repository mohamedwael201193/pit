package auto

import (
	"strings"
	"testing"
	"time"
)

func TestStandDownDoesNotLatch(t *testing.T) {
	kind, hold, latch := ClassifyResearchSkip("", "no_side", false)
	if latch || kind != SkipStoodDown || hold < time.Minute {
		t.Fatalf("%s %v %v", kind, hold, latch)
	}
}

func TestCapitalSkipDoesNotLatch(t *testing.T) {
	kind, _, latch := ClassifyResearchSkip("", "below_min_notional", false)
	if latch || kind != SkipCapital {
		t.Fatalf("%s %v", kind, latch)
	}
}

func TestEligibleLatches(t *testing.T) {
	kind, _, latch := ClassifyResearchSkip("", "", true)
	if !latch || kind != "" {
		t.Fatalf("%s %v", kind, latch)
	}
}

func TestSkipSetExpires(t *testing.T) {
	p := Default()
	p.RememberSkip("BTC", "no_side", SkipStoodDown, time.Minute)
	if _, ok := p.SkipSet(time.Now().Unix())["BTC"]; !ok {
		t.Fatal("missing")
	}
	p.Skips[0].UntilUnix = time.Now().Unix() - 1
	p.PruneSkips(time.Now().Unix())
	if len(p.SkipSet(time.Now().Unix())) != 0 {
		t.Fatal(p.Skips)
	}
}

func TestSearchNoteNamesNext(t *testing.T) {
	got := SearchNote("BTC", SkipStoodDown, "ETH")
	if !strings.Contains(got, "BTC") || !strings.Contains(got, "ETH") || !strings.Contains(got, "next eligible") {
		t.Fatal(got)
	}
}
