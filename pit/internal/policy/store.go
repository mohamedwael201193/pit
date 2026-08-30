package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ClipFloorUSD = 10
	ClipCeilUSD  = 50
	LossCeilUSD  = 500
	OpenCeil     = 5
	SlipCeilBps  = 500
	TTLCeilSec   = 86400
)

var HostAssets = []string{"ETH", "BTC", "SOL", "HYPE", "DOGE", "AVAX"}

func documentPath(dir string) string {
	return filepath.Join(dir, "policy.json")
}

func Clamp(p Policy) Policy {
	d := Default()
	if strings.TrimSpace(p.Version) == "" {
		p.Version = d.Version
	}
	if p.MaxClipUSD < ClipFloorUSD {
		p.MaxClipUSD = ClipFloorUSD
	}
	if p.MaxClipUSD > ClipCeilUSD {
		p.MaxClipUSD = ClipCeilUSD
	}
	if p.DailyLossUSD < 1 {
		p.DailyLossUSD = 1
	}
	if p.DailyLossUSD > LossCeilUSD {
		p.DailyLossUSD = LossCeilUSD
	}
	p.MaxLeverage = 1
	p.AllowedVenues = []string{"hyperliquid"}
	p.AllowedMarketTypes = []string{"perp"}
	p.AllowedAssets = clampAssets(p.AllowedAssets)
	if p.MinSkillCalibration < 0 {
		p.MinSkillCalibration = 0
	}
	if p.MinSkillCalibration > 1 {
		p.MinSkillCalibration = 1
	}
	if p.CooldownSeconds < 0 {
		p.CooldownSeconds = 0
	}
	if p.CooldownSeconds > TTLCeilSec {
		p.CooldownSeconds = TTLCeilSec
	}
	if p.SessionTTLSeconds < 300 {
		p.SessionTTLSeconds = 300
	}
	if p.SessionTTLSeconds > TTLCeilSec {
		p.SessionTTLSeconds = TTLCeilSec
	}
	if p.MaxUncertainty <= 0 {
		p.MaxUncertainty = 1
	}
	if p.MaxUncertainty > 1 {
		p.MaxUncertainty = 1
	}
	if p.MaxSlippageBps < 10 {
		p.MaxSlippageBps = 10
	}
	if p.MaxSlippageBps > SlipCeilBps {
		p.MaxSlippageBps = SlipCeilBps
	}
	if p.MinLiquidityUSD < 0 {
		p.MinLiquidityUSD = 0
	}
	if p.MaxOpenPositions < 1 {
		p.MaxOpenPositions = 1
	}
	if p.MaxOpenPositions > OpenCeil {
		p.MaxOpenPositions = OpenCeil
	}
	if p.MaxConsecutiveLosses < 1 {
		p.MaxConsecutiveLosses = 1
	}
	if p.MaxConsecutiveLosses > 10 {
		p.MaxConsecutiveLosses = 10
	}
	return p
}

func clampAssets(in []string) []string {
	allow := map[string]bool{}
	for _, a := range HostAssets {
		allow[a] = true
	}
	out := make([]string, 0, len(HostAssets))
	seen := map[string]bool{}
	for _, a := range in {
		u := strings.ToUpper(strings.TrimSpace(a))
		if !allow[u] || seen[u] {
			continue
		}
		seen[u] = true
		out = append(out, u)
	}
	if len(out) == 0 {
		return append([]string{}, HostAssets...)
	}
	return out
}

func Save(dir, workspaceID string, p Policy) (Policy, string, error) {
	p = Clamp(p)
	h, err := p.Hash()
	if err != nil {
		return Policy{}, "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return Policy{}, "", err
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return Policy{}, "", err
	}
	if err := os.WriteFile(documentPath(dir), append(b, '\n'), 0o600); err != nil {
		return Policy{}, "", err
	}
	if _, err := PinFile(dir, workspaceID, h); err != nil {
		return Policy{}, "", err
	}
	return p, h, nil
}

func Peek(dir string) Policy {
	d := Default()
	raw, err := os.ReadFile(documentPath(dir))
	if err != nil {
		return d
	}
	var p Policy
	if json.Unmarshal(raw, &p) != nil {
		return d
	}
	return Clamp(p)
}

func Load(dir, workspaceID string) Policy {
	d := Default()
	raw, err := os.ReadFile(documentPath(dir))
	if err != nil {
		return d
	}
	var p Policy
	if json.Unmarshal(raw, &p) != nil {
		return d
	}
	p = Clamp(p)
	pin, err := ReadPin(dir, workspaceID)
	if err != nil || strings.TrimSpace(pin) == "" {
		return d
	}
	h, err := p.Hash()
	if err != nil || MatchPin(pin, h) != nil {
		return d
	}
	return p
}

