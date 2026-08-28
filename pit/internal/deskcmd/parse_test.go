package deskcmd

import "testing"

func TestRefuseExecute(t *testing.T) {
	for _, q := range []string{"buy ETH", "trade now", "I authorize it", "just do it", "flatten ETH", "close my position"} {
		r := Parse(q)
		if r.Execute || r.StartResearch {
			t.Fatalf("%s %+v", q, r)
		}
		if r.Navigate != "preview" {
			t.Fatalf("%s navigate %s", q, r.Navigate)
		}
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
	if r.Navigate != "research" || r.Execute {
		t.Fatal(r)
	}
}

func TestHappeningAndPositions(t *testing.T) {
	r := Parse("What is happening?")
	if r.Execute || r.Navigate != "home" {
		t.Fatalf("%+v", r)
	}
	r = Parse("Show me my current positions")
	if r.Navigate != "positions" || r.Execute {
		t.Fatalf("%+v", r)
	}
}

func TestNoShell(t *testing.T) {
	r := Parse("rm -rf /")
	if r.StartResearch || r.Execute {
		t.Fatal(r)
	}
}
