package cli

import (
	"os"
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestLogoutRemovesSessionKeepsWorkspace(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	sf, err := CreateLocalSession(dir, ws, "mainnet", "v1")
	if err != nil {
		t.Fatal(err)
	}
	if sf.ID == "" {
		t.Fatal("id")
	}
	if err := Logout(dir, false); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadSession(dir); err == nil {
		t.Fatal("session remains")
	}
	st, err := Load(dir)
	if err != nil || st.WorkspaceID != ws {
		t.Fatal("workspace must remain")
	}
}

func TestLogoutForgetUnbinds(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "testnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	if err := Logout(dir, true); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir); err == nil {
		t.Fatal("state remains")
	}
	_ = os.DevNull
}
