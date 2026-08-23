package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestPinWorkspace(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	p := policy.Default()
	if _, err := PinWorkspace(dir, ws, p); err != nil {
		t.Fatal(err)
	}
	if err := CheckPinned(dir, ws, p); err != nil {
		t.Fatal(err)
	}
	p.MaxClipUSD = 99
	if err := CheckPinned(dir, ws, p); err == nil {
		t.Fatal("policy_changed")
	}
	if _, err := PinWorkspace(dir, "not-a-uuid", p); err == nil {
		t.Fatal("id")
	}
}
