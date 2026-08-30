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
	"github.com/mohamedwael201193/pit/internal/policy"
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
		if parsed.Coin == "" {
			parsed.Coin = h.pickBestCoin()
		}
	case "watch.why_not":
		parsed.Reply = h.replyWhyNotTrade()
		parsed.Navigate = "automation"
	case "experience.why":
		coin := parsed.Coin
		if coin == "" {
			coin = h.pickBestCoin()
		}
		parsed.Reply = h.whyThisSetup(coin) + " Chat cannot AUTHORIZE. Private memory never includes the memory key."
		parsed.Navigate = "research"
	case "research.start", "research.best":
		h.researchMu.Lock()
		running := h.job.running
		jobCoin := h.job.coin
		jobID := h.job.ID
		stage := h.job.stage
		h.researchMu.Unlock()
		if running {
			parsed.StartResearch = false
			parsed.Mutate = false
			parsed.Navigate = ""
			parsed.Tool = "research.status"
			parsed.Coin = jobCoin
			parsed.Reply = fmt.Sprintf("Already researching %s (job %s, %s). Chat cannot AUTHORIZE. Stay on Chat for live stages, or open Research.", jobCoin, jobID, stage)
			break
		}
		if live := h.replyWatch(parsed); live != "" {
			parsed.Reply = live + " " + parsed.Reply
		}
	case "policy.get":
		parsed.Reply = h.replyPolicy()
		parsed.Navigate = "security"
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
		mode = "Sleep Mission"
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
	cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), cli.ActivePolicy(h.Dir))
	if lerr != nil {
		return "Hyperliquid did not return a live book. Empty Markets is the honest state. No invented scores."
	}
	view := watch.Public(cands, string(net))
	view = h.annotateWatch(view, cli.ActivePolicy(h.Dir))
	if parsed.Tool == "watch.scan" {
		return fmt.Sprintf("Scanned %d live Hyperliquid perps. %d match your policy. %d execution-feasible for this account. Buying power $%.2f (%s). %s Side is not decided here.", view.Scanned, view.Count, view.ExecFeasibleN, view.BuyingPower, view.PowerSource, view.Copy)
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
		whyBetter = "\nWhy it ranks first: " + view.BestWhy
	}
	execLine := row.WhyExecutable
	if execLine == "" {
		execLine = row.ExecWhy
	}
	return fmt.Sprintf("%s\nMark %s (oracle %s, funding %s, open interest %s, day notional %s).\n%s\nTrend: %s. Host rank %d, not a model score. Policy %s. Layer %s. Freshness %s. Venue hyperliquid. Provenance %s.\nWhy it passes policy: %s\nRequired margin $%.2f · available $%.2f · host-sized $%.2f · min $%.2f.\nWhy it is executable for this user: %s\nWhat would invalidate it: %s\nWhat research will test: %s%s\nSide is not decided here. Chat cannot AUTHORIZE.",
		row.Coin, watch.Price(row.Mark), watch.Price(row.Oracle), watch.FundingPct(row.Funding), watch.Compact(row.OpenInterest), watch.CompactUSD(row.Volume), row.Why, row.Trend, row.Rank, fit, row.Layer, row.Freshness, row.Provenance,
		watch.WhyPolicyFrom(row.Eligible, row.Block, h.policyPinnedNow()), row.RequiredMargin, row.AvailableMargin, row.HostNotional, row.MinNotional, execLine, watch.WhatInvalidatesReason(row.Reason), watch.ResearchWillTestFrom(row.Coin, row.Eligible), whyBetter)
}

func (h *Hub) pickBestCoin() string {
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	if err == nil && strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		return ""
	}
	pol := cli.ActivePolicy(h.Dir)
	cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), pol)
	if lerr != nil {
		return ""
	}
	skip := auto.Load(h.Dir).SkipSet(time.Now().Unix())
	if best, _, ok := watch.NextCandidate(cands, h.capitalNow(), pol, h.sessionAliveNow(), h.policyPinnedNow(), skip); ok {
		return best.Coin
	}
	return ""
}

