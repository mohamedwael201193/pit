package deskcmd

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var namedResearchCoin = regexp.MustCompile(`(?:research|happening with|price of|setup on|why is|tell me about)\s+([a-z0-9]{2,10})\b`)

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
	Hours         int    `json:"hours,omitempty"`
	Hypothesis    string `json:"hypothesis,omitempty"`
	Agent         *Agent `json:"agent,omitempty"`
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
	if wantsWhyNotTrade(low) {
		out.Tool = "watch.why_not"
		out.Reply = "PIT did not invent a fill. The refusal is named from live capital, policy, and the last committee result."
		return out
	}
	if strings.Contains(low, "why this setup") || strings.Contains(low, "verified experience") || (strings.Contains(low, "memory") && strings.Contains(low, "learn")) {
		out.Tool = "experience.why"
		out.Navigate = "research"
		out.Reply = "Verified experience is workspace-local. NOT ENOUGH DATA until enough cases exist. Chat cannot AUTHORIZE."
		return out
	}
	if wantsFlatten(low) {
		out.Tool = "refuse_execute"
		out.Navigate = "preview"
		out.Reply = "PIT will not flatten from chat. Prepare a reduce-only close on this computer, then type AUTHORIZE there."
		return out
	}
	if wantsOvernight(low) {
		out.Tool = "mission.status"
		out.Reply = "Sleep Mission status is host law. Chat cannot arm it. Open Automation to review and arm on this computer."
		return out
	}
	if wantsStopAutonomy(low) {
		out.Tool = "mission.stop"
		out.Mutate = true
		out.Navigate = "automation"
		out.Reply = "Stopping the Sleep Mission. PIT will not place further autonomous orders. Positions are not flattened. Chat cannot AUTHORIZE."
		return out
	}
	if wantsEnableAutonomy(low) {
		out.Tool = "mission.enable_required"
		out.Navigate = "automation"
		out.Hours = parseAutonomyHours(low)
		out.Reply = "Chat cannot arm a Sleep Mission. Open Automation, review the host limits, then confirm ARM SLEEP MISSION on this computer. ENABLE GUARDED AUTONOMY is still accepted."
		return out
	}
	if wantsTradesToday(low) {
		out.Tool = "activity.today"
		out.Navigate = "activity"
		out.Reply = "Today's PIT trades, previews, and fills are on Activity with OID and receipt. Historical fills never appear inside a new preview."
		return out
	}
	if wantsShowTxs(low) {
		out.Tool = "activity.list"
		out.Reply = "Desk ledger with OID, receipts, and explorer links. PIT will not invent a fill. Chat cannot AUTHORIZE."
		return out
	}
	if wantsAcceptPreview(low) {
		out.Tool = "preview.show"
		out.Navigate = "research"
		out.Reply = "Accepting is AUTHORIZE on this computer. Chat cannot AUTHORIZE. Open Research and type AUTHORIZE on the exact preview."
		return out
	}
	if wantsWhyEnter(low) && !wantsResearch(low) {
		out.Tool = "experience.why"
		out.Navigate = "research"
		out.Reply = "Why to enter is host rank plus sealed Direct research. NOT ENOUGH DATA until enough verified cases exist. Chat cannot AUTHORIZE."
		return out
	}
	if wantsOnchainProof(low) {
		out.Tool = "activity.proof"
		out.Reply = "On-chain proof, OID, and explorer links are shown here when they exist. The sealed prompt is never shown. PIT will not invent a hash."
		return out
	}
	if wantsWhyBetter(low) {
		out.Tool = "watch.compare"
		out.Reply = "Best opportunity is the highest host rank among policy-eligible live Hyperliquid books. Rank uses mark/oracle gap, funding, and open interest. It is not a model score."
		return out
	}
	if wantsScanAll(low) {
		out.Tool = "watch.scan"
		out.Reply = "Scanning every live Hyperliquid perp PIT can read, then filtering by your policy. Empty is honest. Side is not decided here."
		return out
	}
	if wantsBest(low) || wantsTradeStrongest(low) || wantsResearchBest(low) {
		out.Tool = "research.best"
		out.StartResearch = true
		out.Mutate = true
		if wantsBestLong(low) {
			out.Hypothesis = "long"
		}
		if wantsBestShort(low) {
			out.Hypothesis = "short"
		}
		out.Reply = "I will scan live Hyperliquid books, rank what is executable under your pinned policy, then start sealed 0G Direct on the strongest candidate. Chat cannot AUTHORIZE."
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
	if (strings.Contains(low, "last research") || strings.Contains(low, "last thesis") || strings.Contains(low, "show me the last research")) && !execAsk {
		out.Tool = "research.result"
		out.Navigate = "research"
		out.Reply = "Last sealed research is on Research with committee roles, TEE check, and the named stand-down or preview. Chat cannot AUTHORIZE."
		return out
	}
	if strings.Contains(low, "compare") && coin == "" {
		out.Tool = "watch.compare"
		out.Reply = "Host rank compares policy-eligible books by mark/oracle gap, funding, and open interest. It is not a model score."
		return out
	}
	if strings.Contains(low, "compare") && coin != "" {
		out.Tool = "watch.compare"
		out.Coin = coin
		out.Reply = "Host rank compares those policy-eligible books by mark/oracle gap, funding, and open interest. It is not a model score."
		return out
	}
	if (strings.Contains(low, "happening") || (strings.Contains(low, "why is") && strings.Contains(low, "interesting")) || (strings.Contains(low, "why is") && strings.Contains(low, "moving"))) && coin != "" {
		out.Tool = "watch.get"
		out.Coin = coin
		out.Reply = fmt.Sprintf("Live Hyperliquid state for %s. Price is the mark, not the order notional. Side is not decided here. Private thesis stays sealed.", coin)
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
	if strings.Contains(low, "learn") || strings.Contains(low, "calibration") || strings.Contains(low, "strategy health") || strings.Contains(low, "skill performance") {
		out.Tool = "calibration.get"
		out.Navigate = "health"
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
		out.Tool = "activity.proof"
		out.Reply = "Proof for the latest sealed research is shown here when a real storage root and chain transaction exist. PIT will not invent a hash."
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
	if wantsPolicyMutate(low) {
		out.Tool = "policy.get"
		out.Navigate = "security"
		out.Reply = "Chat cannot change policy. Open Security, edit the host limits, preview the consequences, then pin on this computer."
		return out
	}
	if strings.Contains(low, "position") || strings.Contains(low, "pnl") || strings.Contains(low, "exposure") || strings.Contains(low, "balance") || strings.Contains(low, "margin") || strings.Contains(low, "explain my risk") || (strings.Contains(low, "risk") && !wantsResearch(low) && !strings.Contains(low, "committee")) {
		out.Tool = "positions.get"
		out.Navigate = "portfolio"
		out.Reply = "Positions are read from your Hyperliquid trading account, not the PIT agent address. Chat cannot flatten."
		return out
	}
	if strings.Contains(low, "policy") && !wantsResearch(low) {
		out.Tool = "policy.get"
		out.Navigate = "security"
		out.Reply = "Policy is host law. Chat cannot raise clip, leverage, or permissions."
		return out
	}
	if wantsResearch(low) {
		out.Tool = "research.start"
		out.StartResearch = true
		out.Coin = coin
		out.Mutate = true
		label := coin
		if label == "" {
			label = "the strongest policy-eligible book"
		}
		out.Reply = fmt.Sprintf("%s is eligible only if your policy allows it. I will start a sealed Direct research pass. That spends private compute, not trading capital. It will not place an order.", label)
		if execAsk {
			out.Reply += " Chat cannot AUTHORIZE. Review the exact preview on this computer, then type AUTHORIZE there."
		}
		return out
	}
	if wantsWatch(low) || wantsPrice(low) || coinSetup(low, coin) {
		out.Tool = "watch.get"
		if coin != "" {
			out.Coin = coin
			out.Reply = fmt.Sprintf("Live Hyperliquid marks for %s under your policy. Price is the mark, not the order notional. Side is not decided here.", coin)
		} else {
			out.Reply = "Live Hyperliquid books with an honest policy-fit flag. Empty is honest. No invented scores."
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
	out.Reply = "Ask a specific desk question: best opportunity, scan markets, research a policy market, evidence, positions, automation, or policy. Chat cannot AUTHORIZE."
	return out
}

func wantsWhyNotTrade(low string) bool {
	return strings.Contains(low, "why didn't you trade") ||
		strings.Contains(low, "why didnt you trade") ||
		strings.Contains(low, "why didn't pit trade") ||
		strings.Contains(low, "why didn't you execute") ||
		strings.Contains(low, "why is nothing executable") ||
		strings.Contains(low, "why nothing executable")
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
	if strings.Contains(low, "authorize") || strings.Contains(low, "authorise") {
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
		strings.Contains(low, "markets") ||
		strings.Contains(low, "interesting")
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
		if wordHas(low, strings.ToLower(coin)) {
			return coin
		}
	}
	if m := namedResearchCoin.FindStringSubmatch(low); len(m) == 2 {
		tok := strings.ToLower(m[1])
		switch tok {
		case "privately", "the", "best", "this", "it", "now", "again", "setup", "one", "my", "a", "on":
			return ""
		}
		return strings.ToUpper(m[1])
	}
	return ""
}

func wordHas(low, needle string) bool {
	if needle == "" {
		return false
	}
	i := 0
	for {
		j := strings.Index(low[i:], needle)
		if j < 0 {
			return false
		}
		j += i
		leftOK := j == 0 || !isCoinChar(low[j-1])
		right := j + len(needle)
		rightOK := right >= len(low) || !isCoinChar(low[right])
		if leftOK && rightOK {
			return true
		}
		i = j + 1
		if i >= len(low) {
			return false
		}
	}
}

func isCoinChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

func wantsStopAutonomy(low string) bool {
	return strings.Contains(low, "stop autonomous") ||
		strings.Contains(low, "stop autonomy") ||
		strings.Contains(low, "stop guarded") ||
		(strings.Contains(low, "stop") && strings.Contains(low, "autonom"))
}

func wantsOvernight(low string) bool {
	return strings.Contains(low, "overnight") ||
		strings.Contains(low, "while i was away") ||
		strings.Contains(low, "while you were away") ||
		strings.Contains(low, "last night") ||
		strings.Contains(low, "good morning") ||
		strings.Contains(low, "while i sleep") ||
		strings.Contains(low, "while i'm asleep") ||
		strings.Contains(low, "while im asleep") ||
		strings.Contains(low, "would pit trade") ||
		strings.Contains(low, "trade while i sleep") ||
		strings.Contains(low, "review sleep")
}

func wantsEnableAutonomy(low string) bool {
	if strings.Contains(low, "review") || strings.Contains(low, "explain autonomy") || strings.Contains(low, "what would") {
		return false
	}
	return strings.Contains(low, "run autonom") ||
		strings.Contains(low, "guarded autonomy") ||
		strings.Contains(low, "enable guarded") ||
		strings.Contains(low, "sleep mission") ||
		strings.Contains(low, "arm sleep") ||
		strings.Contains(low, "start a sleep") ||
		strings.Contains(low, "run this strategy") ||
		(strings.Contains(low, "24 hour") && strings.Contains(low, "autonom"))
}

func parseAutonomyHours(low string) int {
	if strings.Contains(low, "72 hour") || strings.Contains(low, "72h") {
		return 72
	}
	if strings.Contains(low, "24 hour") || strings.Contains(low, "24h") {
		return 24
	}
	if strings.Contains(low, "8 hour") || strings.Contains(low, "8h") {
		return 8
	}
	if strings.Contains(low, "1 hour") || strings.Contains(low, "1h") {
		return 1
	}
	return 8
}

func wantsTradesToday(low string) bool {
	return (strings.Contains(low, "trade") && strings.Contains(low, "today")) ||
		(strings.Contains(low, "every trade") && strings.Contains(low, "today"))
}

func wantsOnchainProof(low string) bool {
	return strings.Contains(low, "on-chain") || strings.Contains(low, "onchain") ||
		(strings.Contains(low, "proof") && (strings.Contains(low, "on-chain") || strings.Contains(low, "explorer") || strings.Contains(low, "oid") || strings.Contains(low, "trade")))
}

func wantsWhyBetter(low string) bool {
	return strings.Contains(low, "better than") || (strings.Contains(low, "why is this") && strings.Contains(low, "better"))
}

func wantsScanAll(low string) bool {
	return strings.Contains(low, "scan everything") || strings.Contains(low, "scan all") ||
		strings.Contains(low, "whole market") ||
		strings.Contains(low, "what is moving") ||
		strings.Contains(low, "what's moving") ||
		strings.Contains(low, "what is executable") ||
		(strings.Contains(low, "scan") && strings.Contains(low, "policy"))
}

func wantsBest(low string) bool {
	return strings.Contains(low, "best opportunity") ||
		strings.Contains(low, "strongest opportunity") ||
		strings.Contains(low, "best setup") ||
		strings.Contains(low, "find me the best") ||
		strings.Contains(low, "find the strongest") ||
		strings.Contains(low, "find the best") ||
		strings.Contains(low, "best trade") ||
		strings.Contains(low, "what can i trade") ||
		strings.Contains(low, "best long") ||
		strings.Contains(low, "best short")
}

func wantsBestLong(low string) bool {
	return strings.Contains(low, "best long") || strings.Contains(low, "find me the best long")
}

func wantsBestShort(low string) bool {
	return strings.Contains(low, "best short") || strings.Contains(low, "find me the best short")
}

func wantsResearchBest(low string) bool {
	return wantsResearch(low) && (strings.Contains(low, "best") || strings.Contains(low, "strongest") || strings.Contains(low, "choose") || strings.Contains(low, "pick the"))
}

func wantsAcceptPreview(low string) bool {
	t := strings.Trim(low, "!?. ")
	if t == "accept" || t == "i accept" || t == "yes enter" {
		return true
	}
	return strings.Contains(low, "i accept") ||
		strings.Contains(low, "accept the preview") ||
		strings.Contains(low, "accept this") ||
		strings.Contains(low, "enter the trade") ||
		strings.Contains(low, "enter this trade") ||
		strings.Contains(low, "enter it now")
}

func wantsWhyEnter(low string) bool {
	return strings.Contains(low, "why should") ||
		strings.Contains(low, "why enter") ||
		strings.Contains(low, "should i enter") ||
		strings.Contains(low, "should we enter")
}

func wantsShowTxs(low string) bool {
	return strings.Contains(low, "show all tx") ||
		strings.Contains(low, "show txs") ||
		strings.Contains(low, "show tx") ||
		strings.Contains(low, "show transaction") ||
		strings.Contains(low, "show the ledger")
}

func wantsTradeStrongest(low string) bool {
	return strings.Contains(low, "trade the strongest") || strings.Contains(low, "trade the best")
}

func wantsPolicyMutate(low string) bool {
	if strings.Contains(low, "raise clip") || strings.Contains(low, "raise leverage") || strings.Contains(low, "increase clip") {
		return true
	}
	touch := strings.Contains(low, "change") || strings.Contains(low, "edit") || strings.Contains(low, "set") || strings.Contains(low, "raise") || strings.Contains(low, "increase")
	target := strings.Contains(low, "policy") || strings.Contains(low, "clip") || strings.Contains(low, "leverage") || strings.Contains(low, "max open")
	return touch && target
}
