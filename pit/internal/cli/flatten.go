package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
)

func CloseSideAndSize(szi string) (side string, sz float64, err error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(szi), 64)
	if err != nil || f == 0 {
		return "", 0, fmt.Errorf("no_position")
	}
	if f > 0 {
		return "sell", f, nil
	}
	return "buy", -f, nil
}

func HostClosePreview(coin, side string, sz float64, forecastID string, book hl.BookSnapshot, pol policy.Policy, sess session.Session, now time.Time, cloid string, nonce int64) (engine.Preview, error) {
	if strings.TrimSpace(forecastID) == "" {
		return engine.Preview{}, fmt.Errorf("forecast_required")
	}
	if sess.ID == "" || sess.Workspace == "" {
		return engine.Preview{}, fmt.Errorf("session_expired")
	}
	if err := hl.MarkFinite(book); err != nil {
		return engine.Preview{}, err
	}
	if sz <= 0 {
		return engine.Preview{}, fmt.Errorf("no_position")
	}
	if book.SzDecimals >= 0 && book.SzDecimals <= 8 {
		pow := math.Pow10(book.SzDecimals)
		sz = math.Round(sz*pow) / pow
	}
	p := engine.Preview{
		Market:        "hyperliquid:perp:" + strings.ToUpper(strings.TrimSpace(coin)),
		Side:          side,
		Sz:            sz,
		OrderType:     "limit",
		LimitPx:       strconv.FormatFloat(book.MarkPx, 'f', -1, 64),
		SlippageBps:   pol.MaxSlippageBps,
		PolicyVersion: pol.Version,
		SessionID:     sess.ID,
		ExpiryUnixMs:  engine.PreviewTTL(now),
		Cloid:         cloid,
		ForecastID:    forecastID,
		Nonce:         nonce,
		WorkspaceID:   sess.Workspace,
		ReduceOnly:    true,
	}
	return engine.BindPreview(p, map[string]any{"sz": 99.0, "side": "buy"})
}

func BindReduceOnlyClose(dir, coin string) (engine.Preview, string, error) {
	st, err := Load(dir)
	if err != nil {
		return engine.Preview{}, "", fmt.Errorf("unbound")
	}
	if st.Kill {
		return engine.Preview{}, "", fmt.Errorf("kill_switch")
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
	c := hl.New(config.For(net))
	rows, err := c.Positions(st.Wallet)
	if err != nil {
		return engine.Preview{}, "", fmt.Errorf("HYPERLIQUID_OUTAGE")
	}
	var szi string
	for _, row := range rows {
		if strings.EqualFold(row.Coin, want) {
			szi = row.Sz
			break
		}
	}
	side, sz, err := CloseSideAndSize(szi)
	if err != nil {
		return engine.Preview{}, "", err
	}
	snap, err := c.PublicBook(want)
	if err != nil || snap.MarkPx <= 0 {
		return engine.Preview{}, "", fmt.Errorf("empty_envelope")
	}
	sum := sha256.Sum256([]byte("reduce-only|" + st.WorkspaceID + "|" + want + "|" + strconv.FormatInt(now.UnixMilli(), 10)))
	forecast := "0x" + hex.EncodeToString(sum[:])
	cloid, err := NewCloid()
	if err != nil {
		return engine.Preview{}, "", err
	}
	card, err := HostClosePreview(want, side, sz, forecast, snap, p, live, now.UTC(), cloid, now.UnixMilli())
	if err != nil {
		return engine.Preview{}, "", err
	}
	hash, err := SavePreview(dir, card)
	if err != nil {
		return engine.Preview{}, "", err
	}
	return card, hash, nil
}

func ReduceOnlyPublic(p engine.Preview, hash string) map[string]any {
	return map[string]any{
		"ok":           true,
		"kind":         "reduce_only_close",
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
		"reduceOnly":   true,
		"note":         "Reduce-only close. This is not a research recommendation. Type AUTHORIZE on this computer. PIT cannot withdraw.",
		"sign":         false,
		"trade":        false,
		"withdraw":     false,
	}
}
