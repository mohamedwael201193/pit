package watch

import (
	"sort"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/feasibility"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/skills"
	"github.com/mohamedwael201193/pit/internal/venue"
)

type PublicCoin struct {
	Coin              string           `json:"coin"`
	Venue             string           `json:"venue"`
	Reason            string           `json:"reason"`
	Why               string           `json:"why"`
	Trend             string           `json:"trend,omitempty"`
	Rank              int              `json:"rank"`
	RankGroup         int              `json:"rankGroup"`
	Freshness         string           `json:"freshness"`
	Mark              float64          `json:"mark"`
	Oracle            float64          `json:"oracle,omitempty"`
	Funding           float64          `json:"funding,omitempty"`
	OpenInterest      float64          `json:"openInterest,omitempty"`
	Volume            float64          `json:"volume,omitempty"`
	SzDecimals        int              `json:"szDecimals,omitempty"`
	Timestamp         string           `json:"timestamp"`
	Provenance        string           `json:"provenance"`
	Source            string           `json:"source"`
	Network           string           `json:"network"`
	Eligible          bool             `json:"eligible"`
	PolicyFit         string           `json:"policyFit"`
	ResearchEligible  bool             `json:"researchEligible"`
	PolicyEligible    bool             `json:"policyEligible"`
	ExecutionFeasible bool             `json:"executionFeasible"`
	PreviewReady      bool             `json:"previewReady"`
	Layer             string           `json:"layer,omitempty"`
	RiskFlags         []string         `json:"riskFlags,omitempty"`
	Block             string           `json:"block,omitempty"`
	ExecGate          string           `json:"execGate,omitempty"`
	ExecWhy           string           `json:"execWhy,omitempty"`
	MinNotional       float64          `json:"minNotional,omitempty"`
	RequiredMargin    float64          `json:"requiredMargin,omitempty"`
	AvailableMargin   float64          `json:"availableMargin,omitempty"`
	PolicyClip        float64          `json:"policyClip,omitempty"`
	PolicyGap         float64          `json:"policyGap,omitempty"`
	HostNotional      float64          `json:"hostNotional,omitempty"`
	HostSz            float64          `json:"hostSz,omitempty"`
	EstimatedSlippage string           `json:"estimatedSlippage,omitempty"`
	WhyExecutable     string           `json:"whyExecutable,omitempty"`
	ExpectedEdge      string           `json:"expectedEdge,omitempty"`
	Invalidation      string           `json:"invalidation,omitempty"`
	WhyRanked         string           `json:"whyRanked,omitempty"`
	Skills            []skills.Finding `json:"skills,omitempty"`
	SkillIDs          []string         `json:"skillIds,omitempty"`
}

type PublicView struct {
	OK            bool           `json:"ok"`
	Sign          bool           `json:"sign"`
	Trade         bool           `json:"trade"`
	Count         int            `json:"count"`
	Scanned       int            `json:"scanned"`
	Copy          string         `json:"copy"`
	Coins         []PublicCoin   `json:"coins"`
	Best          *PublicCoin    `json:"best,omitempty"`
	BestWhy       string         `json:"bestWhy,omitempty"`
	ExecGate      string         `json:"execGate,omitempty"`
	ExecWhy       string         `json:"execWhy,omitempty"`
	PreviewReadyN int            `json:"previewReady,omitempty"`
	ExecFeasibleN int            `json:"executionFeasible,omitempty"`
	BuyingPower   float64        `json:"buyingPower,omitempty"`
	PowerSource   string         `json:"powerSource,omitempty"`
	SpotUSDC      float64        `json:"spotUsdc,omitempty"`
	PerpEquity    float64        `json:"perpEquity,omitempty"`
	Withdrawable  float64        `json:"withdrawable,omitempty"`
	CapitalNote   string         `json:"capitalNote,omitempty"`
	Routes        []CapitalRoute `json:"routes,omitempty"`
	Source        string         `json:"source"`
	Network       string         `json:"network"`
}

