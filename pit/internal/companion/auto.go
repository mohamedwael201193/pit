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
	"github.com/mohamedwael201193/pit/internal/feasibility"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (h *Hub) autoLoop() {
	m := auto.LoadMission(h.Dir)
	if m.Running && (m.Mode == auto.ModeGuarded || m.Mode == auto.ModeResearch) {
		h.autoTick()
	}
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for range t.C {
		h.autoTick()
	}
}

func (h *Hub) autoTick() {
	h.autoMu.Lock()
	defer h.autoMu.Unlock()
	p := auto.Load(h.Dir)
	m := auto.LoadMission(h.Dir)
	if m.Mode == auto.ModeResearch || m.Mode == auto.ModeGuarded {
		p.Watch = true
		p.AutoResearch = true
	}
	if !p.Watch && m.Mode == auto.ModeManual {
		return
	}
	now := time.Now().Unix()
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	kill := false
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		if m.Mode == auto.ModeGuarded && m.Running {
			m.Stage = "waiting"
			m.LastAction = "waiting_bind"
			_ = auto.SaveMission(h.Dir, m)
		}
		return
	}
	if strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	kill = st.Kill
	sessOK := false
	if sf, serr := cli.LoadSession(h.Dir); serr == nil {
		sessOK = session.Alive(sf.Meta().Session(), time.Now().UnixMilli()) && !kill
	}
	openN := h.openPositionCount()
	pol := cli.ActivePolicy(h.Dir)
	m.OpenPositions = openN
	m.CurrentPosition = h.positionSummary()
	if why := auto.MissionHaltReason(m, now, kill, sessOK, 0, pol); why != "" && m.Mode == auto.ModeGuarded {
		auto.Stop(h.Dir, why)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.stopped", Action: "stop", Status: why, Reason: why,
		})
		return
	}
	if why := auto.ExecBlockReason(openN, pol); why != "" && m.Mode == auto.ModeGuarded {
		if m.BlockReason != why {
			appendActivity(h.Dir, activityEvent{
				WorkspaceID: workspaceID(h.Dir), Kind: "mission.exec_blocked", Action: "gate", Status: why, Reason: why,
			})
		}
		auto.RecordBlock(h.Dir, why, m.BestCoin)
		m = auto.LoadMission(h.Dir)
	} else if m.Mode == auto.ModeGuarded && m.BlockReason != "" {
		m.BlockReason = ""
		_ = auto.SaveMission(h.Dir, m)
	}
	h.syncMissionResearch(&m)
	cadence := int64(p.CadenceMinutes) * 60
	if cadence < 60 {
		cadence = 60
	}
	due := p.LastScanUnix == 0 || now-p.LastScanUnix >= cadence || (m.GuardedEnabledUnix > 0 && p.LastScanUnix < m.GuardedEnabledUnix)
	if !due {
		if m.NextScanUnix <= now {
			m.NextScanUnix = p.LastScanUnix + cadence
			if m.NextScanUnix <= now {
				m.NextScanUnix = now + cadence
			}
		}
		if m.Stage == "" || m.Stage == "starting" || m.Stage == "guarded_enabled" {
			m.Stage = "waiting"
			m.LastAction = "waiting_cadence"
		}
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		return
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		m.Stage = "scan_failed"
		m.LastAction = "scan_failed"
		m.LastResult = "market_unavailable"
		_ = auto.SaveMission(h.Dir, m)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.scan_failed", Action: "scan", Status: "market_unavailable", Reason: "market_unavailable",
		})
		return
	}
	m.Stage = "scanning"
	m.LastAction = "scanning"
	_ = auto.SaveMission(h.Dir, m)
	cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), pol)
	p.LastScanUnix = now
	m.NextScanUnix = now + cadence
	m.ScanCount++
	if lerr != nil {
		m.LastAction = "scan_failed"
		m.LastResult = lerr.Error()
		m.Stage = "scan_failed"
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.scan_failed", Action: "scan", Status: "market_unavailable", Reason: "market_unavailable",
		})
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
	eligibleN := 0
	for _, c := range filtered {
		if c.Eligible {
			eligibleN++
		}
	}
	m.Scanned = len(cands)
	m.Eligible = eligibleN
	m.LastResult = "scanned " + strconv.Itoa(len(cands)) + ", " + strconv.Itoa(eligibleN) + " pass policy"
	appendActivity(h.Dir, activityEvent{
		WorkspaceID: workspaceID(h.Dir), Kind: "mission.scanned", Action: "scan", Status: "ok", Reason: m.LastResult,
	})
	matching := make([]watch.Candidate, 0, len(filtered))
	for i := range filtered {
		if !filtered[i].Eligible {
			continue
		}
		if m.MinLiquidityUSD > 0 && filtered[i].Book.OpenInterest < m.MinLiquidityUSD {
			continue
		}
		if m.PauseUncertain && len(filtered[i].Risk) > 0 {
			continue
		}
		if auto.Matches(p.Trigger, filtered[i].Reason) {
			matching = append(matching, filtered[i])
		}
	}
	pool := matching
	if len(pool) == 0 {
		pool = filtered
	}
	acct := h.capitalNow()
	if acct.OpenPositions == 0 {
		acct.OpenPositions = h.openPositionCount()
	}
	var pick *watch.Candidate
	execOK := false
	if best, ok := watch.BestExecutable(pool, acct, pol, sessOK, h.policyPinnedNow()); ok {
		cp := best
		pick = &cp
		execOK = true
	} else if best, ok := watch.Best(pool); ok {
		cp := best
		pick = &cp
	}
	if pick == nil {
		m.LastAction = "no_opportunity"
		m.Stage = "empty"
		m.BestCoin = ""
		block, _ := policy.ExecWhy(h.openPositionCount(), acct.BuyingPower, pol)
		if block == "" {
			block = "no_opportunity"
		}
		auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "rejected", Why: block})
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.empty", Action: "scan", Status: "no_opportunity", Reason: block, Autonomous: m.Mode == auto.ModeGuarded,
		})
		return
	}
	if !execOK {
		fit := feasibility.FitBook(pick.Book, pol, acct, sessOK, h.policyPinnedNow())
		why := fit.Gate
		if why == "" {
			why = "insufficient_margin"
		}
		auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "rejected", Coin: pick.Coin, Why: why})
		m.BlockReason = why
	} else {
		auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "detected", Coin: pick.Coin, Why: pick.Reason})
	}
	m.BestCoin = pick.Coin
	m.BestWhy = watch.WhyHuman(*pick)
	m.LastAction = "ranked:" + pick.Coin
	m.Stage = "ranked"
	view := watch.Public([]watch.Candidate{*pick}, netName)
	if view.BestWhy != "" {
		m.BestWhy = view.BestWhy
	}
	if p.Notify && p.LastNotifyCoin != pick.Coin {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "opportunity", Market: pick.Coin,
			Action: "discover", Status: "ready", Reason: pick.Reason,
			Link: venueTradeLink(netName, pick.Coin),
		})
		p.LastNotifyCoin = pick.Coin
	}
	wantResearch := p.AutoResearch || m.Mode == auto.ModeResearch || m.Mode == auto.ModeGuarded
	h.researchMu.Lock()
	running := h.job.running
	h.researchMu.Unlock()
	if running {
		m.Stage = "researching"
		m.LastAction = "research_running"
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		return
	}
	if wantResearch && p.LastResearchCoin != pick.Coin {
		note := "automation_stops_at_human_approval"
		if m.Mode == auto.ModeGuarded {
			note = "guarded_will_execute_if_eligible"
		}
		auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "researched", Coin: pick.Coin, Why: note})
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "automation.prepared", Market: pick.Coin,
			Action: "prepare", Status: "research", Reason: note,
			Link: venueTradeLink(netName, pick.Coin), Autonomous: m.Mode == auto.ModeGuarded,
		})
		p.LastResearchCoin = pick.Coin
		m.LastAction = "research:" + pick.Coin
		m.Stage = "researching"
		_ = auto.SaveMission(h.Dir, m)
		_ = auto.Save(h.Dir, p)
		h.beginResearch(pick.Coin)
		return
	}
	if wantResearch && p.LastResearchCoin == pick.Coin {
		m.LastAction = "waiting_after_research:" + pick.Coin
		if m.BlockReason != "" {
			m.Stage = "execution-blocked"
		} else {
			m.Stage = "eligible"
		}
	}
	_ = auto.SaveMission(h.Dir, m)
	_ = auto.Save(h.Dir, p)
}

