package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestResearchRefusesUnpinned(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	t.Setenv("PIT_DIRECT_SPONSOR_FILE", "")
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	_, err := RunWorkspaceResearchStage(dir, "BTC", nil, nil)
	if err == nil || err.Error() != "policy_changed" {
		t.Fatal(err)
	}
}

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
