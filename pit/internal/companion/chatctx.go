package companion

import (
	"fmt"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskcmd"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (h *Hub) decorateChat(parsed deskcmd.Result) deskcmd.Result {
	switch parsed.Tool {
	case "status":
		parsed.Reply = h.replyDeskStatus()
	case "greet":
		parsed.Reply = "Hello. " + h.replyDeskStatus()
	case "positions.get":
		parsed.Reply = h.replyPositions()
	case "research.result", "preview.show":
		parsed.Reply = h.replyResearch(parsed.Reply)
	case "refuse_execute":
		parsed.Reply = h.replyCannotExecute(parsed.Reply)
	case "activity.list", "activity.today", "activity.proof":
		parsed.Reply = h.replyActivity(parsed.Reply)
	case "session.status":
		parsed.Reply = h.replySession()
	case "watch.get", "watch.best", "watch.scan", "watch.compare":
		parsed.Reply = h.replyWatch(parsed)
	case "research.start", "research.best":
		if live := h.replyWatch(parsed); live != "" {
			parsed.Reply = live + " " + parsed.Reply
		}
	case "policy.get":
		parsed.Reply = "Policy is pinned host law on this computer. Chat cannot raise clip, leverage, or permissions. Open Security to read the exact constraints."
	case "setup.guide":
		parsed.Reply = h.replyDeskStatus() + " First-run setup walks wallet, network, Hyperliquid, session, Protect, private compute, then policy. Chat cannot AUTHORIZE."
	case "mission.enable_required":
		parsed.Reply = h.replyMissionEnable()
	case "mission.stop":
		parsed.Reply = h.replyMissionStop()
	case "mission.status":
		parsed.Reply = h.replyMissionStatus()
	}
	return parsed
}

func (h *Hub) replyDeskStatus() string {
	h.researchMu.Lock()
	running := h.job.running
	coin := h.job.coin
	stage := h.job.stage
	eligible := h.job.eligible
	hash := h.job.previewHash
	h.researchMu.Unlock()
	st, err := cli.Load(h.Dir)
	wallet := "unbound"
	if err == nil && strings.TrimSpace(st.Wallet) != "" {
		wallet = "bound"
	}
	sf, serr := cli.LoadSession(h.Dir)
	sess := "no session"
	if serr == nil && strings.TrimSpace(sf.AgentAddr) != "" {
		if session.Alive(sf.Meta().Session(), time.Now().UnixMilli()) {
			sess = "session live"
		} else {
			sess = "session expired"
		}
	}
	if running {
		return fmt.Sprintf("Researching %s — %s still running. Compute money, not trading capital. Chat cannot AUTHORIZE. Open Research for the live board.", strings.ToUpper(coin), strings.ReplaceAll(stage, "_", " "))
	}
	m := auto.LoadMission(h.Dir)
	mode := "Manual"
	if m.Mode == auto.ModeGuarded && m.Running {
		mode = "Guarded Autonomy"
	} else if m.Mode == auto.ModeResearch {
		mode = "Research Only"
	}
	if eligible && hash != "" {
		return fmt.Sprintf("Idle. Wallet %s. %s. Mode %s. Exact preview is waiting for AUTHORIZE on Research. Chat cannot AUTHORIZE.", wallet, sess, mode)
	}
	top := h.replyWatch(deskcmd.Result{})
	if strings.Contains(top, "mark") {
		return fmt.Sprintf("Idle. Wallet %s. %s. Mode %s. Best opportunity: %s Chat cannot AUTHORIZE.", wallet, sess, mode, top)
	}
	return fmt.Sprintf("Idle. Wallet %s. %s. Mode %s. Open Markets to discover, Research to investigate. Chat cannot AUTHORIZE.", wallet, sess, mode)
}

