package feasibility

import (
	"fmt"
	"math"
	"strings"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

const MinNotionalUSD = 10.0

const (
	LayerResearch = "research-eligible"
	LayerPolicy   = "policy-eligible"
	LayerExec     = "execution-feasible"
	LayerPreview  = "preview-ready"
	LayerBlocked  = "execution-blocked"
)

type Account struct {
	BuyingPower   float64
	SpotUSDC      float64
	SpotFree      float64
	PerpEquity    float64
	Withdrawable  float64
	PowerSource   string
	Note          string
	OpenPositions int
	OpenOrders    int
	Abstraction   string
	MarginUsed    float64
}

func FromCapital(c hl.Capital) Account {
	return Account{
		BuyingPower:   c.BuyingPower,
		SpotUSDC:      c.SpotUSDC,
		SpotFree:      c.SpotFree,
		PerpEquity:    c.PerpEquity,
		Withdrawable:  c.Withdrawable,
		PowerSource:   c.PowerSource,
		Note:          c.Note,
		OpenPositions: c.OpenPositions,
		OpenOrders:    c.OpenOrders,
		Abstraction:   c.Abstraction,
		MarginUsed:    c.MarginUsed,
	}
}

type Fit struct {
	Coin              string
	ResearchEligible  bool
	PolicyEligible    bool
	ExecutionFeasible bool
	PreviewReady      bool
	Layer             string
	Gate              string
	Why               string
	WhyExecutable     string
	MinNotionalUSD    float64
	RequiredMargin    float64
	AvailableMargin   float64
	PolicyClip        float64
	HostNotional      float64
	HostSz            float64
	SzDecimals        int
	EstimatedSlippage string
	ExpectedEdge      string
	Invalidation      string
}

func policyScan(p policy.Policy, b hl.BookSnapshot) error {
	return policy.Check(p, policy.Context{
		RequestedUSD: p.MaxClipUSD,
		RequestedLev: 1,
		Coin:         b.Coin,
		MarketType:   "perp",
		Venue:        "hyperliquid",
		SessionAlive: true,
		NowUnix:      1,
		ImpactUSD:    b.OpenInterest,
	})
}

func FitBook(b hl.BookSnapshot, p policy.Policy, acct Account, sessionAlive, pinned bool) Fit {
	p = policy.Clamp(p)
	f := Fit{
		Coin:              strings.ToUpper(strings.TrimSpace(b.Coin)),
		MinNotionalUSD:    MinNotionalUSD,
		AvailableMargin:   acct.BuyingPower,
		PolicyClip:        p.MaxClipUSD,
		SzDecimals:        b.SzDecimals,
		EstimatedSlippage: fmt.Sprintf("host ceiling %d bps (live L2 is not sampled for the full universe)", p.MaxSlippageBps),
		ExpectedEdge:      "host rank of venue facts; not a model score",
		Invalidation:      invalidation(b),
	}
	if b.MarkPx <= 0 || f.Coin == "" {
		f.Layer = LayerBlocked
		f.Gate = "no_mark"
		f.Why = "No live Hyperliquid mark. PIT will not invent a book."
		return f
	}
	f.ResearchEligible = true
	if err := policyScan(p, b); err != nil {
		f.Layer = LayerResearch
		f.Gate = err.Error()
		f.Why = fmt.Sprintf("%s is live on Hyperliquid but policy blocked it (%s). Research is refused for this coin.", f.Coin, err.Error())
		f.WhyExecutable = "Not executable: outside pinned policy."
		return f
	}
	f.PolicyEligible = true

	block, why := policy.ExecWhy(acct.OpenPositions, acct.BuyingPower, p)
	if block == "" && acct.BuyingPower+1e-9 < MinNotionalUSD {
		block = "insufficient_margin"
		why = execCapitalWhy(acct)
	}
	if block != "" {
		if block == "insufficient_margin" && strings.TrimSpace(acct.Note) != "" && !strings.Contains(why, acct.Note) {
			why = strings.TrimSpace(acct.Note) + " " + why
		}
		f.Layer = LayerBlocked
		f.Gate = block
		f.Why = why
		f.WhyExecutable = why
		return f
	}

	usd := p.MaxClipUSD
	if acct.BuyingPower > 0 && acct.BuyingPower < usd {
		usd = acct.BuyingPower
	}
	sized, err := engine.SizeOrder(engine.SizerInput{
		MarkPx: b.MarkPx, SzDecimals: b.SzDecimals, MaxClipUSD: p.MaxClipUSD, RequestedUSD: usd,
		Side: "buy", Coin: f.Coin, AllowedCoins: p.AllowedAssets, MaxLeverage: 1, RequestedLev: 1,
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err != nil {
		f.Layer = LayerBlocked
		f.Gate = err.Error()
		f.Why = fmt.Sprintf("%s cannot be sized under this account and policy (%s). PIT will not invent size.", f.Coin, err.Error())
		f.WhyExecutable = f.Why
		return f
	}
	f.HostNotional = sized.NotionalUSD
	f.HostSz = sized.Sz
	f.RequiredMargin = sized.NotionalUSD
	f.ExecutionFeasible = true
	f.Layer = LayerExec
	f.WhyExecutable = fmt.Sprintf(
		"Executable for this user: host-sized notional $%.2f uses $%.2f available (%s), above the $%.0f Hyperliquid minimum, inside the $%.0f policy clip. Side is still not decided here.",
		sized.NotionalUSD, acct.BuyingPower, acct.PowerSource, MinNotionalUSD, p.MaxClipUSD,
	)
	if !pinned {
		f.Gate = "policy_unpinned"
		f.Why = "Policy is not pinned on this computer. Research can still run. AUTHORIZE stays fail-closed."
		f.WhyExecutable = f.WhyExecutable + " Preview is not ready until you pin policy."
		return f
	}
	if !sessionAlive {
		f.Gate = "session_expired"
		f.Why = "No live order/cancel session. Research can still run. AUTHORIZE stays fail-closed."
		f.WhyExecutable = f.WhyExecutable + " Preview is not ready until the Hyperliquid session is alive."
		return f
	}
	f.PreviewReady = true
	f.Layer = LayerPreview
	f.Why = f.WhyExecutable
	return f
}

func execCapitalWhy(acct Account) string {
	short := MinNotionalUSD - acct.BuyingPower
	if short < 0 {
		short = 0
	}
	line := fmt.Sprintf("Available venue margin is $%.2f — $%.2f short of the $%.0f Hyperliquid minimum. PIT will not invent size.", acct.BuyingPower, short, MinNotionalUSD)
	if acct.Note != "" && acct.BuyingPower+1e-9 < MinNotionalUSD {
		return acct.Note + " " + line
	}
	return line
}

func invalidation(b hl.BookSnapshot) string {
	if b.OraclePx > 0 {
		return "Mark/oracle gap closing, funding flipping, or open interest collapsing would remove this host rank."
	}
	return "A missing live mark or a policy fail would remove this book."
}

func RankGroup(f Fit) int {
	switch {
	case f.PreviewReady:
		return 4
	case f.ExecutionFeasible:
		return 3
	case f.PolicyEligible:
		return 2
	case f.ResearchEligible:
		return 1
	default:
		return 0
	}
}

func RoundUSD(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return math.Round(v*10000) / 10000
}
