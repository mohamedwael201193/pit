package companion

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
	"github.com/mohamedwael201193/pit/internal/calib"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskcmd"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/version"
)

func (h *Hub) localActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"events": readActivity(h.Dir, 200), "sign": false, "trade": false,
	})
}

func (h *Hub) localPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	st, err := cli.Load(h.Dir)
	if err != nil {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "unbound", "positions": []any{}, "sign": false, "trade": false})
		return
	}
	user := strings.TrimSpace(st.Wallet)
	agent := ""
	if sf, serr := cli.LoadSession(h.Dir); serr == nil {
		agent = strings.ToLower(strings.TrimSpace(sf.AgentAddr))
	}
	if user == "" {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "unbound", "positions": []any{}, "sign": false, "trade": false})
		return
	}
	if agent != "" && strings.EqualFold(user, agent) {
		writeLocal(w, http.StatusOK, map[string]any{
			"ok": false, "error": "WRONG_ACCOUNT_QUERY", "note": "No positions on this Hyperliquid account.",
			"positions": []any{}, "sign": false, "trade": false,
		})
		return
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": TermWrongNetwork, "positions": []any{}, "sign": false, "trade": false})
		return
	}
	c := hl.New(config.For(net))
	rows, acct, perr := c.Clearinghouse(user)
	if perr != nil {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "HYPERLIQUID_OUTAGE", "positions": []any{}, "sign": false, "trade": false, "account": user})
		return
	}
	pol := policy.Default()
	out := make([]map[string]any, 0, len(rows))
	for _, p := range rows {
		mark := 0.0
		if b, berr := c.PublicBook(p.Coin); berr == nil {
			mark = b.MarkPx
		}
		out = append(out, map[string]any{
			"coin": p.Coin, "sz": p.Sz, "entryPx": p.EntryPx, "unrealizedPnl": p.UnrealizedPnl,
			"markPx": mark, "leverage": p.Leverage, "marginUsed": p.MarginUsed,
			"source": "hyperliquid.clearinghouseState", "account": user,
			"policyClipUsd": pol.MaxClipUSD,
		})
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok": true, "account": user, "queried": "master", "positions": out, "sign": false, "trade": false,
		"lastOrder": cli.LoadLastOrder(h.Dir),
		"summary": map[string]any{
			"accountValue": acct.AccountValue, "totalMarginUsed": acct.TotalMarginUsed,
			"totalNtlPos": acct.TotalNtlPos, "withdrawable": acct.Withdrawable,
		},
	})
}

