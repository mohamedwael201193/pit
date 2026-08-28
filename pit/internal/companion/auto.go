package companion

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (h *Hub) autoLoop() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for range t.C {
		h.autoTick()
	}
}

func (h *Hub) autoTick() {
	p := auto.Load(h.Dir)
	m := auto.LoadMission(h.Dir)
	if m.Mode == auto.ModeResearch || m.Mode == auto.ModeGuarded {
		p.Watch = true
		if m.Mode == auto.ModeResearch || m.Mode == auto.ModeGuarded {
			p.AutoResearch = true
		}
	}
	if !p.Watch && m.Mode == auto.ModeManual {
		return
	}
	now := time.Now().Unix()
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	kill := false
	if err == nil {
		if strings.TrimSpace(st.Network) != "" {
			netName = st.Network
		}
		kill = st.Kill
	}
	sessOK := false
	if sf, serr := cli.LoadSession(h.Dir); serr == nil {
		sessOK = session.Alive(sf.Meta().Session(), time.Now().UnixMilli()) && !kill
	}
	openN := h.openPositionCount()
	pol := policy.Default()
	if why := auto.StopReason(m, now, kill, sessOK, openN, 0, pol); why != "" && m.Mode == auto.ModeGuarded {
		auto.Stop(h.Dir, why)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.stopped", Action: "stop", Status: why, Reason: why,
		})
		return
	}
	if p.LastScanUnix != 0 && now-p.LastScanUnix < int64(p.CadenceMinutes)*60 {
		return
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		auto.RecordAction(h.Dir, "scan_failed", "", "", "", "market_unavailable")
		return
	}
	cands, lerr := watch.Live(hl.New(config.For(net)), watch.PolicyForWatch())
	p.LastScanUnix = now
	m = auto.LoadMission(h.Dir)
	m.NextScanUnix = now + int64(p.CadenceMinutes)*60
	if lerr != nil {
		m.LastAction = "scan_failed"
		m.LastStop = "market_unavailable"
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		return
	}
	filtered := cands
	if len(p.Markets) > 0 {
		allow := map[string]bool{}
		for _, x := range p.Markets {
			allow[strings.ToUpper(strings.TrimSpace(x))] = true
		}
		next := make([]watch.Candidate, 0, len(cands))
		for _, c := range cands {
			if allow[c.Coin] {
				next = append(next, c)
			}
		}
		filtered = next
	}
	if len(m.Assets) > 0 {
		allow := map[string]bool{}
		for _, x := range m.Assets {
			allow[strings.ToUpper(strings.TrimSpace(x))] = true
		}
		next := make([]watch.Candidate, 0, len(filtered))
		for _, c := range filtered {
			if allow[c.Coin] {
				next = append(next, c)
			}
		}
		filtered = next
	}
	var pick *watch.Candidate
	for i := range filtered {
		if m.MinLiquidityUSD > 0 && filtered[i].Book.OpenInterest < m.MinLiquidityUSD {
			continue
		}
		if m.PauseUncertain && len(filtered[i].Risk) > 0 {
			continue
		}
		if auto.Matches(p.Trigger, filtered[i].Reason) {
			pick = &filtered[i]
			break
		}
	}
	if pick == nil {
		m.LastAction = "no_opportunity"
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		return
	}
	m.BestCoin = pick.Coin
	m.BestWhy = watch.WhyHuman(*pick)
	m.LastAction = "scanned:" + pick.Coin
	view := watch.Public([]watch.Candidate{*pick}, netName)
	if view.BestWhy != "" {
		m.BestWhy = view.BestWhy
	}
	if p.Notify && p.LastNotifyCoin != pick.Coin {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "opportunity", Market: pick.Coin,
			Action: "discover", Status: "ready", Reason: pick.Reason,
		})
		p.LastNotifyCoin = pick.Coin
	}
	wantResearch := p.AutoResearch || m.Mode == auto.ModeResearch || m.Mode == auto.ModeGuarded
	if wantResearch && p.LastResearchCoin != pick.Coin {
		h.researchMu.Lock()
		running := h.job.running
		h.researchMu.Unlock()
		if !running {
			p.LastResearchCoin = pick.Coin
			note := "automation_stops_at_human_approval"
			if m.Mode == auto.ModeGuarded {
				note = "guarded_will_execute_if_eligible"
			}
			appendActivity(h.Dir, activityEvent{
				WorkspaceID: workspaceID(h.Dir), Kind: "automation.prepared", Market: pick.Coin,
				Action: "prepare", Status: "research", Reason: note,
			})
			m.LastAction = "research:" + pick.Coin
			_ = auto.SaveMission(h.Dir, m)
			_ = auto.Save(h.Dir, p)
			h.beginResearch(pick.Coin)
			return
		}
	}
	_ = auto.SaveMission(h.Dir, m)
	_ = auto.Save(h.Dir, p)
}

func (h *Hub) openPositionCount() int {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return 0
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return 0
	}
	rows, _, perr := hl.New(config.For(net)).Clearinghouse(st.Wallet)
	if perr != nil {
		return 0
	}
	n := 0
	for _, row := range rows {
		if sz, err := strconv.ParseFloat(strings.TrimSpace(row.Sz), 64); err == nil && sz != 0 {
			n++
		}
	}
	return n
}

