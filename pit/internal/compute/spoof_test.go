package compute

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestMatchProvider(t *testing.T) {
	m := MainnetChat()
	if err := MatchProvider(config.Mainnet, m.Provider, m.Model, m.TeeSigner); err != nil {
		t.Fatal(err)
	}
	if err := MatchProvider(config.Mainnet, "0xdead", m.Model, m.TeeSigner); err == nil {
		t.Fatal("provider")
	}
	if err := MatchProvider(config.Mainnet, m.Provider, "other", m.TeeSigner); err == nil {
		t.Fatal("model")
	}
	if err := MatchProvider(config.Testnet, m.Provider, m.Model, m.TeeSigner); err == nil {
		t.Fatal("copied sku")
	}
}
