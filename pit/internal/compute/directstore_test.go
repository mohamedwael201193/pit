package compute

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/keyring"
)

func TestStoreLoadDeleteDirect(t *testing.T) {
	store, err := keyring.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ws := identity.NewWorkspaceID()
	sku := MainnetChat()
	if err := StoreDirect(store, "mainnet", ws, sku.Provider, "Bearer sk-dashboard"); err == nil {
		t.Fatal("router")
	}
	auth := "Bearer app-sk-dGVzdA=="
	if err := StoreDirect(store, "mainnet", ws, sku.Provider, auth); err != nil {
		t.Fatal(err)
	}
	got, err := LoadDirect(store, "mainnet", ws, sku.Provider)
	if err != nil || got != auth {
		t.Fatalf("%s %v", got, err)
	}
	other := identity.NewWorkspaceID()
	if _, err := LoadDirect(store, "mainnet", other, sku.Provider); err == nil {
		t.Fatal("user B")
	}
	if err := DeleteDirect(store, "mainnet", ws, sku.Provider); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadDirect(store, "mainnet", ws, sku.Provider); err == nil {
		t.Fatal("deleted")
	}
	_ = config.Mainnet
	_ = time.Now()
}
