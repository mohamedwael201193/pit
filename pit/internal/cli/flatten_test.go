package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
)

func TestCloseSideAndSize(t *testing.T) {
	side, sz, err := CloseSideAndSize("0.0041")
	if err != nil || side != "sell" || sz != 0.0041 {
		t.Fatalf("%s %v %v", side, sz, err)
	}
	side, sz, err = CloseSideAndSize("-0.01")
	if err != nil || side != "buy" || sz != 0.01 {
		t.Fatalf("%s %v %v", side, sz, err)
	}
	if _, _, err := CloseSideAndSize("0"); err == nil {
		t.Fatal("zero")
	}
}

func TestBindReduceOnlyCloseNeedsSession(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := BindReduceOnlyClose(dir, "ETH"); err == nil {
		t.Fatal("unbound")
	}
}

func TestReduceOnlyWireSetsR(t *testing.T) {
	p := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "sell", Sz: 0.0041, OrderType: "limit",
		LimitPx: "2500", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
		ReduceOnly: true,
	}
	raw, err := pitexec.WireFromPreview(p, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"r":true`) {
		t.Fatal(string(raw))
	}
	if err := hl.AssertActionType(raw); err != nil {
		t.Fatal(err)
	}
}

func TestHostClosePreviewIgnoresModelSize(t *testing.T) {
	book := hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, SzDecimals: 4}
	pol := policy.Default()
	sess := session.Session{ID: "s", Workspace: "w"}
	p, err := HostClosePreview("ETH", "sell", 0.0041, "0xabc", book, pol, sess, time.Now().UTC(), "0x11111111111111111111111111111111", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !p.ReduceOnly || p.Sz != 0.0041 || p.Side != "sell" {
		t.Fatalf("%+v", p)
	}
}

func TestReduceOnlyPublicNeverSigns(t *testing.T) {
	pub := ReduceOnlyPublic(engine.Preview{Market: "hyperliquid:perp:ETH", Side: "sell", Sz: 0.0041}, "0xabc")
	if pub["sign"] != false || pub["trade"] != false || pub["reduceOnly"] != true {
		t.Fatal(pub)
	}
	if pub["kind"] != "reduce_only_close" {
		t.Fatal(pub["kind"])
	}
}
