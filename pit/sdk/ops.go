package sdk

import (
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (c Client) OpportunityCopy(n int) string {
	return watch.Attention(n)
}

func (c Client) LawCards() []policy.Card {
	return policy.Cards(policy.Default())
}

func (c Client) WatchMayTrade() bool {
	return false
}

func (c Client) SwapOrLPAvailable() bool {
	return false
}
