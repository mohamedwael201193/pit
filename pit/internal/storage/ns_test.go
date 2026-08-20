package storage

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestObjectKeyIsolation(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	ka, err := ObjectKey(config.Mainnet, a, "memory", "book")
	if err != nil {
		t.Fatal(err)
	}
	kb, err := ObjectKey(config.Mainnet, b, "memory", "book")
	if err != nil {
		t.Fatal(err)
	}
	if ka == kb {
		t.Fatal("keys collided")
	}
	if err := AssertWorkspace(ka, a); err != nil {
		t.Fatal(err)
	}
	if err := AssertWorkspace(ka, b); err == nil {
		t.Fatal("cross key")
	}
	if err := RequireHexKey("not-hex"); err == nil {
		t.Fatal("key prefix")
	}
}
