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
	prev := LookupAgent
	t.Cleanup(func() { LookupAgent = prev })
	LookupAgent = func(_, _, _, _ string, _ int64) (bool, int64, error) {
		return false, 0, nil
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

func TestEnsureLocalSessionReusesApprovedExpiredAgent(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	first, minted, err := EnsureLocalSession(dir)
	if err != nil || !minted {
		t.Fatalf("mint %v %v", minted, err)
	}
	first.Expires = 1
	if err := SaveSession(dir, first); err != nil {
		t.Fatal(err)
	}
	prev := LookupAgent
	t.Cleanup(func() { LookupAgent = prev })
	LookupAgent = func(_, _, _, addr string, _ int64) (bool, int64, error) {
		if addr != first.AgentAddr {
			t.Fatalf("addr %s", addr)
		}
		return true, 1803441284611, nil
	}
	got, minted, err := EnsureLocalSession(dir)
	if err != nil || minted {
		t.Fatalf("reuse %v %v", minted, err)
	}
	if got.AgentAddr != first.AgentAddr || got.ID != first.ID {
		t.Fatal("reminted")
	}
	if got.Expires != 1803441284611 {
		t.Fatalf("expires %d", got.Expires)
	}
	if _, err := LiveFromDisk(dir, false, 2); err != nil {
		t.Fatal(err)
	}
}
