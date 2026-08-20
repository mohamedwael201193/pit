package workspace

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestTwoWorkspacesCannotShare(t *testing.T) {
	st := NewStore()
	a, err := identity.NormalizeAddress("0x1111111111111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	b, err := identity.NormalizeAddress("0x2222222222222222222222222222222222222222")
	if err != nil {
		t.Fatal(err)
	}
	wa, err := st.Create(a, config.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	wb, err := st.Create(b, config.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	if wa.ID == wb.ID {
		t.Fatal("same id")
	}
	if wa.Network == wb.Network {
		t.Fatal("networks mixed")
	}
	if _, err := st.Get(wb.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.Get("00000000-0000-4000-8000-000000000000"); err == nil {
		t.Fatal("guess should not found")
	}
}
