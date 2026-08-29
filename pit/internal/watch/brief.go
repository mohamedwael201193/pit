package watch

import "fmt"

func WhyInteresting(c Candidate) string {
	return WhyHuman(c)
}

func WhyPolicy(c Candidate) string {
	return WhyPolicyFrom(c.Eligible, c.Block)
}

func WhyPolicyFrom(eligible bool, block string) string {
	if eligible {
		return "Asset, venue, clip, leverage, and market type are inside the pinned host policy."
	}
	if block != "" {
		return "Policy blocked this book: " + block
	}
	return "Outside the pinned policy universe."
}

func WhatInvalidates(c Candidate) string {
	return WhatInvalidatesReason(c.Reason)
}

func WhatInvalidatesReason(reason string) string {
	switch reason {
	case "mark_below_oracle":
		return "The mark/oracle gap closes, the book stops printing, or policy eligibility is lost."
	case "funding":
		return "Funding prints zero, liquidity thins, or policy eligibility is lost."
	case "blocked":
		return "Already invalid for a trade. Research will not produce an eligible preview."
	default:
		return "The live book goes stale, policy eligibility is lost, or a host gate trips."
	}
}

func ResearchWillTest(c Candidate) string {
	return ResearchWillTestFrom(c.Coin, c.Eligible)
}

func ResearchWillTestFrom(coin string, eligible bool) string {
	if !eligible {
		return fmt.Sprintf("%s is not eligible. Research is not started from this row.", coin)
	}
	return "Researcher, challenger, and risk must all verify. Host then decides. A stand-down is a successful research outcome."
}
