package deskcmd

import (
	"fmt"
	"net/url"
	"strings"
)

type Result struct {
	Reply         string `json:"reply"`
	Tool          string `json:"tool"`
	Mutate        bool   `json:"mutate"`
	Execute       bool   `json:"execute"`
	StartResearch bool   `json:"start_research"`
	Coin          string `json:"coin,omitempty"`
	Navigate      string `json:"navigate,omitempty"`
	OpenURL       string `json:"open_url,omitempty"`
	Sign          bool   `json:"sign"`
	Trade         bool   `json:"trade"`
}

var allowedCoins = map[string]bool{
	"ETH": true, "BTC": true, "SOL": true, "HYPE": true, "DOGE": true, "AVAX": true,
}

func Parse(text string) Result {
	raw := strings.TrimSpace(text)
	low := strings.ToLower(raw)
	out := Result{Sign: false, Trade: false, Execute: false, Mutate: false}
	if raw == "" {
		out.Reply = "Ask about a market, research, policy, session, evidence, or an official page. PIT will not execute from chat."
		out.Tool = "help"
		return out
	}
	if wantsExecute(low) || wantsFlatten(low) {
		out.Tool = "refuse_execute"
		out.Navigate = "preview"
		if wantsFlatten(low) {
			out.Reply = "PIT will not flatten from chat. Prepare a reduce-only close on this computer, then type AUTHORIZE there."
			return out
		}
		out.Reply = "PIT will not buy, sell, or authorize from chat. Review the exact preview on this computer, then type AUTHORIZE there."
		return out
	}
	if strings.Contains(low, "forget") && (strings.Contains(low, "memory") || strings.Contains(low, "everything")) {
		out.Tool = "memory.forget"
		out.Mutate = true
		out.Reply = "Forget wipes workspace working memory and lessons. It never deletes venue positions or receipts."
		return out
	}
	if strings.Contains(low, "happening") || strings.Contains(low, "what is pit doing") {
		out.Tool = "status"
		out.Navigate = "home"
		out.Reply = "Desk shows live status. Chat cannot see the sealed book. Research, preview, and session stay on this computer."
		return out
	}
	if strings.Contains(low, "learn") || strings.Contains(low, "calibration") {
		out.Tool = "calibration.get"
		out.Navigate = "settings"
		out.Reply = "Calibration stays NOT ENOUGH DATA until enough real outcomes exist. PIT will not invent a lesson."
		return out
	}
	if strings.Contains(low, "hyperliquid") && (strings.Contains(low, "ready") || strings.Contains(low, "approved")) {
		out.Tool = "session.status"
		out.Navigate = "security"
		out.Reply = "Hyperliquid is ready only after this computer has a session and the PIT agent is listed. Open Security. Chat cannot approve the agent."
		return out
	}
	if strings.Contains(low, "open pairing") || strings.Contains(low, "pair this") {
		out.Tool = "links.open"
		out.OpenURL = "https://pit0g.vercel.app/pair"
		out.Reply = "Opening pairing. The website never receives a session key."
		return out
	}
	if strings.Contains(low, "open hyperliquid api") {
		out.Tool = "links.open"
		out.OpenURL = "https://app.hyperliquid.xyz/API"
		out.Navigate = "security"
		out.Reply = "Opening the official Hyperliquid API page. Approve the PIT agent shown on this computer. PIT cannot withdraw."
		return out
	}
	if strings.Contains(low, "open hyperliquid") {
		out.Tool = "links.open"
		out.OpenURL = "https://app.hyperliquid.xyz"
		out.Reply = "Opening the official Hyperliquid app for this network."
		return out
	}
	if strings.Contains(low, "private compute") || strings.Contains(low, "pc.0g") || strings.Contains(low, "0g private") {
		out.Tool = "links.open"
		out.OpenURL = "https://pc.0g.ai/sdk/dashboard/funds"
		out.Reply = "Opening 0G Private Compute (Advanced funds). That is compute money, not Hyperliquid trading capital."
		return out
	}
	if strings.Contains(low, "session") && (strings.Contains(low, "permission") || strings.Contains(low, "what is")) {
		out.Tool = "session.status"
		out.Reply = "Order and cancel only. Withdraw, transfer, and leverage stay denied. Chat cannot widen that."
		return out
	}
	if strings.Contains(low, "evidence") || strings.Contains(low, "tee") || strings.Contains(low, "proof") {
		out.Tool = "research.result"
		out.Navigate = "research"
		out.Reply = "Opening research evidence on this computer. The sealed prompt stays private."
		return out
	}
	if strings.Contains(low, "preview") {
		out.Tool = "preview.show"
		out.Navigate = "research"
		out.Reply = "The exact preview is on the desk. Chat cannot AUTHORIZE it."
		return out
	}
	if strings.Contains(low, "changed") || strings.Contains(low, "recent") || strings.Contains(low, "activity") || strings.Contains(low, "yesterday") {
		out.Tool = "activity.list"
		out.Navigate = "activity"
		out.Reply = "Recent desk events are on Activity. Historical fills never appear inside a new preview."
		return out
	}
	if strings.Contains(low, "opportunit") || strings.Contains(low, "match my policy") || strings.Contains(low, "watch") {
		out.Tool = "watch.get"
		out.Navigate = "watch"
		out.Reply = "Watch lists policy-eligible public markets. Empty is honest."
		return out
	}
	if strings.Contains(low, "compare") {
		out.Tool = "watch.get"
		out.Navigate = "watch"
		out.Reply = "Comparison uses the public book under your policy. Private research stays sealed until you start it."
		return out
	}
	if strings.Contains(low, "blocked") || strings.Contains(low, "rejected") || strings.Contains(low, "stood down") || strings.Contains(low, "kill") {
		out.Tool = "research.result"
		out.Navigate = "research"
		out.Reply = "A committee stand-down or policy block is a verified result, not a crash. Open Research for the named reason."
		return out
	}
	if strings.Contains(low, "position") || strings.Contains(low, "pnl") || strings.Contains(low, "exposure") {
		out.Tool = "positions.get"
		out.Navigate = "positions"
		out.Reply = "Positions are read from your Hyperliquid trading account, not the PIT agent address."
		return out
	}
	if strings.Contains(low, "policy") {
		out.Tool = "policy.get"
		out.Navigate = "policy"
		out.Reply = "Policy is host law. Chat cannot raise clip, leverage, or permissions."
		return out
	}
	if strings.Contains(low, "research") || strings.Contains(low, "prepare") {
		coin := firstCoin(low)
		if coin == "" {
			coin = "ETH"
		}
		out.Tool = "research.start"
		out.StartResearch = true
		out.Coin = coin
		out.Mutate = true
		out.Reply = fmt.Sprintf("%s is eligible only if your policy allows it. I will start a sealed Direct research pass. That spends private compute, not trading capital. It will not place an order.", coin)
		return out
	}
	if strings.Contains(low, "why") {
		coin := firstCoin(low)
		out.Tool = "watch.get"
		out.Navigate = "watch"
		if coin == "" {
			out.Reply = "Public Watch shows mark, oracle, funding, and open interest. Private thesis stays sealed."
		} else {
			out.Coin = coin
			out.Reply = fmt.Sprintf("I can explain the public %s book under your policy. Private research stays sealed until you start it.", coin)
		}
		return out
	}
	if _, err := url.ParseRequestURI(raw); err == nil {
		out.Tool = "help"
		out.Reply = "PIT will not fetch arbitrary URLs. Use Open Hyperliquid or Open 0G Private Compute."
		return out
	}
	out.Tool = "help"
	out.Reply = "I can research a policy market, explain Watch, show evidence, open official Hyperliquid or 0G pages, or explain policy and session. I cannot execute."
	return out
}

func wantsFlatten(low string) bool {
	return strings.Contains(low, "flatten") || strings.Contains(low, "close my position") || strings.Contains(low, "close the position")
}

func wantsExecute(low string) bool {
	needles := []string{
		"buy now", "sell now", "trade now", "just do it", "just buy", "just sell",
		"i authorize", "i authorise", "execute this", "place the order",
		"go long now", "go short now",
	}
	for _, n := range needles {
		if strings.Contains(low, n) {
			return true
		}
	}
	if strings.HasPrefix(low, "buy ") || strings.HasPrefix(low, "sell ") {
		return true
	}
	return false
}

func firstCoin(low string) string {
	for coin := range allowedCoins {
		if strings.Contains(low, strings.ToLower(coin)) {
			return coin
		}
	}
	return ""
}
