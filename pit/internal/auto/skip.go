package auto

import (
	"fmt"
	"strings"
	"time"
)

const (
	SkipStoodDown  = "stood_down"
	SkipFailed     = "failed"
	SkipCapital    = "capital"
	SkipResearched = "already_researched"
)

type ResearchSkip struct {
	Coin      string `json:"coin"`
	Why       string `json:"why"`
	Kind      string `json:"kind"`
	UntilUnix int64  `json:"until_unix"`
}

func (p *Prefs) PruneSkips(now int64) {
	if p == nil || len(p.Skips) == 0 {
		return
	}
	next := p.Skips[:0]
	for _, s := range p.Skips {
		if s.UntilUnix > now && strings.TrimSpace(s.Coin) != "" {
			next = append(next, s)
		}
	}
	if len(next) == 0 {
		p.Skips = nil
		return
	}
	p.Skips = next
}

func (p Prefs) SkipSet(now int64) map[string]string {
	out := map[string]string{}
	for _, s := range p.Skips {
		if s.UntilUnix > now {
			out[strings.ToUpper(strings.TrimSpace(s.Coin))] = s.Kind
			if s.Why != "" {
				out[strings.ToUpper(strings.TrimSpace(s.Coin))] = s.Why
			}
		}
	}
	return out
}

func (p *Prefs) RememberSkip(coin, why, kind string, hold time.Duration) {
	if p == nil {
		return
	}
	coin = strings.ToUpper(strings.TrimSpace(coin))
	if coin == "" {
		return
	}
	if hold <= 0 {
		hold = 10 * time.Minute
	}
	now := time.Now().Unix()
	p.PruneSkips(now)
	until := now + int64(hold.Seconds())
	for i, s := range p.Skips {
		if strings.EqualFold(s.Coin, coin) {
			p.Skips[i] = ResearchSkip{Coin: coin, Why: why, Kind: kind, UntilUnix: until}
			return
		}
	}
	p.Skips = append(p.Skips, ResearchSkip{Coin: coin, Why: why, Kind: kind, UntilUnix: until})
}

func ClassifyResearchSkip(jobErr, deny string, eligible bool) (kind string, hold time.Duration, latch bool) {
	if eligible {
		return "", 0, true
	}
	switch strings.ToLower(strings.TrimSpace(deny)) {
	case "no_side", "challenger_killed", "risk_killed":
		return SkipStoodDown, 10 * time.Minute, false
	case "insufficient_margin", "below_min_notional":
		return SkipCapital, 2 * time.Minute, false
	}
	if strings.TrimSpace(jobErr) != "" {
		return SkipFailed, 3 * time.Minute, false
	}
	if strings.TrimSpace(deny) != "" {
		return "rejected", 5 * time.Minute, false
	}
	return SkipResearched, 2 * time.Minute, false
}

func (p Prefs) LatestSkip() ResearchSkip {
	if len(p.Skips) == 0 {
		return ResearchSkip{}
	}
	best := p.Skips[0]
	for _, s := range p.Skips[1:] {
		if s.UntilUnix > best.UntilUnix {
			best = s
		}
	}
	return best
}

func SearchNote(skipped, why, next string) string {
	skipped = strings.ToUpper(strings.TrimSpace(skipped))
	next = strings.ToUpper(strings.TrimSpace(next))
	human := HumanSkip(why)
	if skipped == "" {
		return ""
	}
	if next == "" {
		return fmt.Sprintf("%s: %s. No remaining eligible market.", skipped, human)
	}
	return fmt.Sprintf("%s: %s. PIT is checking the next eligible market (%s).", skipped, human, next)
}

func HumanSkip(why string) string {
	switch strings.ToLower(strings.TrimSpace(why)) {
	case SkipStoodDown, "no_side", "challenger_killed", "risk_killed", "research_stood_down":
		return "no side survived challenge"
	case SkipCapital, "insufficient_margin", "below_min_notional":
		return "blocked by capital for this market"
	case SkipFailed:
		return "research failed"
	case SkipResearched:
		return "already researched this cycle"
	default:
		if h := HumanWhy(why); h != "" {
			return h
		}
		if why == "" {
			return "not a trade"
		}
		return why
	}
}
