package experience

import (
	"strings"
	"testing"
)

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
	if !strings.Contains(got, "Resting is not a fill") {
		t.Fatal(got)
	}
}

func TestRestingIsNotFilled(t *testing.T) {
	rows := make([]Case, MinSamples)
	for i := range rows {
		rows[i] = Case{Coin: "ETH", Decision: DecisionResting}
	}
	got := WhyThisSetup("ETH", rows)
	if strings.Contains(got, "5 confirmed fills") {
		t.Fatal(got)
	}
	if !strings.Contains(got, "5 resting OIDs") {
		t.Fatal(got)
	}
}
