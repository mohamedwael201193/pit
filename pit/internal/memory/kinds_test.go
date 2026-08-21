package memory

import "testing"

func TestKindsIncludeThesisAndRegime(t *testing.T) {
	for _, k := range []string{"observation", "thesis", "forecast", "execution", "outcome", "error", "skill_performance", "calibration", "regime", "risk_event"} {
		if !ValidKind(k) {
			t.Fatal(k)
		}
	}
	if ValidKind("global") {
		t.Fatal("global")
	}
}
