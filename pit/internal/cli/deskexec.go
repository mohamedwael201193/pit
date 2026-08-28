package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pitexec "github.com/mohamedwael201193/pit/internal/exec"
)

type OrderResult struct {
	OK     bool    `json:"ok"`
	Posted bool    `json:"posted"`
	OID    string  `json:"oid,omitempty"`
	Cloid  string  `json:"cloid,omitempty"`
	Hash   string  `json:"hash,omitempty"`
	Market string  `json:"market,omitempty"`
	Side   string  `json:"side,omitempty"`
	Sz     float64 `json:"sz,omitempty"`
	Agent  string  `json:"agent,omitempty"`
	Error  string  `json:"error,omitempty"`
	Sign   bool    `json:"sign"`
	Trade  bool    `json:"trade"`
	Order  bool    `json:"order"`
	Cancel bool    `json:"cancel"`
}

func LastOrderPath(dir string) string {
	return filepath.Join(dir, "last-order.json")
}

func LoadLastOrder(dir string) map[string]any {
	b, err := os.ReadFile(LastOrderPath(dir))
	if err != nil || strings.Contains(strings.ToLower(string(b)), "app-sk-") {
		return nil
	}
	var m map[string]any
	if json.Unmarshal(b, &m) != nil {
		return nil
	}
	return m
}

func saveLastOrder(dir string, body map[string]any) {
	raw, err := json.MarshalIndent(body, "", "  ")
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return
	}
	_ = os.WriteFile(LastOrderPath(dir), raw, 0o600)
}

func ConfirmDeskAuthorize(typed string, sessionAlive bool) error {
	if !sessionAlive {
		return fmt.Errorf("session_expired")
	}
	if strings.TrimSpace(typed) != ConfirmToken {
		return fmt.Errorf("need_exact_AUTHORIZE")
	}
	return nil
}

func ExecuteDeskOrder(dir, typed, presentedHash string) OrderResult {
	out := OrderResult{Sign: false, Trade: false, Order: true, Cancel: true}
	st, err := Load(dir)
	if err != nil {
		out.Error = "unbound"
		return out
	}
	live, err := LiveFromDisk(dir, st.Kill, time.Now().UnixMilli())
	if err != nil {
		out.Error = "session_expired"
		return out
	}
	if err := ConfirmDeskAuthorize(typed, true); err != nil {
		out.Error = err.Error()
		return out
	}
	card, hash, err := LoadPreview(dir)
	if err != nil {
		out.Error = "preview_required"
		return out
	}
	if presentedHash != "" && presentedHash != hash {
		out.Error = "preview_hash_mismatch"
		return out
	}
	if card.WorkspaceID != live.Workspace || card.SessionID != live.ID {
		out.Error = "wrong_workspace"
		return out
	}
	now := time.Now().UnixMilli()
	if err := pitexec.RequirePreview(card, hash, now); err != nil {
		out.Error = err.Error()
		return out
	}
	h := pitexec.HashForAuthorize(hash)
	used := map[string]struct{}{}
	if err := pitexec.Prepare(pitexec.Intent{Action: "order", Preview: card, Hash: h, Workspace: live.Workspace}, now, used); err != nil {
		out.Error = err.Error()
		return out
	}
	coin, err := pitexec.CoinFromMarket(card.Market)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	book, err := LiveAsset(st.Network, coin)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	raw, err := pitexec.WireFromPreview(card, book.Asset)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	if err := RememberAuthorized(dir, st.Network, live.Workspace, card.Cloid, h); err != nil {
		out.Error = err.Error()
		return out
	}
	env, signErr := SignBound(dir, live, st.Network, raw, now)
	if signErr != nil || !env.Signed() {
		out.Error = "exchange_unsigned"
		return out
	}
	linked, linkErr := LiveLinked(st.Network, st.Wallet, live.Workspace, live.AgentAddr, now)
	if linkErr != nil {
		out.Error = linkErr.Error()
		out.Agent = live.AgentAddr
		return out
	}
	out.Agent = live.AgentAddr
	if err := pitexec.RefusePostUntilLinked(linked); err != nil {
		out.Error = "approveAgent_required"
		return out
	}
	body, err := PostLinked(st.Network, env, linked, h)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	oid := pitexec.ReceiptOID(body)
	_ = RememberPosted(dir, st.Network, live.Workspace, card.Cloid, oid)
	out.OK = true
	out.Posted = true
	out.OID = oid
	out.Cloid = card.Cloid
	out.Hash = hash
	out.Market = card.Market
	out.Side = card.Side
	out.Sz = card.Sz
	saveLastOrder(dir, map[string]any{
		"ok": true, "posted": true, "oid": oid, "cloid": card.Cloid, "hash": hash,
		"market": card.Market, "side": card.Side, "sz": card.Sz, "sign": false, "trade": false,
	})
	return out
}

func ExecuteDeskCancel(dir, typed string) OrderResult {
	out := OrderResult{Sign: false, Trade: false, Order: true, Cancel: true}
	if strings.TrimSpace(typed) != ConfirmToken {
		out.Error = "need_exact_AUTHORIZE"
		return out
	}
	st, err := Load(dir)
	if err != nil {
		out.Error = "unbound"
		return out
	}
	live, err := LiveFromDisk(dir, st.Kill, time.Now().UnixMilli())
	if err != nil {
		out.Error = "session_expired"
		return out
	}
	card, hash, err := LoadPreview(dir)
	if err != nil {
		out.Error = "preview_required"
		return out
	}
	cloid := card.Cloid
	rec, err := LookupAction(dir, st.Network, live.Workspace, cloid)
	if err != nil || rec.OID == "" {
		if last := LoadLastOrder(dir); last != nil {
			if posted, _ := last["cloid"].(string); posted != "" {
				cloid = posted
				rec, err = LookupAction(dir, st.Network, live.Workspace, cloid)
			}
		}
	}
	if err != nil || rec.OID == "" {
		out.Error = "cancel_requires_posted_order"
		return out
	}
	coin, err := pitexec.CoinFromMarket(card.Market)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	book, err := LiveAsset(st.Network, coin)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	raw, err := CancelWire(book.Asset, cloid)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	now := time.Now().UnixMilli()
	env, signErr := SignBound(dir, live, st.Network, raw, now)
	if signErr != nil || !env.Signed() {
		out.Error = "exchange_unsigned"
		return out
	}
	linked, linkErr := LiveLinked(st.Network, st.Wallet, live.Workspace, live.AgentAddr, now)
	if linkErr != nil {
		out.Error = linkErr.Error()
		out.Agent = live.AgentAddr
		return out
	}
	out.Agent = live.AgentAddr
	if err := pitexec.RefusePostUntilLinked(linked); err != nil {
		out.Error = "approveAgent_required"
		return out
	}
	if _, err := PostLinked(st.Network, env, linked, pitexec.HashForAuthorize(hash)); err != nil {
		out.Error = err.Error()
		return out
	}
	out.OK = true
	out.Posted = true
	out.Cloid = cloid
	out.Hash = hash
	out.Market = card.Market
	saveLastOrder(dir, map[string]any{
		"ok": true, "posted": true, "cancelled": true, "cloid": cloid, "hash": hash,
		"market": card.Market, "sign": false, "trade": false,
	})
	return out
}