func (h *Hub) replyWatch(parsed deskcmd.Result) string {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		if parsed.Coin != "" {
			return parsed.Coin + " is in the policy universe. Bind this computer to load the live Hyperliquid mark. Side is not decided here."
		}
		return "Markets lists live Hyperliquid books after this computer is bound. Empty is honest until then."
	}
	netName := "mainnet"
	if strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		return "Wrong network. Markets stay unread."
	}
	cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), watch.PolicyForWatch())
	if lerr != nil {
		return "Hyperliquid did not return a live book. Empty Markets is the honest state. No invented scores."
	}
	view := watch.Public(cands, string(net))
	if parsed.Tool == "watch.scan" {
		return fmt.Sprintf("Scanned %d live Hyperliquid perps. %d match your policy. %s Side is not decided here.", view.Scanned, view.Count, view.Copy)
	}
	if len(view.Coins) == 0 {
		return "No live Hyperliquid books matched this scan. Empty is honest. Side is not decided on Markets."
	}
	row := view.Coins[0]
	if view.Best != nil && (parsed.Tool == "watch.best" || parsed.Tool == "watch.compare" || parsed.Coin == "") {
		row = *view.Best
	}
	if parsed.Coin != "" {
		for _, c := range view.Coins {
			if c.Coin == parsed.Coin {
				row = c
				break
			}
		}
	}
	fit := row.PolicyFit
	if fit == "" {
		if row.Eligible {
			fit = "PASS"
		} else {
			fit = "BLOCKED"
		}
	}
	whyBetter := ""
	if parsed.Tool == "watch.compare" && view.Best != nil {
		whyBetter = " " + view.BestWhy
	}
	return fmt.Sprintf("%s mark %g (oracle %g, funding %g, open interest %.0f, day notional %g). %s Trend: %s. Host rank %d, not a model score. Policy %s. Freshness %s. Venue hyperliquid. Provenance %s.%s Side is not decided here.",
		row.Coin, row.Mark, row.Oracle, row.Funding, row.OpenInterest, row.Volume, row.Why, row.Trend, row.Rank, fit, row.Freshness, row.Provenance, whyBetter)
}

func (h *Hub) pickBestCoin() string {
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	if err == nil && strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		return "ETH"
	}
	cands, lerr := watch.Live(hl.New(config.For(net)), watch.PolicyForWatch())
	if lerr != nil {
		return "ETH"
	}
	if best, ok := watch.Best(cands); ok {
		return best.Coin
	}
	return "ETH"
}

func (h *Hub) replyMissionEnable() string {
	return "Chat cannot enable Guarded Autonomy. Open Automation, review the host limits, then type ENABLE GUARDED AUTONOMY. The model cannot change those limits."
}

func (h *Hub) replyMissionStop() string {
	m := auto.LoadMission(h.Dir)
	why := m.LastStop
	if why == "" {
		why = "stopped"
	}
	return "Guarded Autonomy is stopped (" + why + "). PIT will not place further orders until you enable it again on Automation. Chat cannot AUTHORIZE."
}

func (h *Hub) replyMissionStatus() string {
	m := auto.LoadMission(h.Dir)
	return fmt.Sprintf("Mode %s. Running %v. Best %s. Last action %s. Stop reason %s. Trades today %d. Chat cannot enable Guarded Autonomy.", m.Mode, m.Running, m.BestCoin, m.LastAction, m.LastStop, m.TradesToday)
}

func (h *Hub) replyPositions() string {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		return "This computer is not bound. Positions are read from your Hyperliquid trading account after you connect a wallet."
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		return "Wrong network. Positions stay unread."
	}
	rows, acct, perr := hl.New(config.For(net)).Clearinghouse(st.Wallet)
	if perr != nil {
		return "Hyperliquid did not return clearinghouse state. Open Positions when the venue is reachable. PIT cannot withdraw."
	}
	if len(rows) == 0 {
		last := cli.LoadLastOrder(h.Dir)
		oid := ""
		if s, ok := last["oid"].(string); ok && s != "" {
			oid = " Last fill OID " + s + " is historical, not a new preview."
		}
		return "No open positions on the trading account. Equity " + acct.AccountValue + "." + oid
	}
	parts := make([]string, 0, len(rows))
	for _, p := range rows {
		parts = append(parts, fmt.Sprintf("%s %s @ %s uPnL %s", p.Coin, p.Sz, p.EntryPx, p.UnrealizedPnl))
	}
	return "Venue positions (master account, not the PIT agent): " + strings.Join(parts, "; ") + ". PIT cannot withdraw. Chat cannot flatten."
}

