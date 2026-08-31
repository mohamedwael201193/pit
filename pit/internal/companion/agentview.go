package companion

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskcmd"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (h *Hub) decorateWatchAgent(parsed deskcmd.Result) deskcmd.Result {
	agent := &deskcmd.Agent{Kind: kindForTool(parsed.Tool)}
	st, err := cli.Load(h.Dir)
	if err != nil || strings.TrimSpace(st.Wallet) == "" {
		agent.Executive = "This computer is not bound. Connect a wallet on Setup."
		parsed.Reply = agent.Executive
		parsed.Agent = agent
		return parsed
	}
	netName := "mainnet"
	if strings.TrimSpace(st.Network) != "" {
		netName = st.Network
	}
	net, nerr := config.ParseNetwork(netName)
	if nerr != nil {
		agent.Executive = "Wrong network. Markets stay unread."
		parsed.Reply = agent.Executive
		parsed.Agent = agent
		return parsed
	}
	cands, lerr := watch.LiveUniverse(hl.New(config.For(net)), cli.ActivePolicy(h.Dir))
	if lerr != nil {
		agent.Executive = "Hyperliquid did not return a live book. Empty is honest. No invented scores."
		parsed.Reply = agent.Executive
		parsed.Agent = agent
		return parsed
	}
	view := h.annotateWatch(watch.Public(cands, string(net)), cli.ActivePolicy(h.Dir))
	agent.Scanned = view.Scanned
	agent.Eligible = view.Count
	agent.Executable = view.ExecFeasibleN
	agent.BuyingPower = view.BuyingPower
	agent.PowerSource = view.PowerSource
	agent.Why = view.BestWhy
	row := watch.PublicCoin{}
	if len(view.Coins) > 0 {
		row = view.Coins[0]
	}
	if view.Best != nil && (parsed.Tool == "watch.best" || parsed.Tool == "watch.compare" || parsed.Tool == "research.best" || parsed.Coin == "") {
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
	if row.Coin != "" {
		agent.Best = row.Coin
		agent.Mark = row.Mark
		agent.MinNotional = row.MinNotional
		agent.HostNotional = row.HostNotional
		agent.Funding = row.Funding
		agent.OpenInterest = row.OpenInterest
		agent.Freshness = row.Freshness
	}
	for _, c := range view.Coins {
		if len(agent.Coins) >= 8 {
			break
		}
		if !c.Eligible && !c.ExecutionFeasible && parsed.Tool != "watch.scan" {
			continue
		}
		agent.Coins = append(agent.Coins, deskcmd.AgentCoin{
			Coin:              c.Coin,
			Mark:              c.Mark,
			MinNotional:       c.MinNotional,
			HostNotional:      c.HostNotional,
			Funding:           c.Funding,
			OpenInterest:      c.OpenInterest,
			ExecutionFeasible: c.ExecutionFeasible,
			Eligible:          c.Eligible,
			Why:               firstNonEmpty(c.WhyRanked, c.Why),
			ExecWhy:           firstNonEmpty(c.WhyExecutable, c.ExecWhy),
			Block:             c.Block,
			Trend:             c.Trend,
		})
	}
	if len(agent.Coins) == 0 {
		for i, c := range view.Coins {
			if i >= 6 {
				break
			}
			agent.Coins = append(agent.Coins, deskcmd.AgentCoin{
				Coin:              c.Coin,
				Mark:              c.Mark,
				MinNotional:       c.MinNotional,
				HostNotional:      c.HostNotional,
				Funding:           c.Funding,
				OpenInterest:      c.OpenInterest,
				ExecutionFeasible: c.ExecutionFeasible,
				Eligible:          c.Eligible,
				Why:               c.Why,
				Block:             c.Block,
			})
		}
	}
	best := agent.Best
	if best == "" {
		best = "none"
	}
	switch parsed.Tool {
	case "watch.scan":
		agent.Kind = "scan"
		agent.Executive = fmt.Sprintf("SCAN %d discovered, %d policy eligible, %d execution feasible. Strongest: %s. Price is the mark, not the order notional.", agent.Scanned, agent.Eligible, agent.Executable, best)
		parsed.Reply = agent.Executive
	case "watch.get":
		agent.Kind = "book"
		if parsed.Coin != "" && row.Coin != "" {
			agent.Executive = fmt.Sprintf("%s mark %s. Venue min $%.2f. Host clip $%.2f. Price is not the order notional.", row.Coin, watch.Price(row.Mark), row.MinNotional, row.HostNotional)
		} else {
			agent.Executive = fmt.Sprintf("Scanned %d live Hyperliquid perps. %d executable. Strongest: %s.", agent.Scanned, agent.Executable, best)
		}
		parsed.Reply = agent.Executive
	case "watch.compare":
		agent.Kind = "compare"
		names := make([]string, 0, 3)
		for _, c := range agent.Coins {
			if len(names) >= 3 {
				break
			}
			if c.Coin != "" {
				names = append(names, c.Coin)
			}
		}
		if len(names) == 0 && best != "none" {
			names = []string{best}
		}
		agent.Executive = fmt.Sprintf("%s rank among %d executable of %d scanned. Host rank uses mark, funding, and open interest.", strings.Join(names, ", "), agent.Executable, agent.Scanned)
		parsed.Reply = agent.Executive
	case "watch.why_not":
		agent.Kind = "no_trade"
		agent.Executive = fmt.Sprintf("NO TRADE. Scanned %d. Policy eligible %d. Executable %d. Strongest ranked: %s.", agent.Scanned, agent.Eligible, agent.Executable, best)
		parsed.Reply = agent.Executive
	default:
		agent.Kind = "hunt"
		if agent.Best == "" {
			agent.Executive = fmt.Sprintf("Scanned %d. %d policy eligible. %d executable. No book is execution-feasible on this computer right now.", agent.Scanned, agent.Eligible, agent.Executable)
			parsed.Reply = "No executable book on this computer right now."
		} else {
			hypo := strings.ToLower(strings.TrimSpace(parsed.Hypothesis))
			scan := "Scanning live Hyperliquid markets…"
			if hypo == "long" || hypo == "short" {
				scan = fmt.Sprintf("Scanning live Hyperliquid markets for a %s…", hypo)
			}
			agent.Executive = fmt.Sprintf("%s %d scanned, %d executable. Strongest: %s.", scan, agent.Scanned, agent.Executable, agent.Best)
			parsed.Reply = scan
		}
	}
	parsed.Agent = agent
	if parsed.StartResearch && parsed.Coin == "" && agent.Best != "" && agent.Best != "none" {
		skip := h.huntSkipSet()
		if skip[strings.ToUpper(agent.Best)] == "" {
			parsed.Coin = agent.Best
		} else if next := h.pickBestCoinSkipping(skip); next != "" {
			parsed.Coin = next
			agent.Best = next
		}
	}
	return parsed
}

func chatDonePayload(parsed deskcmd.Result, thread, model string, stream bool) map[string]any {
	body := map[string]any{
		"ok": true, "done": true, "reply": parsed.Reply, "tool": parsed.Tool, "mutate": parsed.Mutate,
		"execute": false, "start_research": parsed.StartResearch, "coin": parsed.Coin,
		"navigate": parsed.Navigate, "open_url": parsed.OpenURL, "thread": thread,
		"hours": parsed.Hours, "model": model, "stream": stream, "sign": false, "trade": false,
	}
	if parsed.Hypothesis != "" {
		body["hypothesis"] = parsed.Hypothesis
	}
	if parsed.Agent != nil {
		body["agent"] = parsed.Agent
	}
	return body
}

func kindForTool(tool string) string {
	switch tool {
	case "watch.scan":
		return "scan"
	case "watch.why_not":
		return "no_trade"
	case "watch.compare":
		return "compare"
	case "research.best", "research.start":
		return "hunt"
	default:
		return "scan"
	}
}