func (h *Hub) maybeGuardedExecute(hash, coin string, started time.Time) {
	if strings.TrimSpace(hash) == "" {
		return
	}
	st, err := cli.Load(h.Dir)
	kill := false
	if err == nil {
		kill = st.Kill
	}
	sessOK := false
	if sf, serr := cli.LoadSession(h.Dir); serr == nil {
		sessOK = session.Alive(sf.Meta().Session(), time.Now().UnixMilli()) && !kill
	}
	pol := policy.Default()
	g := auto.ExecGate{
		PreviewHash: hash,
		StartedUnix: started.Unix(),
		OpenCount:   h.openPositionCount(),
		SessionOK:   sessOK,
		Kill:        kill,
		Now:         time.Now().Unix(),
		Policy:      pol,
		Coin:        coin,
	}
	if err := auto.AllowHostExecute(h.Dir, g); err != nil {
		auto.RecordAction(h.Dir, "refused:"+err.Error(), coin, hash, "", err.Error())
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.refused", Market: coin,
			Action: "guarded", Status: err.Error(), PreviewHash: hash, Reason: err.Error(),
		})
		return
	}
	got := cli.ExecuteDeskOrder(h.Dir, cli.ConfirmToken, hash)
	if got.Error != "" {
		auto.RecordAction(h.Dir, "exec_failed:"+got.Error, coin, hash, "", got.Error)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "order.rejected", Market: coin,
			Action: "guarded", Status: got.Error, PreviewHash: hash, Reason: got.Error,
		})
		return
	}
	auto.RecordAction(h.Dir, "executed", coin, hash, got.OID, "")
	appendActivity(h.Dir, activityEvent{
		WorkspaceID: workspaceID(h.Dir), Kind: "order.submitted", Market: got.Market,
		Action: "guarded", Status: "submitted", OID: got.OID, PreviewHash: got.Hash,
		Reason: "guarded_autonomy",
	})
}

func (h *Hub) localAutomation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
			http.Error(w, "origin_denied", http.StatusForbidden)
			return
		}
		p := auto.Load(h.Dir)
		p.Execute = false
		body := auto.Public(h.Dir)
		body["prefs"] = p
		body["note"] = "Watch, discover, research, notify, and prepare. Execute requires Guarded Autonomy on Automation."
		writeLocal(w, http.StatusOK, body)
		return
	}
	if !desktopOnly(w, r) {
		return
	}
	var body auto.Prefs
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Execute {
		writeLocal(w, http.StatusOK, map[string]any{
			"ok": false, "error": auto.RefuseExecute().Error(), "execute": false, "sign": false, "trade": false,
		})
		return
	}
	body.Execute = false
	cur := auto.Load(h.Dir)
	body.LastScanUnix = cur.LastScanUnix
	body.LastNotifyCoin = cur.LastNotifyCoin
	body.LastResearchCoin = cur.LastResearchCoin
	if err := auto.Save(h.Dir, body); err != nil {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": err.Error(), "execute": false, "sign": false, "trade": false})
		return
	}
	got := auto.Load(h.Dir)
	writeLocal(w, http.StatusOK, map[string]any{"ok": true, "prefs": got, "execute": false, "sign": false, "trade": false})
}

func (h *Hub) localMission(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
			http.Error(w, "origin_denied", http.StatusForbidden)
			return
		}
		writeLocal(w, http.StatusOK, auto.Public(h.Dir))
		return
	}
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Typed           string   `json:"typed"`
		Mode            string   `json:"mode"`
		Hours           int      `json:"hours"`
		MaxTrades       int      `json:"max_trades"`
		StopLossUSD     float64  `json:"stop_loss_usd"`
		MinLiquidityUSD float64  `json:"min_liquidity_usd"`
		PauseUncertain  bool     `json:"pause_uncertain"`
		Assets          []string `json:"assets"`
		Stop            bool    `json:"stop"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Stop || strings.EqualFold(strings.TrimSpace(body.Typed), auto.StopToken) {
		m := auto.Stop(h.Dir, "user_stop")
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.stopped", Action: "stop", Status: "user_stop", Reason: "user_stop",
		})
		out := auto.Public(h.Dir)
		out["mission"] = m
		writeLocal(w, http.StatusOK, out)
		return
	}
	if strings.TrimSpace(body.Typed) != "" || strings.EqualFold(body.Mode, auto.ModeGuarded) {
		hash, _ := policy.Default().Hash()
		m, err := auto.EnableGuarded(h.Dir, body.Typed, body.Hours, hash)
		if err != nil {
			writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": err.Error(), "execute": false, "sign": false, "trade": false})
			return
		}
		m.MaxTrades = body.MaxTrades
		m.StopLossUSD = body.StopLossUSD
		m.MinLiquidityUSD = body.MinLiquidityUSD
		m.PauseUncertain = body.PauseUncertain
		m.Assets = body.Assets
		_ = auto.SaveMission(h.Dir, m)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.enabled", Action: "guarded", Status: "running", Reason: "ENABLE GUARDED AUTONOMY",
		})
		writeLocal(w, http.StatusOK, auto.Public(h.Dir))
		return
	}
	if body.Mode != "" {
		m, err := auto.SetMode(h.Dir, body.Mode)
		if err != nil {
			writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": err.Error(), "execute": false, "sign": false, "trade": false})
			return
		}
		m.MaxTrades = body.MaxTrades
		m.StopLossUSD = body.StopLossUSD
		m.MinLiquidityUSD = body.MinLiquidityUSD
		m.PauseUncertain = body.PauseUncertain
		m.Assets = body.Assets
		_ = auto.SaveMission(h.Dir, m)
		writeLocal(w, http.StatusOK, auto.Public(h.Dir))
		return
	}
	m := auto.LoadMission(h.Dir)
	m.MaxTrades = body.MaxTrades
	m.StopLossUSD = body.StopLossUSD
	m.MinLiquidityUSD = body.MinLiquidityUSD
	m.PauseUncertain = body.PauseUncertain
	m.Assets = body.Assets
	_ = auto.SaveMission(h.Dir, m)
	writeLocal(w, http.StatusOK, auto.Public(h.Dir))
}
