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

func Parse(text string) Result {
	raw := strings.TrimSpace(text)
	low := strings.ToLower(raw)
	out := Result{Sign: false, Trade: false, Execute: false, Mutate: false}
	if raw == "" {
		out.Reply = "Ask about a live market, private research, positions, policy, or an official page. Chat cannot AUTHORIZE."
		out.Tool = "help"
		return out
	}
	if wantsFlatten(low) {
		out.Tool = "refuse_execute"
		out.Navigate = "preview"
		out.Reply = "PIT will not flatten from chat. Prepare a reduce-only close on this computer, then type AUTHORIZE there."
		return out
	}
	execAsk := wantsExecute(low)
	coin := firstCoin(low)

	if strings.Contains(low, "forget") && (strings.Contains(low, "memory") || strings.Contains(low, "everything")) {
		out.Tool = "memory.forget"
		out.Mutate = true
		out.Reply = "Forget wipes workspace working memory and lessons. It never deletes venue positions or receipts."
		return out
	}
	if wantsOnboard(low) && coin == "" {
		out.Tool = "setup.guide"
		out.Navigate = "setup"
		out.Reply = "Opening first-run setup. Wallet, network, Hyperliquid, session, Protect, private compute, then policy. Chat cannot AUTHORIZE."
		return out
	}
	if strings.Contains(low, "take me") && coin == "" {
		out.Tool = "status"
		out.Navigate = "home"
		out.Reply = "Desk shows the next required step. Chat cannot AUTHORIZE."
		return out
	}
	if strings.Contains(low, "happening") || strings.Contains(low, "what is pit doing") || strings.Contains(low, "waiting for") {
		out.Tool = "status"
		out.Navigate = "chat"
		out.Reply = "Desk shows live status. Chat cannot see the sealed book."
		return out
	}
	if cannotExecute(low) {
		out.Tool = "refuse_execute"
		out.Navigate = "preview"
		out.Reply = "PIT will not buy, sell, or authorize from chat. Review the exact preview on this computer, then type AUTHORIZE there."
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
	if strings.Contains(low, "open 0g") || strings.Contains(low, "private compute") || strings.Contains(low, "pc.0g") || strings.Contains(low, "0g private") {
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
	if strings.Contains(low, "preview") || strings.Contains(low, "prepare this trade") || strings.Contains(low, "prepare the trade") || strings.Contains(low, "prepare the exact") {
		if execAsk {
			out.Tool = "refuse_execute"
			out.Navigate = "preview"
			out.Reply = "PIT will not buy, sell, or authorize from chat. Review the exact preview on this computer, then type AUTHORIZE there."
			return out
		}
		out.Tool = "preview.show"
		out.Navigate = "research"
		out.Reply = "The exact preview is on Research. Chat cannot AUTHORIZE it."
		return out
	}
	if strings.Contains(low, "changed") || strings.Contains(low, "recent") || strings.Contains(low, "activity") || strings.Contains(low, "yesterday") {
		out.Tool = "activity.list"
		out.Navigate = "activity"
		out.Reply = "Recent desk events are on Activity. Historical fills never appear inside a new preview."
		return out
	}
	if strings.Contains(low, "blocked") || strings.Contains(low, "rejected") || strings.Contains(low, "stood down") || strings.Contains(low, "committee") {
		out.Tool = "research.result"
		out.Navigate = "research"
		out.Reply = "A committee stand-down or policy block is a verified result, not a crash. Open Research for the named reason."
		return out
	}
	if strings.Contains(low, "position") || strings.Contains(low, "pnl") || strings.Contains(low, "exposure") || strings.Contains(low, "explain my risk") || (strings.Contains(low, "risk") && !wantsResearch(low) && !strings.Contains(low, "committee")) {
		out.Tool = "positions.get"
		out.Navigate = "positions"
		out.Reply = "Positions are read from your Hyperliquid trading account, not the PIT agent address. Chat cannot flatten."
		return out
	}
	if strings.Contains(low, "policy") && !wantsResearch(low) {
		out.Tool = "policy.get"
		out.Navigate = "policy"
		out.Reply = "Policy is host law. Chat cannot raise clip, leverage, or permissions."
		return out
	}
	if wantsResearch(low) {
		if coin == "" {
			coin = "ETH"
		}
		out.Tool = "research.start"
		out.StartResearch = true
		out.Coin = coin
		out.Mutate = true
		out.Navigate = "research"
		out.Reply = fmt.Sprintf("%s is eligible only if your policy allows it. I will start a sealed Direct research pass. That spends private compute, not trading capital. It will not place an order.", coin)
		if execAsk {
			out.Reply += " Chat cannot AUTHORIZE. Review the exact preview on Research, then type AUTHORIZE there."
		}
		return out
	}
	if wantsWatch(low) || wantsPrice(low) || coinSetup(low, coin) {
		out.Tool = "watch.get"
		out.Navigate = "watch"
		if coin != "" {
			out.Coin = coin
			out.Reply = fmt.Sprintf("Live Hyperliquid marks for %s under your policy. Side is not decided here. Private thesis stays sealed until you start research.", coin)
		} else {
			out.Reply = "Watch lists policy-eligible public markets from the live Hyperliquid book. Empty is honest. No invented scores."
		}
		if execAsk {
			out.Navigate = "preview"
			out.Reply += " Chat cannot AUTHORIZE. Open Research for the exact preview."
		}
		return out
	}
	if execAsk {
		out.Tool = "refuse_execute"
		out.Navigate = "preview"
		out.Reply = "PIT will not buy, sell, or authorize from chat. Review the exact preview on this computer, then type AUTHORIZE there."
		return out
	}
	if greetingOnly(low) {
		out.Tool = "greet"
		out.Navigate = "chat"
		out.Reply = "Hello. Ask what is interesting, a live price, private research, or your positions. Chat cannot AUTHORIZE."
		return out
	}
	if _, err := url.ParseRequestURI(raw); err == nil {
		out.Tool = "help"
		out.Reply = "PIT will not fetch arbitrary URLs. Use Open Hyperliquid or Open 0G Private Compute."
		return out
	}
	out.Tool = "help"
	out.Reply = "Ask a specific desk question: what is interesting, a live price, research a policy market, evidence, positions, or policy. Chat cannot AUTHORIZE."
	return out
}

func wantsFlatten(low string) bool {
	return strings.Contains(low, "flatten") || strings.Contains(low, "close my position") || strings.Contains(low, "close the position")
}

func wantsExecute(low string) bool {
	needles := []string{
		"buy now", "sell now", "trade now", "do trade", "do the trade", "execute trade",
		"just do it", "just buy", "just sell", "trade it", "make a trade",
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

func cannotExecute(low string) bool {
	return strings.Contains(low, "can't execute") ||
		strings.Contains(low, "cannot execute") ||
		strings.Contains(low, "why can't pit") ||
		strings.Contains(low, "why cant pit") ||
		strings.Contains(low, "why can't this") ||
		strings.Contains(low, "why won't pit execute")
}

func wantsOnboard(low string) bool {
	return strings.Contains(low, "first run") ||
		strings.Contains(low, "open setup") ||
		strings.Contains(low, "onboard") ||
		(strings.Contains(low, "setup") && strings.Contains(low, "wizard"))
}

func coinSetup(low string, coin string) bool {
	if coin == "" {
		return false
	}
	return strings.Contains(low, "setup") || strings.Contains(low, "the setup")
}

func wantsResearch(low string) bool {
	if strings.Contains(low, "research") {
		return true
	}
	for _, n := range []string{"reasearch", "reaserch", "reseach", "reserach", "resarch", "reasearh"} {
		if strings.Contains(low, n) {
			return true
		}
	}
	return strings.Contains(low, "sealed pass") || strings.Contains(low, "research privately")
}

func wantsWatch(low string) bool {
	return strings.Contains(low, "opportunit") ||
		strings.Contains(low, "match my policy") ||
		strings.Contains(low, "watch") ||
		strings.Contains(low, "interesting") ||
		strings.Contains(low, "strongest") ||
		strings.Contains(low, "compare") ||
		strings.Contains(low, "today")
}

func wantsPrice(low string) bool {
	return strings.Contains(low, "price") || strings.Contains(low, "mark px") || strings.Contains(low, "how much is")
}

func greetingOnly(low string) bool {
	t := strings.Trim(low, "!?. ")
	switch t {
	case "hi", "hello", "hey", "yo", "gm", "thanks", "thank you", "hi there":
		return true
	}
	return false
}

func firstCoin(low string) string {
	order := []string{"ETH", "BTC", "SOL", "HYPE", "DOGE", "AVAX"}
	for _, coin := range order {
		if strings.Contains(low, strings.ToLower(coin)) {
			return coin
		}
	}
	return ""
}