func toPublic(c Candidate, net, now string) PublicCoin {
	fit := "BLOCKED"
	if c.Eligible {
		fit = "PASS"
	}
	facts := skills.Apply(c.Book, nil)
	return PublicCoin{
		Coin:              c.Coin,
		Venue:             "hyperliquid",
		Reason:            c.Reason,
		Why:               WhyHuman(c),
		Trend:             Trend(c.Book),
		Rank:              Rank(c),
		Freshness:         "live",
		Mark:              c.Book.MarkPx,
		Oracle:            c.Book.OraclePx,
		Funding:           c.Book.Funding,
		OpenInterest:      c.Book.OpenInterest,
		Volume:            c.Book.DayNtlVlm,
		SzDecimals:        c.Book.SzDecimals,
		Timestamp:         now,
		Provenance:        "hyperliquid.info metaAndAssetCtxs",
		Source:            "hyperliquid",
		Network:           net,
		Eligible:          c.Eligible,
		PolicyFit:         fit,
		ResearchEligible:  c.Eligible,
		PolicyEligible:    c.Eligible,
		ExecutionFeasible: false,
		MinNotional:       venue.PerpMinNotionalUSD(c.Book.MarkPx, c.Book.SzDecimals),
		RiskFlags:         c.Risk,
		Block:             c.Block,
		Skills:            facts,
		SkillIDs:          skills.IDs(facts),
	}
}

