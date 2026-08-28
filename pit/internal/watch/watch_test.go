package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestScanRespectsAllowlistAndEmpty(t *testing.T) {
	p := policy.Default()
	cands, err := Scan([]hl.BookSnapshot{
		{Coin: "XYZ", MarkPx: 1, OraclePx: 1, OpenInterest: 1e9},
		{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, OpenInterest: 1e9},
	}, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 || cands[0].Coin != "ETH" {
		t.Fatalf("%+v", cands)
	}
	all := Universe([]hl.BookSnapshot{
		{Coin: "XYZ", MarkPx: 1, OraclePx: 1, OpenInterest: 1e9},
		{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, OpenInterest: 1e9},
	}, p)
	if len(all) != 2 {
		t.Fatalf("universe %+v", all)
	}
	best, ok := Best(all)
	if !ok || best.Coin != "ETH" || !best.Eligible {
		t.Fatalf("best %+v %v", best, ok)
	}
	c2, _ := Scan(nil, p)
	if EmptyCopy(len(c2)) != "No opportunities match your policy." {
		t.Fatal("empty copy")
	}
}
