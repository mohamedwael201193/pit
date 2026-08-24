package redteam

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestTenantMix(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	if err := TenantMix(a, b); err != nil {
		t.Fatal(err)
	}
	if err := TenantMix(a, a); err == nil {
		t.Fatal("same")
	}
}
