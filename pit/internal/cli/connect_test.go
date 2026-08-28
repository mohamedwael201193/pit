package cli

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestBindConnectionPreviewNeedsSession(t *testing.T) {
	dir := t.TempDir()
	if err := Save(dir, DiskState{WorkspaceID: identity.NewWorkspaceID(), Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	_, _, err := BindConnectionPreview(dir, "ETH")
	if err == nil || err.Error() != "session_expired" {
		t.Fatalf("%v", err)
	}
}

func TestBindConnectionPreviewNeedsApproval(t *testing.T) {
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
	if _, _, err := EnsureLocalSession(dir); err != nil {
		t.Fatal(err)
	}
	_, _, err := BindConnectionPreview(dir, "ETH")
	if err == nil || err.Error() != "approveAgent_required" {
		t.Fatalf("%v", err)
	}
}

func TestConnectionPreviewPublicNeverSigns(t *testing.T) {
	p := engine.Preview{Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004}
	pub := ConnectionPreviewPublic(p, "0xabc")
	if pub["sign"] != false || pub["trade"] != false || pub["withdraw"] != false {
		t.Fatal(pub)
	}
	if pub["kind"] != "connection_test" {
		t.Fatal(pub)
	}
	note, _ := pub["note"].(string)
	if !strings.Contains(strings.ToLower(note), "not a research recommendation") {
		t.Fatal(note)
	}
}
