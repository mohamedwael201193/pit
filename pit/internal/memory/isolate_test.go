package memory

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestIsolateWorkspaces(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	if err := Isolate(config.Mainnet, a, b, "thesis", "note"); err != nil {
		t.Fatal(err)
	}
	if err := Isolate(config.Mainnet, a, a, "thesis", "note"); err == nil {
		t.Fatal("same")
	}
}
