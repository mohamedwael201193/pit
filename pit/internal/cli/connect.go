package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func BindConnectionPreview(dir, coin string) (engine.Preview, string, error) {
	st, err := Load(dir)
	if err != nil {
		return engine.Preview{}, "", fmt.Errorf("unbound")
	}
	now := time.Now()
	live, err := LiveFromDisk(dir, st.Kill, now.UnixMilli())
	if err != nil {
		return engine.Preview{}, "", fmt.Errorf("session_expired")
	}
	linked, _, err := LookupAgent(st.Network, st.Wallet, live.Workspace, live.AgentAddr, now.UnixMilli())
	if err != nil {
		return engine.Preview{}, "", err
	}
	if !linked {
		return engine.Preview{}, "", fmt.Errorf("approveAgent_required")
	}
	p := policy.Default()
	_ = CheckPinned(dir, st.WorkspaceID, p)
	want := strings.ToUpper(strings.TrimSpace(coin))
	if want == "" {
		want = "ETH"
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return engine.Preview{}, "", err
	}
	snap, err := hl.New(config.For(net)).PublicBook(want)
	if err != nil || snap.MarkPx <= 0 {
		return engine.Preview{}, "", fmt.Errorf("empty_envelope")
	}
	sum := sha256.Sum256([]byte("connection-test|" + st.WorkspaceID + "|" + want + "|" + strconv.FormatInt(now.UnixMilli(), 10)))
	forecast := "0x" + hex.EncodeToString(sum[:])
	cloid, err := NewCloid()
	if err != nil {
		return engine.Preview{}, "", err
	}
	card, err := HostPreview(want, "buy", forecast, snap, p, live, now.UTC(), cloid, now.UnixMilli())
	if err != nil {
		return engine.Preview{}, "", err
	}
	hash, err := SavePreview(dir, card)
	if err != nil {
		return engine.Preview{}, "", err
	}
	return card, hash, nil
}

func ConnectionPreviewPublic(p engine.Preview, hash string) map[string]any {
	return map[string]any{
		"ok":           true,
		"kind":         "connection_test",
		"eligible":     true,
		"market":       p.Market,
		"side":         p.Side,
		"sz":           p.Sz,
		"orderType":    p.OrderType,
		"limitPx":      p.LimitPx,
		"hash":         hash,
		"cloid":        p.Cloid,
		"expiryUnixMs": p.ExpiryUnixMs,
		"forecastId":   p.ForecastID,
		"note":         "Connection test. This is not a research recommendation. Host sized a policy clip. Type AUTHORIZE on this computer to send it.",
		"sign":         false,
		"trade":        false,
		"withdraw":     false,
	}
}
