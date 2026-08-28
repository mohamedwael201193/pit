package companion

import (
	"fmt"
	"strings"
	"time"

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
	case "session.status":
		parsed.Reply = h.replySession()
	case "watch.get":
		parsed.Reply = h.replyWatch(parsed)
	case "research.start":
		if live := h.replyWatch(parsed); live != "" {
			parsed.Reply = live + " " + parsed.Reply
		}
	case "policy.get":
		parsed.Reply = "Policy is pinned host law on this computer. Chat cannot raise clip, leverage, or permissions. Open Policy to read the exact constraints."
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
	if eligible && hash != "" {
		return fmt.Sprintf("Idle. Wallet %s. %s. Exact preview is waiting for AUTHORIZE on Research. Chat cannot AUTHORIZE.", wallet, sess)
	}
	return fmt.Sprintf("Idle. Wallet %s. %s. Open Watch to discover, Research to investigate. Chat cannot AUTHORIZE.", wallet, sess)
}

func (h *Hub) replyWatch(parsed deskcmd.Result) string {
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		if parsed.Coin != "" {
			return parsed.Coin + " is in the policy universe. Bind this computer to load the live Hyperliquid mark. Side is not decided here."
		}
		return "Watch lists policy-eligible public markets after this computer is bound. Empty is honest until then."
	}
	netName := "mainnet"
	if strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		return "Wrong network. Watch stays unread."
	}
	cands, lerr := watch.Live(hl.New(config.For(net)), watch.PolicyForWatch())
	if lerr != nil {
		return "Hyperliquid did not return a live book. Empty Watch is the honest state. No invented scores."
	}
	view := watch.Public(cands, string(net))
	if len(view.Coins) == 0 {
		return "No opportunities match your policy. Empty is honest. Side is not decided on Watch."
	}
	row := view.Coins[0]
	if parsed.Coin != "" {
		for _, c := range view.Coins {
			if c.Coin == parsed.Coin {
				row = c
				break
			}
		}
	}
	fit := "PASS"
	if !row.Eligible {
		fit = "BLOCKED"
	}
	return fmt.Sprintf("%s mark %g (oracle %g, funding %g, open interest %.0f). %s Trend: %s. Host rank %d, not a model score. Policy %s. Freshness %s. Side is not decided here.",
		row.Coin, row.Mark, row.Oracle, row.Funding, row.OpenInterest, row.Why, row.Trend, row.Rank, fit, row.Freshness)
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
	if h.job.deny != "" {
		return fmt.Sprintf("Committee stood down (%s). Three-role verification can still be OK. No order was placed. Open Research for the named reason.", h.job.deny)
	}
	if h.job.eligible && h.job.previewHash != "" {
		return "Exact preview is on Research. Chat cannot AUTHORIZE it. Type AUTHORIZE only on that card."
	}
	if h.job.err != "" {
		return "Last research ended: " + h.job.err + ". That is not a fake success. Open Research."
	}
	return fallback
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