func Public(cands []Candidate, net string) PublicView {
	if cands == nil {
		cands = []Candidate{}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	coins := make([]PublicCoin, 0, len(cands))
	eligible := 0
	var best *PublicCoin
	for _, c := range cands {
		row := toPublic(c, net, now)
		coins = append(coins, row)
		if c.Eligible {
			eligible++
			if best == nil {
				copyRow := row
				best = &copyRow
			}
		}
	}
	bestWhy := ""
	if best != nil {
		bestWhy = "Highest host rank among policy-eligible live Hyperliquid books. Rank uses mark/oracle gap, funding, and open interest already on the venue. It is not a model score."
	}
	return PublicView{
		OK:      true,
		Sign:    false,
		Trade:   false,
		Count:   eligible,
		Scanned: len(coins),
		Copy:    Attention(eligible),
		Coins:   coins,
		Best:    best,
		BestWhy: bestWhy,
		Source:  "hyperliquid",
		Network: net,
	}
}

func EmptyPublic(net string) PublicView {
	v := Public(nil, net)
	v.Copy = Attention(0)
	return v
}

func PolicyForWatch() policy.Policy {
	return policy.Default()
}

func bookOf(c PublicCoin) hl.BookSnapshot {
	return hl.BookSnapshot{
		Coin:         c.Coin,
		MarkPx:       c.Mark,
		OraclePx:     c.Oracle,
		Funding:      c.Funding,
		OpenInterest: c.OpenInterest,
		DayNtlVlm:    c.Volume,
		SzDecimals:   c.SzDecimals,
	}
}

func attachFit(row PublicCoin, f feasibility.Fit) PublicCoin {
	row.ResearchEligible = f.ResearchEligible
	row.PolicyEligible = f.PolicyEligible
	row.ExecutionFeasible = f.ExecutionFeasible
	row.PreviewReady = f.PreviewReady
	row.Layer = f.Layer
	row.ExecGate = f.Gate
	row.ExecWhy = f.Why
	row.WhyExecutable = f.WhyExecutable
	row.MinNotional = f.MinNotionalUSD
	row.RequiredMargin = f.RequiredMargin
	row.AvailableMargin = f.AvailableMargin
	row.PolicyClip = f.PolicyClip
	if f.MinNotionalUSD > 0 && f.PolicyClip > 0 {
		row.PolicyGap = f.MinNotionalUSD - f.PolicyClip
	}
	row.HostNotional = f.HostNotional
	row.HostSz = f.HostSz
	row.EstimatedSlippage = f.EstimatedSlippage
	row.ExpectedEdge = f.ExpectedEdge
	row.Invalidation = f.Invalidation
	row.RankGroup = feasibility.RankGroup(f)
	if f.WhyExecutable != "" && f.ExecutionFeasible {
		row.WhyRanked = "Ranked because this account can actually size it, not because the signal is the largest in the universe."
	} else if f.PolicyEligible {
		row.WhyRanked = "Policy-eligible. Not execution-feasible for this account right now."
	} else {
		row.WhyRanked = "Live book only. Not policy-eligible."
	}
	return row
}

func ApplyCapital(view PublicView, acct feasibility.Account, p policy.Policy, sessionAlive, pinned bool) PublicView {
	previewN, execN := 0, 0
	for i, row := range view.Coins {
		f := feasibility.FitBook(bookOf(row), p, acct, sessionAlive, pinned)
		view.Coins[i] = attachFit(row, f)
		if view.Coins[i].PreviewReady {
			previewN++
		}
		if view.Coins[i].ExecutionFeasible {
			execN++
		}
	}
	sort.SliceStable(view.Coins, func(i, j int) bool {
		return execLess(view.Coins[i], view.Coins[j])
	})
	view.PreviewReadyN = previewN
	view.ExecFeasibleN = execN
	view.BuyingPower = acct.BuyingPower
	view.PowerSource = acct.PowerSource
	view.SpotUSDC = acct.SpotUSDC
	view.PerpEquity = acct.PerpEquity
	view.Withdrawable = acct.Withdrawable
	view.CapitalNote = acct.Note
	block, why := policy.ExecWhy(acct.OpenPositions, acct.BuyingPower, p)
	if block == "" && acct.BuyingPower+1e-9 < feasibility.MinNotionalUSD {
		block = "insufficient_margin"
		why = acct.Note
	}
	view.ExecGate = block
	view.ExecWhy = why
	if view.ExecGate == "" && execN == 0 {
		for i := range view.Coins {
			if view.Coins[i].ExecGate == "policy_clip_tight" {
				view.ExecGate = "policy_clip_tight"
				view.ExecWhy = view.Coins[i].ExecWhy
				break
			}
		}
	}
	var best *PublicCoin
	for i := range view.Coins {
		row := view.Coins[i]
		if row.PreviewReady || row.ExecutionFeasible {
			copyRow := row
			best = &copyRow
			break
		}
	}
	if best == nil {
		for i := range view.Coins {
			if view.Coins[i].PolicyEligible {
				copyRow := view.Coins[i]
				best = &copyRow
				break
			}
		}
	}
	view.Best = best
	if best != nil && (best.PreviewReady || best.ExecutionFeasible) {
		view.BestWhy = "Highest host rank among books this account can actually size under pinned policy, venue minimum, and available margin. BTC is not auto-first. Side is not decided here."
	} else if best != nil {
		extra := view.ExecWhy
		if extra == "" {
			extra = why
		}
		view.BestWhy = "Highest host rank among policy-eligible live books. None are execution-feasible for this account right now. " + extra
	}
	view.Routes = DecideRoutes(view)
	return view
}

func assetIndex(coin string) int {
	u := strings.ToUpper(strings.TrimSpace(coin))
	for i, a := range policy.HostAssets {
		if a == u {
			return i
		}
	}
	return len(policy.HostAssets)
}

func execLess(a, b PublicCoin) bool {
	if a.RankGroup != b.RankGroup {
		return a.RankGroup > b.RankGroup
	}
	if a.MinNotional != b.MinNotional {
		if a.MinNotional <= 0 {
			return false
		}
		if b.MinNotional <= 0 {
			return true
		}
		return a.MinNotional < b.MinNotional
	}
	ai, bi := assetIndex(a.Coin), assetIndex(b.Coin)
	if ai != bi {
		return ai < bi
	}
	return a.Rank > b.Rank
}

func BestExecutable(cands []Candidate, acct feasibility.Account, p policy.Policy, sessionAlive, pinned bool) (Candidate, bool) {
	return BestExecutableExcept(cands, acct, p, sessionAlive, pinned, nil)
}

func BestExecutableExcept(cands []Candidate, acct feasibility.Account, p policy.Policy, sessionAlive, pinned bool, skip map[string]string) (Candidate, bool) {
	type scored struct {
		c   Candidate
		g   int
		r   int
		min float64
		idx int
	}
	rows := make([]scored, 0, len(cands))
	for _, c := range cands {
		if skip != nil {
			if _, hit := skip[strings.ToUpper(c.Coin)]; hit {
				continue
			}
		}
		f := feasibility.FitBook(c.Book, p, acct, sessionAlive, pinned)
		rows = append(rows, scored{c: c, g: feasibility.RankGroup(f), r: Rank(c), min: f.MinNotionalUSD, idx: assetIndex(c.Coin)})
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].g != rows[j].g {
			return rows[i].g > rows[j].g
		}
		if rows[i].min != rows[j].min {
			if rows[i].min <= 0 {
				return false
			}
			if rows[j].min <= 0 {
				return true
			}
			return rows[i].min < rows[j].min
		}
		if rows[i].idx != rows[j].idx {
			return rows[i].idx < rows[j].idx
		}
		return rows[i].r > rows[j].r
	})
	for _, row := range rows {
		if row.g >= 3 {
			return row.c, true
		}
	}
	return Candidate{}, false
}

func NextCandidate(cands []Candidate, acct feasibility.Account, p policy.Policy, sessionAlive, pinned bool, skip map[string]string) (Candidate, bool, bool) {
	if best, ok := BestExecutableExcept(cands, acct, p, sessionAlive, pinned, skip); ok {
		return best, true, true
	}
	if best, ok := BestExcept(cands, skip); ok {
		return best, false, true
	}
	return Candidate{}, false, false
}