func formatBookFloors(bp float64, coins []watch.PublicCoin) string {
	parts := make([]string, 0, 8)
	n := 0
	for _, c := range coins {
		if !c.PolicyEligible && !c.Eligible {
			continue
		}
		n++
		min := c.MinNotional
		if min <= 0 {
			parts = append(parts, c.Coin+" min unknown")
			continue
		}
		if bp+1e-9 < min {
			parts = append(parts, fmt.Sprintf("%s $%.2f short of $%.2f", c.Coin, min-bp, min))
			continue
		}
		if c.ExecutionFeasible {
			parts = append(parts, fmt.Sprintf("%s executable at min $%.2f", c.Coin, min))
			continue
		}
		if c.ExecGate == "policy_clip_tight" || (bp+1e-9 >= min && c.PolicyClip > 0 && c.PolicyClip+1e-9 < min) {
			gap := min - c.PolicyClip
			if gap < 0 {
				gap = 0
			}
			parts = append(parts, fmt.Sprintf("%s policy cap $%.2f too tight for min $%.2f (account $%.2f)", c.Coin, gap, min, bp))
			continue
		}
		parts = append(parts, fmt.Sprintf("%s min $%.2f", c.Coin, min))
	}
	if n == 0 {
		return "No policy-eligible books in this Watch."
	}
	return fmt.Sprintf("%d pass policy. %s", n, strings.Join(parts, "; "))
}

func (h *Hub) replyWhyNotTrade() string {
	cap := h.capitalNow()
	pol := cli.ActivePolicy(h.Dir)
	away := auto.LoadAway(h.Dir)
	p := auto.Load(h.Dir)
	floors := "Watch did not return sized books on this computer."
	nextLine := "No remaining candidate qualifies under this account and law."
	bestCoin := ""
	st, err := cli.Load(h.Dir)
	netName := "mainnet"
	if err == nil && strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	if net, nerr := config.ParseNetwork(netName); nerr == nil {
		cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), pol)
		if lerr == nil {
			view := watch.Public(cands, string(net))
			view = h.annotateWatch(view, pol)
			floors = formatBookFloors(cap.BuyingPower, view.Coins)
			skip := p.SkipSet(time.Now().Unix())
			best, exe, ok := watch.NextCandidate(cands, cap, pol, h.sessionAliveNow(), h.policyPinnedNow(), skip)
			if ok {
				bestCoin = best.Coin
				why := ""
				for _, c := range view.Coins {
					if strings.EqualFold(c.Coin, best.Coin) {
						why = c.WhyExecutable
						if why == "" {
							why = c.ExecWhy
						}
						break
					}
				}
				if why == "" {
					why = "This candidate is not executable for this account."
				}
				layer := "not executable"
				if exe {
					layer = "execution-feasible"
				}
				nextLine = fmt.Sprintf("Next candidate %s is %s. %s", best.Coin, layer, why)
			}
		}
	}
	lesson := h.whyThisSetup(bestCoin)
	return fmt.Sprintf("Why PIT did not trade: Buying power $%.2f (%s). %s Scan continues past a blocked or stood-down book. %s While you were away: %d detected, %d researched, %d refused, %d autonomous trades. %s Chat cannot AUTHORIZE. PIT will not invent size.",
		cap.BuyingPower, cap.PowerSource, floors, nextLine, away.Detected, away.Researched, away.Rejected, away.Traded, lesson)
}

func (h *Hub) replyMissionEnable() string {
	return "Chat cannot arm a Sleep Mission. Open Automation, review the host limits, then confirm ARM SLEEP MISSION on this computer. ENABLE GUARDED AUTONOMY is still accepted. The model cannot change those limits."
}