func Diff(from, to Policy) []string {
	from, to = Clamp(from), Clamp(to)
	out := []string{}
	if from.MaxClipUSD != to.MaxClipUSD {
		out = append(out, fmt.Sprintf("Max trade moves from $%.0f to $%.0f. Host sizes every clip to this ceiling.", from.MaxClipUSD, to.MaxClipUSD))
	}
	if from.DailyLossUSD != to.DailyLossUSD {
		out = append(out, fmt.Sprintf("Daily loss halt moves from $%.0f to $%.0f. Positions are not flattened when it trips.", from.DailyLossUSD, to.DailyLossUSD))
	}
	if from.MaxOpenPositions != to.MaxOpenPositions {
		out = append(out, fmt.Sprintf("Open position ceiling moves from %d to %d. A new order is refused while this ceiling is full. Existing positions are not flattened.", from.MaxOpenPositions, to.MaxOpenPositions))
	}
	if from.MaxConsecutiveLosses != to.MaxConsecutiveLosses {
		out = append(out, fmt.Sprintf("Consecutive loss halt moves from %d to %d.", from.MaxConsecutiveLosses, to.MaxConsecutiveLosses))
	}
	if from.MaxSlippageBps != to.MaxSlippageBps {
		out = append(out, fmt.Sprintf("Slippage band moves from %d bps to %d bps.", from.MaxSlippageBps, to.MaxSlippageBps))
	}
	if from.MinLiquidityUSD != to.MinLiquidityUSD {
		out = append(out, fmt.Sprintf("Liquidity floor moves from $%.0f to $%.0f.", from.MinLiquidityUSD, to.MinLiquidityUSD))
	}
	if from.CooldownSeconds != to.CooldownSeconds {
		out = append(out, fmt.Sprintf("Cooldown moves from %ds to %ds.", from.CooldownSeconds, to.CooldownSeconds))
	}
	if from.MaxUncertainty != to.MaxUncertainty {
		out = append(out, fmt.Sprintf("Uncertainty ceiling moves from %.2f to %.2f.", from.MaxUncertainty, to.MaxUncertainty))
	}
	if from.SessionTTLSeconds != to.SessionTTLSeconds {
		out = append(out, fmt.Sprintf("Policy session TTL moves from %ds to %ds. Venue approval can last longer.", from.SessionTTLSeconds, to.SessionTTLSeconds))
	}
	if from.KillSwitch != to.KillSwitch {
		if to.KillSwitch {
			out = append(out, "Kill switch will be ON. New AUTHORIZE posts are refused until you turn it off.")
		} else {
			out = append(out, "Kill switch will be OFF.")
		}
	}
	if strings.Join(from.AllowedAssets, ",") != strings.Join(to.AllowedAssets, ",") {
		out = append(out, "Allowed assets become "+join(to.AllowedAssets)+". Coins outside this list are blocked before a sealed request starts.")
	}
	if len(out) == 0 {
		out = append(out, "No field change. Pinning writes the same host law again.")
	}
	out = append(out, "Leverage stays 1x. Venue, market type, withdraw, and transfer cannot be changed. The model cannot pin this.")
	return out
}

func AllowedRefused(p Policy) (allowed, refused []string) {
	p = Clamp(p)
	allowed = []string{
		fmt.Sprintf("Size a perp clip up to $%.0f at 1x on Hyperliquid for %s.", p.MaxClipUSD, join(p.AllowedAssets)),
		fmt.Sprintf("Hold at most %d open positions. Scan and private research continue if that ceiling is full.", p.MaxOpenPositions),
	}
	if !p.KillSwitch {
		allowed = append(allowed, "Accept AUTHORIZE on this computer when session, policy pin, and venue margin all pass.")
	}
	refused = []string{
		"Withdraw, transfer, or change venue leverage.",
		"Place an order from chat or from a model.",
		"Invent size below a book's Hyperliquid minimum notional.",
		"Silently raise clip, assets, or kill switch.",
	}
	if p.KillSwitch {
		refused = append(refused, "Any new AUTHORIZE while kill switch is on.")
	}
	return allowed, refused
}

func ExecWhy(openPositions int, availableUSD float64, p Policy) (block, why string) {
	p = Clamp(p)
	if p.KillSwitch {
		return "kill_switch", "Kill switch is on. You flip it on Security. The model cannot."
	}
	if p.MaxOpenPositions > 0 && openPositions >= p.MaxOpenPositions {
		return "max_open_positions", "An open position already fills the host ceiling. Scan and private research continue. Existing positions are not flattened."
	}
	if availableUSD+1e-9 < ClipFloorUSD {
		if availableUSD <= 0 {
			return "insufficient_margin", "Available venue margin is $0.00. Hyperliquid needs at least the $10 minimum. PIT will not invent size."
		}
		short := ClipFloorUSD - availableUSD
		return "insufficient_margin", fmt.Sprintf("Available venue margin is $%.2f — $%.2f short of Hyperliquid's $%d rounded open floor. Per-book minimums can be higher after szDecimals. PIT will not invent size.", availableUSD, short, ClipFloorUSD)
	}
	if p.MaxClipUSD+1e-9 < ClipFloorUSD {
		return "below_min_notional", "Policy clip is below the $10 venue minimum. Raise max trade on Security, then pin."
	}
	return "", ""
}
