package deskcmd

import (
	"strings"
	"testing"
)

func TestWhyDidntYouTrade(t *testing.T) {
	r := Parse("Why didn't you trade?")
	if r.Execute || r.Tool != "watch.why_not" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Why is nothing executable?")
	if r.Execute || r.Tool != "watch.why_not" {
		t.Fatalf("%+v", r)
	}
}

func TestRefuseExecute(t *testing.T) {
	for _, q := range []string{"buy ETH", "I authorize it", "just do it", "flatten ETH", "close my position"} {
		r := Parse(q)
		if r.Execute || r.StartResearch {
			t.Fatalf("%s %+v", q, r)
		}
		if r.Navigate != "preview" {
			t.Fatalf("%s navigate %s", q, r.Navigate)
		}
	}
	r := Parse("trade now")
	if r.Execute || r.StartResearch || r.Tool != "preview.show" || r.Navigate != "" {
		t.Fatalf("trade now %+v", r)
	}
}

func TestResearchETH(t *testing.T) {
	r := Parse("Research ETH privately.")
	if !r.StartResearch || r.Coin != "ETH" || r.Execute {
		t.Fatalf("%+v", r)
	}
}

func TestOpenOfficialPages(t *testing.T) {
	r := Parse("Open 0G Private Compute")
	if r.OpenURL != "https://pc.0g.ai/sdk/dashboard/funds" {
		t.Fatal(r.OpenURL)
	}
	r = Parse("Open Hyperliquid")
	if r.OpenURL != "https://app.hyperliquid.xyz" {
		t.Fatal(r.OpenURL)
	}
}

func TestWhyAndForget(t *testing.T) {
	r := Parse("Why is ETH interesting?")
	if r.Execute || r.StartResearch || r.Tool != "watch.get" {
		t.Fatalf("%+v", r)
	}
	r = Parse("forget everything")
	if r.Tool != "memory.forget" || r.Execute {
		t.Fatal(r)
	}
	r = Parse("Show me the evidence")
	if r.Tool != "activity.proof" || r.Execute || r.Navigate == "markets" {
		t.Fatal(r)
	}
}

func TestOperatorIntents(t *testing.T) {
	r := Parse("Is Hyperliquid ready?")
	if r.Tool != "session.status" || r.Execute || r.Navigate != "security" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show the current preview.")
	if r.Tool != "preview.show" || r.Navigate != "research" || r.Execute {
		t.Fatalf("%+v", r)
	}
	r = Parse("Explain my policy.")
	if r.Tool != "policy.get" || r.Navigate != "security" || r.Execute {
		t.Fatalf("%+v", r)
	}
	r = Parse("Open Hyperliquid API")
	if r.OpenURL != "https://app.hyperliquid.xyz/API" || r.Execute {
		t.Fatalf("%+v", r)
	}
}

func TestHappeningAndPositions(t *testing.T) {
	r := Parse("What is happening?")
	if r.Execute || r.Tool != "status" {
		t.Fatalf("%+v", r)
	}
	r = Parse("What is happening with BTC?")
	if r.Execute || r.Tool != "watch.get" || r.Coin != "BTC" || r.Navigate == "markets" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Why is ETH moving?")
	if r.Execute || r.Tool != "watch.get" || r.Coin != "ETH" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Compare BTC and ETH.")
	if r.Execute || r.Tool != "watch.compare" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show me the last research.")
	if r.Execute || r.Tool != "research.result" || r.Navigate != "research" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show me my current positions")
	if r.Navigate != "portfolio" || r.Execute {
		t.Fatalf("%+v", r)
	}
}

func TestETHSetupAndCannotExecute(t *testing.T) {
	r := Parse("What is the ETH setup?")
	if r.Execute || r.StartResearch || r.Tool != "watch.get" || r.Coin != "ETH" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Why can't PIT execute this?")
	if r.Execute || r.StartResearch || r.Tool != "refuse_execute" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Open 0G")
	if r.OpenURL != "https://pc.0g.ai/sdk/dashboard/funds" {
		t.Fatal(r.OpenURL)
	}
	r = Parse("Show me today's opportunities")
	if r.Execute || r.StartResearch || r.Tool != "watch.get" {
		t.Fatalf("%+v", r)
	}
}

func TestGreetingIsNotCannedHelp(t *testing.T) {
	r := Parse("HI")
	if r.Execute || r.Tool != "greet" {
		t.Fatalf("%+v", r)
	}
	if strings.Contains(r.Reply, "I can research a policy market") {
		t.Fatal(r.Reply)
	}
}

func TestTypoResearchStartsSealedPass(t *testing.T) {
	r := Parse("HI WHAT PRICE OF ETH NOW AND GIVE ME REASEARCH AND DO TRADE")
	if r.Execute || !r.StartResearch || r.Coin != "ETH" {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Reply, "cannot AUTHORIZE") {
		t.Fatal(r.Reply)
	}
}

func TestStrongestOpportunityStartsResearch(t *testing.T) {
	r := Parse("Find the strongest opportunity.")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" || r.Navigate != "" {
		t.Fatalf("%+v", r)
	}
}

