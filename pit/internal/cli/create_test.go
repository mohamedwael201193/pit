package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/session"
)

func TestCreateLocalSessionDoesNotExport(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	sf, err := CreateLocalSession(dir, ws, "mainnet", "v1")
	if err != nil {
		t.Fatal(err)
	}
	if sf.AgentAddr == "" || sf.Expires == 0 {
		t.Fatalf("%+v", sf)
	}
	s, err := LiveFromDisk(dir, false, sf.Expires-1)
	if err != nil {
		t.Fatal(err)
	}
	if s.Workspace != ws {
		t.Fatal(s.Workspace)
	}
	if _, err := LiveFromDisk(dir, false, sf.Expires); err == nil {
		t.Fatal("expired")
	}
	if session.Alive(s, sf.Expires) {
		t.Fatal("boundary")
	}
}

func TestEnsureLocalSessionDoesNotRotateLiveAgent(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	first, minted, err := EnsureLocalSession(dir)
	if err != nil || !minted {
		t.Fatalf("mint %v %v", minted, err)
	}
	second, minted, err := EnsureLocalSession(dir)
	if err != nil || minted {
		t.Fatalf("reuse %v %v", minted, err)
	}
	if first.ID != second.ID || first.AgentAddr != second.AgentAddr {
		t.Fatal("rotated")
	}
	pub := SessionPublic(first)
	if pub["withdraw"] != false || pub["sign"] != false {
		t.Fatal(pub)
	}
	if _, ok := pub["private_key"]; ok {
		t.Fatal("secret")
	}
}
