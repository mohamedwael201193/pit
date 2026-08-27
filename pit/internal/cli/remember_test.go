package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/ledger"
)

func TestRememberAuthorizedOnce(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	cloid := "0x11111111111111111111111111111111"
	if err := RememberAuthorized(dir, "testnet", ws, cloid, "hash"); err != nil {
		t.Fatal(err)
	}
	if err := RememberAuthorized(dir, "testnet", ws, cloid, "hash"); err == nil {
		t.Fatal("dup")
	}
	got, err := LookupAction(dir, "testnet", ws, cloid)
	if err != nil || got.Status != ledger.StatusAuthorized {
		t.Fatalf("%v %+v", err, got)
	}
	if err := RememberPosted(dir, "testnet", ws, cloid, "42"); err != nil {
		t.Fatal(err)
	}
	got, err = LookupAction(dir, "testnet", ws, cloid)
	if err != nil || got.Status != ledger.StatusReceipt || got.OID != "42" {
		t.Fatalf("%v %+v", err, got)
	}
}
