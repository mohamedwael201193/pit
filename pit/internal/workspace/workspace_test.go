package workspace

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestIsolation(t *testing.T) {
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
		t.Fatal("ids collided")
	}
	if wa.Network == wb.Network {
		t.Fatal("networks must stay with the workspace")
	}
	if err := st.PutMemory(wa.ID, "observation", "book", []byte("secret-a")); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetMemory(wb.ID, "observation", "book"); err == nil {
		t.Fatal("B must not read A memory")
	}
	got, err := st.GetMemory(wa.ID, "observation", "book")
	if err != nil || string(got.Payload) != "secret-a" {
		t.Fatalf("A read: %v %s", err, got.Payload)
	}
	if _, err := st.Get("00000000-0000-4000-8000-000000000000"); err == nil || err.Error() != "not found" {
		t.Fatalf("uuid guess must be not found, got %v", err)
	}
	keys, err := st.ListMemory(wb.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Fatalf("B listed %v", keys)
	}
}
