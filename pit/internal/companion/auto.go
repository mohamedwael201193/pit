package companion

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/httpx"
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
	if !p.Watch {
		return
	}
	now := time.Now().Unix()
	if p.LastScanUnix == 0 {
		p.LastScanUnix = now
		_ = auto.Save(h.Dir, p)
		return
	}
	if now-p.LastScanUnix < int64(p.CadenceMinutes)*60 {
		return
	}
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	if err == nil && strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		return
	}
	cands, lerr := watch.Live(hl.New(config.For(net)), watch.PolicyForWatch())
	p.LastScanUnix = now
	if lerr != nil || len(cands) == 0 {
		_ = auto.Save(h.Dir, p)
		return
	}
	filtered := cands
	if len(p.Markets) > 0 {
		allow := map[string]bool{}
		for _, m := range p.Markets {
			allow[strings.ToUpper(strings.TrimSpace(m))] = true
		}
		next := make([]watch.Candidate, 0, len(cands))
		for _, c := range cands {
			if allow[c.Coin] {
				next = append(next, c)
			}
		}
		filtered = next
	}
	var pick *watch.Candidate
	for i := range filtered {
		if auto.Matches(p.Trigger, filtered[i].Reason) {
			pick = &filtered[i]
			break
		}
	}
	if pick == nil {
		_ = auto.Save(h.Dir, p)
		return
	}
	if p.Notify && p.LastNotifyCoin != pick.Coin {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "opportunity", Market: pick.Coin,
			Action: "discover", Status: "ready", Reason: pick.Reason,
		})
		p.LastNotifyCoin = pick.Coin
	}
	if p.AutoResearch && p.LastResearchCoin != pick.Coin {
		h.researchMu.Lock()
		running := h.job.running
		h.researchMu.Unlock()
		if !running {
			p.LastResearchCoin = pick.Coin
			appendActivity(h.Dir, activityEvent{
				WorkspaceID: workspaceID(h.Dir), Kind: "automation.prepared", Market: pick.Coin,
				Action: "prepare", Status: "research", Reason: "automation_stops_at_human_approval",
			})
			_ = auto.Save(h.Dir, p)
			h.beginResearch(pick.Coin)
			return
		}
	}
	_ = auto.Save(h.Dir, p)
}

func (h *Hub) localAutomation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
			http.Error(w, "origin_denied", http.StatusForbidden)
			return
		}
		p := auto.Load(h.Dir)
		p.Execute = false
		writeLocal(w, http.StatusOK, map[string]any{
			"ok": true, "prefs": p, "execute": false, "sign": false, "trade": false,
			"note": "Automation may watch, discover, research, notify, and prepare. It cannot AUTHORIZE.",
		})
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
