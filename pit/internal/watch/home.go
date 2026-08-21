package watch

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

type HomeCard struct {
	Coin    string
	Eligible bool
	Reason  string
}

func Home(cands []Candidate, p policy.Policy) []HomeCard {
	_ = p
	out := make([]HomeCard, 0, len(cands))
	for _, c := range cands {
		out = append(out, HomeCard{Coin: c.Coin, Eligible: true, Reason: c.Reason})
	}
	return out
}

func Attention(n int) string {
	if n <= 0 {
		return "No opportunities match your policy."
	}
	return fmt.Sprintf("%d opportunities match your policy.", n)
}

func BlockedCard(coin, why string) HomeCard {
	return HomeCard{Coin: coin, Eligible: false, Reason: why}
}