func (h *Hub) replyResearch(fallback string) string {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	if h.job.running {
		return fmt.Sprintf("Research is still running on %s. Live view can lag. That is not a TEE failure. Open Research.", strings.ToUpper(h.job.coin))
	}
	roles := roleLine(h.job.roles)
	if h.job.deny != "" {
		return fmt.Sprintf("Committee stood down (%s). %s One verified role is never a committee result. No order was placed. Open Research for the named reason.", h.job.deny, roles)
	}
	if h.job.eligible && h.job.previewHash != "" {
		return "Exact preview is on Research. " + roles + " Chat cannot AUTHORIZE it. Type AUTHORIZE only on that card."
	}
	if h.job.err != "" {
		return "Last research ended: " + h.job.err + ". That is not a fake success. Open Research."
	}
	if roles != "" {
		return fallback + " " + roles
	}
	return fallback
}

func (h *Hub) replyCannotExecute(fallback string) string {
	h.researchMu.Lock()
	eligible := h.job.eligible
	hash := h.job.previewHash
	running := h.job.running
	h.researchMu.Unlock()
	sf, serr := cli.LoadSession(h.Dir)
	live := serr == nil && session.Alive(sf.Meta().Session(), time.Now().UnixMilli())
	if running {
		return "Research is still running. Chat cannot AUTHORIZE. Wait for the exact preview on Research."
	}
	if eligible && hash != "" {
		if !live {
			return "An exact preview is waiting, but there is no live order/cancel session. Create a session on Security, then type AUTHORIZE on Research. Chat cannot AUTHORIZE."
		}
		return "An exact preview is waiting on Research. Chat cannot AUTHORIZE. Type AUTHORIZE only on that card. Withdraw stays impossible."
	}
	return fallback
}

func (h *Hub) replyActivity(fallback string) string {
	evs := readActivity(h.Dir, 8)
	if len(evs) == 0 {
		return fallback
	}
	last := evs[len(evs)-1]
	return fmt.Sprintf("Last desk event: %s %s %s. Historical fills never appear inside a new preview. Open Activity.", last.Kind, last.Market, last.Status)
}

func roleLine(roles []map[string]any) string {
	if len(roles) == 0 {
		return ""
	}
	parts := make([]string, 0, len(roles))
	ok := 0
	for _, rm := range roles {
		name := strings.TrimSpace(fmtString(rm["role"]))
		if name == "" {
			continue
		}
		side := strings.TrimSpace(fmtString(rm["proposed_side"]))
		ver := strings.TrimSpace(fmtString(rm["verify_e2ee"]))
		bit := name
		if strings.EqualFold(ver, "OK") {
			ok++
			bit += " verified"
		}
		if side != "" {
			bit += " " + side
		}
		if kill, _ := rm["kill"].(bool); kill {
			bit += " stood down"
		}
		parts = append(parts, bit)
	}
	if len(parts) == 0 {
		return ""
	}
	note := strings.Join(parts, "; ") + "."
	if ok > 0 && ok < 3 {
		note += " Incomplete roles are not a committee result."
	}
	return note
}

func (h *Hub) replySession() string {
	sf, err := cli.LoadSession(h.Dir)
	if err != nil || strings.TrimSpace(sf.AgentAddr) == "" {
		return "No local trading session. Create a secure session on this computer, then approve that PIT agent on Hyperliquid API. Chat cannot approve it."
	}
	name := strings.TrimSpace(sf.Workspace)
	if n, nerr := session.AgentName(sf.Workspace); nerr == nil && n != "" {
		name = n
	}
	if session.Alive(sf.Meta().Session(), time.Now().UnixMilli()) {
		return name + " is live on this computer. Order and cancel only. Withdraw, transfer, and leverage stay denied. Open Security to check Hyperliquid listing."
	}
	return name + " exists but the local session is not live. Create a secure session. Hyperliquid may already list this agent — PIT will reuse it."
}
