package watch

import "testing"

func TestCompactHuman(t *testing.T) {
	if Compact(4.549873164451658e+09) != "4.55B" {
		t.Fatal(Compact(4.549873164451658e+09))
	}
	if CompactUSD(77823) != "$77823" && Compact(77823) != "77823" {
		t.Fatal(Compact(77823))
	}
	if FundingPct(0.0001) != "0.0100%" {
		t.Fatal(FundingPct(0.0001))
	}
}
