package skills

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestRegistryVersionedAndHostOnly(t *testing.T) {
	reg := Registry()
	if len(reg) < 8 {
		t.Fatalf("%d", len(reg))
	}
	seen := map[string]bool{}
	for _, s := range reg {
		if s.Authority != "host" {
			t.Fatalf("%+v", s)
		}
		if s.Version == "" || s.ID == "" {
			t.Fatalf("%+v", s)
		}
		if seen[s.ID] {
			t.Fatalf("dup %s", s.ID)
		}
		seen[s.ID] = true
	}
}

func TestApplyDoesNotInventCandles(t *testing.T) {
	out := Apply(hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, SzDecimals: 4}, nil)
	found := false
	for _, f := range out {
		if f.Skill.ID == "candle_structure" {
			found = true
			if f.Present {
				t.Fatal("invented candles")
			}
		}
	}
	if !found {
		t.Fatal("missing candle skill")
	}
	ids := IDs(out)
	if len(ids) != len(out) {
		t.Fatal("provenance")
	}
}
