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
)

func (h *Hub) decorateChat(parsed deskcmd.Result) deskcmd.Result {
	switch parsed.Tool {
	case "status":
		parsed.Reply = h.replyDeskStatus()
	case "positions.get":
		parsed.Reply = h.replyPositions()
	case "research.result", "preview.show":
		parsed.Reply = h.replyResearch(parsed.Reply)
	case "session.status":
		parsed.Reply = h.replySession()
	}
	return parsed
}

func (h *Hub) replyDeskStatus() string {
	h.researchMu.Lock()
	running := h.job.running
	coin := h.job.coin
	stage := h.job.stage
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
	return fmt.Sprintf("Idle. Wallet %s. %s. Open Desk for what needs you. Chat cannot AUTHORIZE.", wallet, sess)
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
