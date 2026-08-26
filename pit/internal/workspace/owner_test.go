package workspace

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestAssertOwnerRejectsOtherWallet(t *testing.T) {
	st := NewStore()
	a, _ := identity.NormalizeAddress("0x1111111111111111111111111111111111111111")
	b, _ := identity.NormalizeAddress("0x2222222222222222222222222222222222222222")
	wa, err := st.Create(a, config.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.AssertOwner(wa.ID, a); err != nil {
		t.Fatal(err)
	}
	if err := st.AssertOwner(wa.ID, b); err == nil {
		t.Fatal("other")
	}
}
