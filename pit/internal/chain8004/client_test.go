package chain8004

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestReporterNotOwner(t *testing.T) {
	o, _ := identity.NormalizeAddress("0x1111111111111111111111111111111111111111")
	r, _ := identity.NormalizeAddress("0x2222222222222222222222222222222222222222")
	if _, err := Prepare(config.Mainnet, o, o); err == nil {
		t.Fatal("same")
	}
	reg, err := Prepare(config.Mainnet, o, r)
	if err != nil {
		t.Fatal(err)
	}
	if err := FeedbackAllowed(reg.Owner, reg.Reporter, o); err == nil {
		t.Fatal("self feedback")
	}
	if err := SameChainIDs(config.Mainnet, config.Testnet, "1", "1"); err == nil {
		t.Fatal("portable")
	}
}