func (h *Hub) replyMissionStop() string {
	m := auto.LoadMission(h.Dir)
	why := m.LastStop
	if why == "" {
		why = "stopped"
	}
	return "Sleep Mission is stopped (" + why + "). PIT will not place further autonomous orders until you arm it again on Automation. Positions were not flattened. Chat cannot AUTHORIZE."
}

func (h *Hub) replyMissionStatus() string {
	m := auto.LoadMission(h.Dir)
	log := auto.LoadEvents(h.Dir)
	return fmt.Sprintf("Sleep state %s. Mode %s. Running %v. Mission %s. Detected %d, researched %d, challenger rejects %d, executions %d, fills %d. Chat cannot arm a Sleep Mission. Confirm ARM SLEEP MISSION or ENABLE GUARDED AUTONOMY on Automation.", m.SleepState, m.Mode, m.Running, m.MissionID, log.Detected, log.Researched, log.Challenger, log.Executions, log.Fills)
}

func (h *Hub) replyPolicy() string {
	p := policy.Peek(h.Dir)
	open := h.openPositionCount()
	avail := h.availableUSD()
	cap := h.capitalNow()
	block, why := policy.ExecWhy(open, avail, p)
	gate := "Execution gate is clear."
	if block != "" {
		gate = "Execution blocked: " + strings.ReplaceAll(block, "_", " ") + ". " + why
	}
	pin := "Policy is drafted, not pinned. Markets PASS is the default universe — not locked host law. Pin on Security. Chat cannot pin."
	if h.policyPinnedNow() {
		pin = "Host law is pinned on this computer."
	}
	spot := ""
	if cap.SpotUSDC > 0 {
		spot = fmt.Sprintf(" Spot USDC %.4f (%s).", cap.SpotUSDC, cap.PowerSource)
	}
	return fmt.Sprintf("%s Host policy on this computer: max trade $%.0f, leverage 1x, assets %s, venue hyperliquid, max open %d, daily loss $%.0f, slippage %d bps, cooldown %ds, uncertainty %.2f, session TTL %ds. Chat cannot pin or mutate this. Open Security to edit and pin. %s%s",
		pin, p.MaxClipUSD, strings.Join(p.AllowedAssets, " "), p.MaxOpenPositions, p.DailyLossUSD, p.MaxSlippageBps, p.CooldownSeconds, p.MaxUncertainty, p.SessionTTLSeconds, gate, spot)
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
		return "Hyperliquid did not return clearinghouse state. Open Portfolio when the venue is reachable. PIT cannot withdraw."
	}
	av, _ := hl.New(config.For(net)).Account(st.Wallet)
	spotNote := ""
	if av.SpotUSDC > 0 {
		spotNote = fmt.Sprintf(" Spot USDC %.4f.", av.SpotUSDC)
	}
	if len(rows) == 0 {
		last := cli.LoadLastOrder(h.Dir)
		oid := ""
		if s, ok := last["oid"].(string); ok && s != "" {
			oid = " Last fill OID " + s + " is historical, not a new preview."
		}
		cap := h.capitalNow()
		return fmt.Sprintf("No open positions on the trading account. Trading equity $%.2f (%s). Perp account value %s.%s%s", cap.BuyingPower, cap.PowerSource, acct.AccountValue, spotNote, oid)
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
		kind := TerminalKind(false, h.job.err, h.job.deny, namedRolesVerified(h.job.roles), h.job.eligible, h.job.roles)
		if kind == TermReadyStoodDown {
			return fmt.Sprintf("Committee stood down (%s). That is a verified no-trade, not a crash. %s PIT will check the next eligible market. No order was placed. Chat cannot AUTHORIZE.", h.job.deny, roles)
		}
		if kind == TermMarketDenied {
			return fmt.Sprintf("This market is not sizeable (%s). %s PIT will not invent size. Chat cannot AUTHORIZE.", h.job.deny, roles)
		}
		return fmt.Sprintf("Committee result: %s (%s). %s One verified role is never a committee result. No order was placed. Open Research.", kind, h.job.deny, roles)
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
