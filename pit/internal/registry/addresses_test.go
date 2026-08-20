package registry

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestChainSpecific(t *testing.T) {
	m := For(config.Mainnet)
	g := For(config.Testnet)
	if m.Identity8004 == g.Identity8004 {
		t.Fatal("8004 ids are not portable")
	}
	if m.Serving == g.Serving {
		t.Fatal("serving")
	}
}
