package experience

import (
	"fmt"
	"strings"
)

const MinSamples = 5

const (
	DecisionReady     = "READY"
	DecisionNoTrade   = "NO_TRADE"
	DecisionFailed    = "FAILED"
	DecisionCapital   = "CAPITAL_BLOCKED"
	DecisionExecution = "EXECUTION_BLOCKED"
	DecisionFilled    = "FILLED"
)

type Case struct {
	Unix        int64  `json:"unix"`
	Coin        string `json:"coin"`
	Decision    string `json:"decision"`
	Executable  bool   `json:"executable"`
	Interesting bool   `json:"interesting"`
	PreviewHash string `json:"preview_hash,omitempty"`
	OID         string `json:"oid,omitempty"`
	Why         string `json:"why,omitempty"`
	Root        string `json:"root,omitempty"`
}

func Classify(eligible bool, deny, jobErr string) string {
	if eligible {
		return DecisionReady
	}
	d := strings.ToLower(strings.TrimSpace(deny))
	e := strings.ToLower(strings.TrimSpace(jobErr))
	switch {
	case d == "insufficient_margin" || e == "insufficient_margin":
		return DecisionCapital
	case d == "below_min_notional" || d == "below_minimum" || e == "below_min_notional":
		return DecisionExecution
	case d == "no_side" || d == "challenger_killed" || d == "risk_killed" || d == "stood_down":
		return DecisionNoTrade
	case e != "":
		return DecisionFailed
	default:
		return DecisionNoTrade
	}
}

func ForCoin(all []Case, coin string) []Case {
	want := strings.ToUpper(strings.TrimSpace(coin))
	out := make([]Case, 0)
	for _, c := range all {
		if strings.ToUpper(c.Coin) == want {
			out = append(out, c)
		}
	}
	return out
}

func WhyThisSetup(coin string, cases []Case) string {
	n := len(cases)
	c := strings.ToUpper(strings.TrimSpace(coin))
	if c == "" {
		c = "this market"
	}
	if n < MinSamples {
		return fmt.Sprintf("NOT ENOUGH DATA (%d/%d verified cases for %s). PIT will not invent that this desk learned.", n, MinSamples, c)
	}
	noTrade, ready, failed, filled := 0, 0, 0, 0
	for _, row := range cases {
		switch row.Decision {
		case DecisionReady:
			ready++
		case DecisionFailed:
			failed++
		case DecisionFilled:
			filled++
		default:
			noTrade++
		}
	}
	return fmt.Sprintf("%s has %d verified internal cases. %d ready previews, %d successful stand-downs or blocks, %d confirmed fills, %d failed research. This is workspace memory, not a guarantee.", c, n, ready, noTrade, filled, failed)
}