func (h *Hub) localChat(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Text   string `json:"text"`
		Thread string `json:"thread"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	thread := strings.TrimSpace(body.Thread)
	if thread == "" {
		thread = "desk"
	}
	if secretful(body.Text) {
		writeLocal(w, http.StatusOK, map[string]any{
			"ok": false, "error": "secret_refused", "reply": "PIT will not store or send secrets in chat.",
			"execute": false, "sign": false, "trade": false,
		})
		return
	}
	parsed := h.decorateChat(deskcmd.Parse(body.Text))
	if parsed.Tool == "memory.forget" {
		forgetMemoryFiles(h.Dir)
		appendActivity(h.Dir, activityEvent{WorkspaceID: workspaceID(h.Dir), Kind: "memory.forgot", Action: "forget", Status: "ok"})
	}
	if parsed.Tool == "mission.stop" {
		auto.Stop(h.Dir, "chat_stop")
		appendActivity(h.Dir, activityEvent{WorkspaceID: workspaceID(h.Dir), Kind: "mission.stopped", Action: "stop", Status: "chat_stop", Reason: "chat_stop"})
		parsed.Reply = h.replyMissionStop()
	}
	if parsed.StartResearch {
		if parsed.Coin == "" {
			parsed.Coin = h.pickBestCoin()
		}
		parsed.Reply = parsed.Reply + " Chat cannot AUTHORIZE. The desk will start the sealed job on this computer."
	}
	appendChatThread(h.Dir, "user", body.Text, "", thread)
	appendChatThread(h.Dir, "pit", parsed.Reply, parsed.Tool, thread)
	writeLocal(w, http.StatusOK, map[string]any{
		"ok": true, "reply": parsed.Reply, "tool": parsed.Tool, "mutate": parsed.Mutate,
		"execute": false, "start_research": parsed.StartResearch, "coin": parsed.Coin,
		"navigate": parsed.Navigate, "open_url": parsed.OpenURL, "thread": thread, "sign": false, "trade": false,
	})
}

func (h *Hub) localChatLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	thread := strings.TrimSpace(r.URL.Query().Get("thread"))
	writeLocal(w, http.StatusOK, map[string]any{"messages": readChatThread(h.Dir, thread, 80), "sign": false, "trade": false})
}

func (h *Hub) localMemoryForget(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	forgetMemoryFiles(h.Dir)
	appendActivity(h.Dir, activityEvent{WorkspaceID: workspaceID(h.Dir), Kind: "memory.forgot", Action: "forget", Status: "ok"})
	writeLocal(w, http.StatusOK, map[string]any{"ok": true, "sign": false, "trade": false})
}

func (h *Hub) localCalibration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	n := 0
	if raw, err := os.ReadFile(filepath.Join(h.Dir, "forecasts.json")); err == nil {
		var arr []any
		if json.Unmarshal(raw, &arr) == nil {
			n = len(arr)
		}
	}
	need := calib.NeedResolved()
	copy := "NOT ENOUGH DATA"
	if n >= need {
		copy = "Calibration uses resolved forecasts only."
	}
	_ = calib.RefuseSparse(n)
	writeLocal(w, http.StatusOK, map[string]any{
		"n": n, "need": need, "copy": copy, "enough": n >= need, "sign": false, "trade": false,
	})
}

func (h *Hub) localSecurity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	checks := cli.Doctor(h.Dir)
	named := map[string]cli.Check{}
	for _, c := range checks {
		named[c.Name] = c
	}
	sessionAlive := false
	if st, err := cli.Load(h.Dir); err == nil {
		if sf, serr := cli.LoadSession(h.Dir); serr == nil {
			sessionAlive = session.Alive(sf.Meta().Session(), time.Now().UnixMilli()) && !st.Kill
		}
	}
	tee := named["tee"]
	teeState := "READY"
	teeWhy := "No sealed research yet. Idle is not a failure."
	if tee.OK {
		teeWhy = tee.Detail
	} else if strings.Contains(strings.ToLower(tee.Detail), "fail") || strings.Contains(strings.ToLower(tee.Detail), "invalid") {
		teeState = "BLOCKED"
		teeWhy = tee.Detail
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"domains": []map[string]any{
			securityDomain("wallet", named["wallet"], "Connect wallet", "Open pairing"),
			{"id": "compute", "state": map[bool]string{true: "READY", false: "NEEDS ACTION"}[named["direct_auth"].OK], "why": named["direct_auth"].Detail, "means": "Sealed research needs Protect my strategy and Direct funds.", "do": "Protect my strategy, then Open 0G Private Compute", "href": "https://pc.0g.ai/sdk/dashboard/funds", "hrefLabel": "Open 0G Private Compute"},
			{"id": "tee", "state": teeState, "why": teeWhy, "means": "VerifyE2EE must match the on-chain signer after a committee.", "do": "Run sealed research when compute is ready", "href": "", "hrefLabel": ""},
			{"id": "hyperliquid", "state": map[bool]string{true: "READY", false: "NEEDS ACTION"}[named["hl_agent"].OK], "why": named["hl_agent"].Detail, "means": "Hyperliquid must list the PIT agent. PIT cannot withdraw.", "do": "Approve PIT on Hyperliquid", "href": "https://app.hyperliquid.xyz/API", "hrefLabel": "Open Hyperliquid API"},
			{"id": "session", "state": map[bool]string{true: "READY", false: "NEEDS ACTION"}[sessionAlive], "why": "Local order/cancel session on this computer.", "means": "PIT still cannot withdraw.", "do": "Create / refresh secure session", "href": "", "hrefLabel": "Create session"},
			securityDomain("policy", named["policy"], "Pin policy", "Open Policy"),
			securityDomain("storage", named["storage"], "Install official storage client", "Open documentation"),
			{"id": "identity", "state": "READY", "why": "Desk identity is optional. Transfer of Agentic ID is not live on mainnet.", "means": "Trading does not wait on mint.", "do": "None required", "href": "", "hrefLabel": ""},
			{"id": "workspace", "state": "READY", "why": "This computer’s workspace is isolated by wallet.", "means": "Another wallet cannot read this session.", "do": "", "href": "", "hrefLabel": ""},
		},
		"sign": false, "trade": false, "version": version.Number,
	})
}

func securityDomain(id string, c cli.Check, do, hrefLabel string) map[string]any {
	state := "NEEDS ACTION"
	if c.Name == "" {
		state = "NEEDS ACTION"
	} else if c.OK {
		state = "READY"
	}
	return map[string]any{
		"id": id, "state": state, "why": c.Detail, "means": do, "do": do, "hrefLabel": hrefLabel,
	}
}

func (h *Hub) localIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"itransfer": "UNAVAILABLE",
		"iclone":    "UNAVAILABLE",
		"note":      "Transfer of Agentic ID is not live on mainnet. Trading does not require it.",
		"sign":      false, "trade": false,
	})
}

func (h *Hub) localUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	h.researchMu.Lock()
	running := h.job.running
	h.researchMu.Unlock()
	writeLocal(w, http.StatusOK, map[string]any{
		"version": version.Number, "research_running": running,
		"restart_allowed": !running, "signed": false, "authenticode": false,
		"note": "This build is checksum-verified, not OS-signed.",
		"sign": false, "trade": false,
	})
}

func (h *Hub) localExplain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	snap := h.snapshotResearch()
	snap["what_host_saw"] = "Public Hyperliquid mark/oracle/funding/open interest for the researched market."
	snap["why_size"] = "Host sized to policy clip and venue minimum. Model size was ignored."
	snap["private_prompt"] = false
	writeLocal(w, http.StatusOK, snap)
}

type chatLine struct {
	TS     int64  `json:"ts"`
	Role   string `json:"role"`
	Text   string `json:"text"`
	Tool   string `json:"tool,omitempty"`
	Thread string `json:"thread,omitempty"`
	Sign   bool   `json:"sign"`
	Trade  bool   `json:"trade"`
}

func chatPath(dir string) string { return filepath.Join(dir, "chat-transcript.jsonl") }

func appendChat(dir, role, text, tool string) {
	appendChatThread(dir, role, text, tool, "desk")
}

func appendChatThread(dir, role, text, tool, thread string) {
	if secretful(text) {
		return
	}
	if strings.TrimSpace(thread) == "" {
		thread = "desk"
	}
	raw, err := json.Marshal(chatLine{TS: time.Now().UnixMilli(), Role: role, Text: text, Tool: tool, Thread: thread})
	if err != nil {
		return
	}
	sealed, err := sealBytes(dir, raw)
	if err != nil {
		return
	}
	f, err := os.OpenFile(chatPath(dir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write([]byte(sealed + "\n"))
	if role == "user" {
		touchThread(dir, thread, text)
	}
}

func readChat(dir string, limit int) []chatLine {
	return readChatThread(dir, "", limit)
}

func readChatThread(dir, thread string, limit int) []chatLine {
	raw, err := os.ReadFile(chatPath(dir))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return nil
	}
	var out []chatLine
	for _, ln := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		plain, err := openBytes(dir, ln)
		if err != nil || secretful(string(plain)) {
			continue
		}
		var row chatLine
		if json.Unmarshal(plain, &row) != nil {
			continue
		}
		if row.Thread == "" {
			row.Thread = "desk"
		}
		if thread != "" && row.Thread != thread {
			continue
		}
		out = append(out, row)
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out
}

func (h *Hub) localModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	if err == nil && strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": TermWrongNetwork, "models": []any{}, "sign": false, "trade": false})
		return
	}
	sku := compute.ForNetwork(net)
	label := "Private + Unverified"
	if sku.ProvenE2EE {
		label = "Private + Verified"
	}
	direct := map[string]any{
		"model": sku.Model, "verifiability": sku.Verifiability, "proven_e2ee": sku.ProvenE2EE,
		"label": label, "path": "Direct", "provider": sku.Provider,
		"role_separation": true, "private_book": sku.ProvenE2EE,
		"capability": "sealed private research",
		"note":       "Private only when TeeML is proven on this network. Role separation on one SKU, not three independent models.",
		"latency":    "Live sealed round-trip. Typical 30–90s per role when the provider is live. Not a timer.",
		"cost":       "Estimated ~3 0G locked for one sealed committee. Public catalog pricing is not this path.",
	}
	private := []map[string]any{}
	other := []map[string]any{}
	if sku.ProvenE2EE {
		private = append(private, direct)
	} else {
		other = append(other, direct)
	}
	hostChat := map[string]any{
		"model": "host-parsed", "label": "Desk command", "path": "host",
		"private_book": false, "proven_e2ee": false, "capability": "chat",
		"note": "Desk chat is host-parsed on this computer. It is not a sealed model and cannot AUTHORIZE.",
	}
	other = append(other, hostChat)
	unsupported := []map[string]any{{
		"model": "public-catalog", "label": "Public compute catalog", "path": "unsupported",
		"private_book": false, "proven_e2ee": false, "capability": "not used",
		"note": "SKUs that only exist on the public compute catalog are unsupported for private research. PIT never routes the sealed book there.",
	}}
	models := private
	if len(models) == 0 {
		models = []map[string]any{}
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok": true, "network": string(net), "sign": false, "trade": false,
		"note":   "Private book uses Direct only. Presence in a catalog does not make a model private.",
		"models": models,
		"groups": map[string]any{
			"private_verified": private,
			"other_chat":       other,
			"unsupported":      unsupported,
		},
	})
}

func secretful(s string) bool {
	low := strings.ToLower(s)
	if strings.Contains(low, "app-sk-") || strings.Contains(low, "sk-") && strings.Contains(low, "bearer") {
		return true
	}
	if strings.Contains(low, "mnemonic") || strings.Contains(low, "seed phrase") {
		return true
	}
	return false
}
