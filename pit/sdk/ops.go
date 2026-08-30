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

func (c Client) CanArmMission() bool {
	return false
}

func (c Client) MissionStatus() map[string]any {
	return map[string]any{
		"arm": false, "authorize": false, "sign": false, "trade": false,
		"copy": "SDK is read-only. Arm a Sleep Mission on PIT Desktop or CLI TTY.",
	}
}

func (c Client) SwapOrLPAvailable() bool {
	return false
}
