package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func BindResearchPreview(dir, coin string, snap hl.BookSnapshot, pol policy.Policy, st DiskState, rep compute.AskReport) compute.AskReport {
	live, err := LiveFromDisk(dir, st.Kill, time.Now().UnixMilli())
	if err != nil {
		rep.Deny = "session_expired"
		rep.Eligible = false
		rep.Preview = map[string]any{"eligible": false, "deny": "session_expired"}
		return rep
	}
	got := engine.EvaluateCommittee(
		snap.MarkPx, snap.SzDecimals, pol.MaxClipUSD, pol.MaxClipUSD, "", coin, pol.AllowedAssets,
		1, pol.MaxLeverage, pol.KillSwitch || st.Kill, rep.Researcher, rep.Challenger, rep.Risk,
	)
	if !got.Eligible || got.Size == nil {
		deny := got.Deny
		if deny == "" {
			deny = "no_side"
		}
		rep.Deny = deny
		rep.Eligible = false
		rep.Preview = map[string]any{"eligible": false, "deny": deny, "reasons": got.Reasons}
		return rep
	}
	sum := sha256.Sum256([]byte(st.WorkspaceID + "|" + coin + "|" + strconv.FormatInt(time.Now().UnixMilli(), 10)))
	forecast := "0x" + hex.EncodeToString(sum[:])
	cloid, err := NewCloid()
	if err != nil {
		rep.Deny = err.Error()
		rep.Eligible = false
		return rep
	}
	p, err := HostPreview(coin, got.Size.Side, forecast, snap, pol, live, time.Now().UTC(), cloid, time.Now().UnixMilli())
	if err != nil {
		rep.Deny = err.Error()
		rep.Eligible = false
		rep.Preview = map[string]any{"eligible": false, "deny": err.Error()}
		return rep
	}
	hash, err := SavePreview(dir, p)
	if err != nil {
		rep.Deny = err.Error()
		rep.Eligible = false
		return rep
	}
	rep.Eligible = true
	rep.PreviewHash = hash
	rep.Preview = map[string]any{
		"eligible":      true,
		"market":        p.Market,
		"side":          p.Side,
		"sz":            p.Sz,
		"orderType":     p.OrderType,
		"limitPx":       p.LimitPx,
		"slippageBps":   p.SlippageBps,
		"policyVersion": p.PolicyVersion,
		"expiryUnixMs":  p.ExpiryUnixMs,
		"cloid":         p.Cloid,
		"forecastId":    p.ForecastID,
		"hash":          hash,
		"notionalUsd":   got.Size.NotionalUSD,
	}
	return rep
}
