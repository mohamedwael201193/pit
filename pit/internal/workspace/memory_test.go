package workspace

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestMemoryCrossAccess(t *testing.T) {
	st := NewStore()
	a, _ := identity.NormalizeAddress("0x1111111111111111111111111111111111111111")
	b, _ := identity.NormalizeAddress("0x2222222222222222222222222222222222222222")
	wa, err := st.Create(a, config.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	wb, err := st.Create(b, config.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.PutMemory(wa.ID, "thesis", "note", []byte("private-a")); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetMemory(wb.ID, "thesis", "note"); err == nil {
		t.Fatal("cross")
	}
	got, err := st.GetMemory(wa.ID, "thesis", "note")
	if err != nil || string(got.Payload) != "private-a" {
		t.Fatal(err)
	}
}
