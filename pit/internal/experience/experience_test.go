package experience

import "testing"

func TestClassifyStandDownIsNoTrade(t *testing.T) {
	if Classify(false, "no_side", "") != DecisionNoTrade {
		t.Fatal("stand-down")
	}
	if Classify(true, "", "") != DecisionReady {
		t.Fatal("ready")
	}
	if Classify(false, "insufficient_margin", "") != DecisionCapital {
		t.Fatal("capital")
	}
}

func TestWhyThisSetupNeedsSamples(t *testing.T) {
	got := WhyThisSetup("BTC", []Case{{Coin: "BTC", Decision: DecisionNoTrade}})
	if got[:15] != "NOT ENOUGH DATA" {
		t.Fatal(got)
	}
	rows := make([]Case, MinSamples)
	for i := range rows {
		rows[i] = Case{Coin: "ETH", Decision: DecisionNoTrade}
	}
	got = WhyThisSetup("ETH", rows)
	if got[:15] == "NOT ENOUGH DATA" {
		t.Fatal(got)
	}
}