func TestFindBestTradeStaysOnChat(t *testing.T) {
	r := Parse("Find the best trade available right now.")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" || r.Navigate == "markets" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Find me the best long")
	if r.Execute || !r.StartResearch || r.Hypothesis != "long" {
		t.Fatalf("%+v", r)
	}
	r = Parse("What can I trade with my current capital?")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" {
		t.Fatalf("%+v", r)
	}
	r = Parse("What happened with the strategy?")
	if r.Coin == "ETH" || r.Coin == "BTC" || r.Coin == "SOL" {
		t.Fatalf("ticker substring %s", r.Coin)
	}
	r = Parse("Tell me what is happening with BTC")
	if r.Coin != "BTC" || r.StartResearch {
		t.Fatalf("%+v", r)
	}
}

func TestAutonomyIntents(t *testing.T) {
	r := Parse("Find me the best opportunity right now.")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" || r.Navigate == "markets" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Scan everything allowed by my policy.")
	if r.Execute || r.Tool != "watch.scan" || r.Navigate == "markets" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Research the best setup.")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" {
		t.Fatalf("%+v", r)
	}
	r = Parse("research in market and choose best opportunity now in market and give me research on it and why should enter")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" {
		t.Fatalf("%+v", r)
	}
	r = Parse("I accept")
	if r.Execute || r.StartResearch || r.Tool != "preview.show" || r.Navigate != "" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Trade this")
	if r.Execute || r.StartResearch || r.Tool != "preview.show" || r.Navigate != "" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Reject this")
	if r.Execute || r.StartResearch || r.Tool != "watch.why_not" {
		t.Fatalf("%+v", r)
	}
	r = Parse("AUTHORIZE")
	if r.Execute || r.StartResearch || r.Tool != "refuse_execute" {
		t.Fatalf("%+v", r)
	}
	r = Parse("show all txs")
	if r.Execute || r.Tool != "activity.list" {
		t.Fatalf("%+v", r)
	}
	r = Parse("why should I enter it")
	if r.Execute || r.StartResearch || r.Tool != "experience.why" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Trade the strongest setup.")
	if r.Execute || !r.StartResearch || r.Tool != "research.best" || r.Navigate == "markets" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Run autonomously for 24 hours.")
	if r.Execute || r.StartResearch || r.Tool != "mission.enable_required" || r.Navigate != "automation" || r.Hours != 24 {
		t.Fatalf("%+v", r)
	}
	r = Parse("Enable guarded autonomy for 8 hours.")
	if r.Execute || r.Tool != "mission.enable_required" || r.Hours != 8 || r.Navigate != "automation" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Stop autonomous trading.")
	if r.Execute || r.Tool != "mission.stop" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Start a sleep mission.")
	if r.Execute || r.Tool != "mission.enable_required" || r.Navigate != "automation" {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Reply, "ENABLE GUARDED AUTONOMY") || !strings.Contains(r.Reply, "ARM SLEEP MISSION") {
		t.Fatalf("%s", r.Reply)
	}
	r = Parse("What happened overnight?")
	if r.Execute || r.Tool != "mission.status" || r.Navigate == "markets" {
		t.Fatalf("%+v", r)
	}
	r = Parse("How did last night's mission perform?")
	if r.Execute || r.Tool != "mission.status" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show every trade PIT made today.")
	if r.Execute || r.Tool != "activity.today" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show me the on-chain proof.")
	if r.Execute || r.Tool != "activity.proof" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show me the proof for that trade.")
	if r.Execute || r.Tool != "activity.proof" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Why is this opportunity better than the others?")
	if r.Execute || r.Tool != "watch.compare" {
		t.Fatalf("%+v", r)
	}
}

func TestDoTradeAloneRefuses(t *testing.T) {
	r := Parse("do trade")
	if r.Execute || r.StartResearch || r.Tool != "refuse_execute" {
		t.Fatalf("%+v", r)
	}
}

func TestNoShell(t *testing.T) {
	r := Parse("rm -rf /")
	if r.StartResearch || r.Execute {
		t.Fatal(r)
	}
}

func TestResearchWithoutNamedCoinDoesNotInventETH(t *testing.T) {
	r := Parse("Research privately.")
	if !r.StartResearch || r.Execute || r.Coin != "" {
		t.Fatalf("%+v", r)
	}
}

func TestResearchPumpDoesNotBecomeETH(t *testing.T) {
	r := Parse("Research PUMP privately.")
	if !r.StartResearch || r.Coin != "PUMP" || r.Execute {
		t.Fatalf("%+v", r)
	}
}

func TestHyperliquidDoesNotMatchHYPE(t *testing.T) {
	r := Parse("Open Hyperliquid")
	if r.Coin == "HYPE" {
		t.Fatal("hyperliquid matched HYPE")
	}
	if r.OpenURL != "https://app.hyperliquid.xyz" {
		t.Fatal(r.OpenURL)
	}
	r = Parse("What is happening with BTC?")
	if r.Coin != "BTC" || r.Tool != "watch.get" {
		t.Fatalf("%+v", r)
	}
}

func TestChatCannotMutatePolicy(t *testing.T) {
	for _, q := range []string{"raise clip to 1000", "set leverage 20x", "edit my policy", "increase max open positions"} {
		r := Parse(q)
		if r.Execute || r.Mutate {
			t.Fatalf("%s %+v", q, r)
		}
		if r.Tool != "policy.get" {
			t.Fatalf("%s tool %s", q, r.Tool)
		}
		if r.Navigate != "security" {
			t.Fatalf("%s nav %s", q, r.Navigate)
		}
		if !strings.Contains(strings.ToLower(r.Reply), "cannot") {
			t.Fatalf("%s reply %s", q, r.Reply)
		}
	}
}
