package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
)

func TestHostPreviewIgnoresModelSize(t *testing.T) {
	ws := identity.NewWorkspaceID()
	sess := session.Session{ID: "s1", Workspace: ws, AgentAddr: "0xa", Expires: time.Now().Add(time.Hour).UnixMilli(), PolicyVer: "v1"}
	p, err := HostPreview("ETH", "buy", "f1", hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, SzDecimals: 4}, policy.Default(), sess, time.Now(), "0x11111111111111111111111111111111", 1)
	if err != nil {
		t.Fatal(err)
	}
	if p.Sz > 0.01 {
		t.Fatalf("clip %v", p.Sz)
	}
	if p.Side != "buy" || !strings.Contains(p.Market, "ETH") {
		t.Fatalf("%+v", p)
	}
}

func TestHostPreviewNeedsForecast(t *testing.T) {
	_, err := HostPreview("ETH", "buy", "", hl.BookSnapshot{MarkPx: 2500, SzDecimals: 4}, policy.Default(), session.Session{ID: "s", Workspace: "w"}, time.Now(), "0x1", 1)
	if err == nil || err.Error() != "forecast_required" {
		t.Fatalf("%v", err)
	}
}

func TestSaveLoadPreviewHash(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	p := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit",
		LimitPx: "2500", PolicyVersion: "v1", SessionID: "s", WorkspaceID: ws,
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
	}
	h, err := SavePreview(dir, p)
	if err != nil {
		t.Fatal(err)
	}
	got, gh, err := LoadPreview(dir)
	if err != nil || gh != h || got.Sz != p.Sz {
		t.Fatalf("%v %s %s", err, gh, h)
	}
	p.Sz = 9
	h2, _ := engine.CanonicalHash(p)
	if "0x"+h2 == h {
		t.Fatal("mutation must change hash")
	}
}

func TestBindResearchPreviewNeedsSession(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	st := DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}
	if err := Save(dir, st); err != nil {
		t.Fatal(err)
	}
	rep := BindResearchPreview(dir, "ETH", hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, SzDecimals: 4}, policy.Default(), st, compute.AskReport{})
	if rep.Eligible || rep.Deny != "session_expired" {
		t.Fatalf("%+v", rep)
	}
}
