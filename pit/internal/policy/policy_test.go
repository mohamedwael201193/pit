package policy

import "testing"

func TestDefaultAllowsTinyETH(t *testing.T) {
	p := Default()
	err := Check(p, Context{
		RequestedUSD: 10, RequestedLev: 1, Coin: "ETH", MarketType: "perp",
		Venue: "hyperliquid", SessionAlive: true, NowUnix: 1_700_000_000,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPolicyDenies(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Policy, *Context)
	}{
		{"kill", func(p *Policy, c *Context) { p.KillSwitch = true }},
		{"clip", func(_ *Policy, c *Context) { c.RequestedUSD = 11 }},
		{"lev", func(_ *Policy, c *Context) { c.RequestedLev = 50 }},
		{"coin", func(_ *Policy, c *Context) { c.Coin = "XYZ" }},
		{"venue", func(_ *Policy, c *Context) { c.Venue = "other" }},
		{"type", func(_ *Policy, c *Context) { c.MarketType = "spot" }},
		{"session", func(_ *Policy, c *Context) { c.SessionAlive = false }},
		{"loss", func(p *Policy, c *Context) { p.DailyLossUSD = 10; c.RealizedPnLUSD = -11 }},
		{"unc", func(p *Policy, c *Context) { p.MaxUncertainty = 0.2; c.Uncertainty = 0.9 }},
		{"slip", func(p *Policy, c *Context) { p.MaxSlippageBps = 10; c.SlippageBps = 80 }},
		{"liq", func(p *Policy, c *Context) { p.MinLiquidityUSD = 100; c.ImpactUSD = 1 }},
		{"cd", func(p *Policy, c *Context) { p.CooldownSeconds = 60; c.LastFillUnix = c.NowUnix - 1 }},
		{"cal", func(p *Policy, c *Context) { p.MinSkillCalibration = 0.8; c.SkillCalib = 0.1 }},
	}
	for _, tc := range cases {
		pp := Default()
		c := Context{
			RequestedUSD: 10, RequestedLev: 1, Coin: "ETH", MarketType: "perp",
			Venue: "hyperliquid", SessionAlive: true, NowUnix: 1_700_000_000,
		}
		tc.mut(&pp, &c)
		if err := Check(pp, c); err == nil {
			t.Fatalf("%s expected deny", tc.name)
		}
	}
}

func TestHashStable(t *testing.T) {
	a, err := Default().Hash()
	if err != nil || a == "" {
		t.Fatal(err)
	}
	b, _ := Default().Hash()
	if a != b {
		t.Fatal("hash drift")
	}
}