func (h *Hub) syncMissionResearch(m *auto.Mission) {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	if !h.job.running {
		return
	}
	m.Stage = "researching"
	m.LastAction = "research:" + h.job.coin
	if h.job.coin != "" {
		m.BestCoin = h.job.coin
	}
}

func (h *Hub) positionSummary() string {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return ""
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return ""
	}
	rows, _, perr := hl.New(config.For(net)).Clearinghouse(st.Wallet)
	if perr != nil {
		return ""
	}
	parts := make([]string, 0, len(rows))
	for _, row := range rows {
		sz := strings.TrimSpace(row.Sz)
		if sz == "" || sz == "0" {
			continue
		}
		parts = append(parts, strings.ToUpper(row.Coin)+" "+sz)
	}
	return strings.Join(parts, ", ")
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
	pol := cli.ActivePolicy(h.Dir)
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
		auto.RecordBlock(h.Dir, err.Error(), coin)
		auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "rejected", Coin: coin, Why: err.Error()})
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.refused", Market: coin,
			Action: "guarded", Status: err.Error(), PreviewHash: hash, Reason: err.Error(), Autonomous: true,
		})
		return
	}
	got := cli.ExecuteDeskOrder(h.Dir, cli.ConfirmToken, hash)
	if got.Error != "" {
		auto.RecordStage(h.Dir, "exec_failed", "exec_failed:"+got.Error, got.Error, coin)
		auto.RecordAction(h.Dir, "exec_failed:"+got.Error, coin, hash, "", "")
		auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "rejected", Coin: coin, Why: got.Error})
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "order.rejected", Market: coin,
			Action: "guarded", Status: got.Error, PreviewHash: hash, Reason: got.Error, Autonomous: true,
		})
		return
	}
	auto.RecordAction(h.Dir, "executed", coin, hash, got.OID, "")
	auto.RecordStage(h.Dir, "executed", "executed", "oid:"+got.OID, coin)
	auto.AppendAway(h.Dir, auto.AwayEvent{Kind: "traded", Coin: coin, Why: "guarded_autonomy", OID: got.OID})
	link := venueTradeLink(workspaceNetwork(h.Dir), got.Market)
	appendActivity(h.Dir, activityEvent{
		WorkspaceID: workspaceID(h.Dir), Kind: "order.submitted", Market: got.Market,
		Action: "guarded", Status: got.Status, OID: got.OID, PreviewHash: got.Hash,
		Reason: "guarded_autonomy", Link: link, Autonomous: true,
	})
	h.recordPostedOrder(got, "guarded", h.currentJobID())
	h.fileOrder(got, h.currentJobID())
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
		writeLocal(w, http.StatusOK, h.missionPublic())
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
		Stop            bool     `json:"stop"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Stop || strings.EqualFold(strings.TrimSpace(body.Typed), auto.StopToken) {
		m := auto.Stop(h.Dir, "user_stop")
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.stopped", Action: "stop", Status: "user_stop", Reason: "user_stop", Autonomous: true,
		})
		out := h.missionPublic()
		out["mission"] = m
		writeLocal(w, http.StatusOK, out)
		return
	}
	if strings.TrimSpace(body.Typed) != "" || strings.EqualFold(body.Mode, auto.ModeGuarded) {
		if st, lerr := cli.Load(h.Dir); lerr == nil && strings.TrimSpace(st.Wallet) != "" {
			if err := cli.CheckPinned(h.Dir, st.WorkspaceID, cli.ActivePolicy(h.Dir)); err != nil {
				writeLocal(w, http.StatusOK, map[string]any{
					"ok": false, "error": "need_pin", "explain": auto.HumanWhy("need_pin"), "execute": false, "sign": false, "trade": false,
				})
				return
			}
		}
		hash, _ := cli.ActivePolicy(h.Dir).Hash()
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
			WorkspaceID: workspaceID(h.Dir), Kind: "mission.enabled", Action: "guarded", Status: "running", Reason: "ENABLE GUARDED AUTONOMY", Autonomous: true,
		})
		go h.autoTick()
		writeLocal(w, http.StatusOK, h.missionPublic())
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
		writeLocal(w, http.StatusOK, h.missionPublic())
		return
	}
	m := auto.LoadMission(h.Dir)
	m.MaxTrades = body.MaxTrades
	m.StopLossUSD = body.StopLossUSD
	m.MinLiquidityUSD = body.MinLiquidityUSD
	m.PauseUncertain = body.PauseUncertain
	m.Assets = body.Assets
	_ = auto.SaveMission(h.Dir, m)
	writeLocal(w, http.StatusOK, h.missionPublic())
}

func (h *Hub) missionPublic() map[string]any {
	out := auto.Public(h.Dir)
	h.researchMu.Lock()
	out["research_running"] = h.job.running
	out["research_stage"] = h.job.stage
	out["research_coin"] = h.job.coin
	out["research_job_id"] = h.job.ID
	h.researchMu.Unlock()
	if last := cli.LoadLastOrder(h.Dir); last != nil {
		out["last_order"] = last
	}
	return out
}
