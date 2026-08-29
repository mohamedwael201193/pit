package hl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Capital is the venue truth PIT uses for execution feasibility.
// Unified / portfolio accounts report trading collateral in spotClearinghouseState.
// Standard accounts keep perp withdrawable separate from spot USDC.
type Capital struct {
	Address       string
	Abstraction   string
	SpotUSDC      float64
	SpotHold      float64
	SpotFree      float64
	PerpEquity    float64
	Withdrawable  float64
	MarginUsed    float64
	TotalNtlPos   float64
	OpenPositions int
	OpenOrders    int
	BuyingPower   float64
	PowerSource   string
	Note          string
	Positions     []Position
}

func spotUSDCHold(raw json.RawMessage) (total, hold float64) {
	var st struct {
		Balances []struct {
			Coin  string `json:"coin"`
			Total string `json:"total"`
			Hold  string `json:"hold"`
		} `json:"balances"`
	}
	if err := json.Unmarshal(raw, &st); err != nil {
		return 0, 0
	}
	for _, b := range st.Balances {
		if b.Coin == "USDC" {
			return asFloat(b.Total), asFloat(b.Hold)
		}
	}
	return 0, 0
}

func parseAbstraction(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return "unknown"
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		s = strings.TrimSpace(s)
		if s != "" {
			return s
		}
	}
	var obj map[string]any
	if json.Unmarshal(raw, &obj) != nil {
		return "unknown"
	}
	for _, k := range []string{"userAbstraction", "abstraction", "mode", "dexAbstraction"} {
		if v, ok := obj[k]; ok {
			switch t := v.(type) {
			case string:
				if strings.TrimSpace(t) != "" {
					return t
				}
			case bool:
				if t {
					return k
				}
			}
		}
	}
	return "unknown"
}

func unifyName(s string) string {
	n := strings.ToLower(strings.TrimSpace(s))
	switch {
	case n == "unifiedaccount" || n == "unified_account" || n == "unified":
		return "unifiedAccount"
	case n == "portfoliomargin" || n == "portfolio_margin" || n == "portfolio":
		return "portfolioMargin"
	case n == "disabled" || n == "standard" || n == "false":
		return "disabled"
	case n == "" || n == "unknown":
		return "unknown"
	default:
		return s
	}
}

func DecidePower(abstraction string, spotFree, withdrawable, perpEquity float64) (power float64, source, note string) {
	abs := unifyName(abstraction)
	switch abs {
	case "unifiedAccount", "portfolioMargin":
		power = spotFree
		if withdrawable > power {
			power = withdrawable
		}
		source = "unified_spot"
		note = "Hyperliquid unified/portfolio accounts use spot USDC as the source of truth for trading collateral."
		return
	case "disabled":
		power = withdrawable
		source = "perp_withdrawable"
		if spotFree > 0 && withdrawable+1e-9 < 10 {
			note = "This account is in standard mode. Spot USDC is not perp margin. PIT will not invent a transfer."
			source = "spot_not_perp_margin"
			power = withdrawable
		}
		return
	default:
		if withdrawable+1e-9 >= 10 {
			return withdrawable, "perp_withdrawable", "Using Hyperliquid perp withdrawable. Account mode was not reported."
		}
		if perpEquity+1e-9 >= 10 {
			return perpEquity, "perp_equity", "Using perp account value. Account mode was not reported."
		}
		if spotFree > 0 {
			return 0, "spot_not_perp_margin", "Spot USDC is present and perp withdrawable is below the $10 venue minimum. PIT will not treat spot as perp buying power unless Hyperliquid reports unified account mode."
		}
		return 0, "unfunded", "No executable perp margin. PIT will not invent size."
	}
}

func (c *Client) userAbstraction(user string) string {
	raw, err := c.postInfo(map[string]any{"type": "userAbstraction", "user": user})
	if err != nil {
		return "unknown"
	}
	return unifyName(parseAbstraction(raw))
}

func (c *Client) openOrderCount(user string) int {
	raw, err := c.postInfo(map[string]any{"type": "openOrders", "user": user})
	if err != nil {
		return 0
	}
	var rows []json.RawMessage
	if json.Unmarshal(raw, &rows) != nil {
		return 0
	}
	return len(rows)
}

type capitalCacheEntry struct {
	at  time.Time
	got Capital
}

var (
	capitalMu    sync.Mutex
	capitalCache = map[string]capitalCacheEntry{}
)

const capitalTTL = 5 * time.Second

func (c *Client) Capital(user string) (Capital, error) {
	user = strings.TrimSpace(user)
	if user == "" {
		return Capital{}, fmt.Errorf("unbound")
	}
	key := c.InfoURL + "|" + strings.ToLower(user)
	capitalMu.Lock()
	if e, ok := capitalCache[key]; ok && time.Since(e.at) < capitalTTL {
		got := e.got
		capitalMu.Unlock()
		return got, nil
	}
	capitalMu.Unlock()

	perpRaw, err := c.postInfo(map[string]any{"type": "clearinghouseState", "user": user})
	if err != nil {
		return Capital{}, err
	}
	spotRaw, err := c.postInfo(map[string]any{"type": "spotClearinghouseState", "user": user})
	if err != nil {
		return Capital{}, err
	}
	sum := ParseClearinghouse(perpRaw)
	pos := ParsePositions(perpRaw)
	spot, hold := spotUSDCHold(spotRaw)
	if spot == 0 {
		spot = spotUSDCFromClearinghouse(spotRaw)
	}
	free := spot - hold
	if free < 0 {
		free = 0
	}
	openN := 0
	for _, p := range pos {
		if sz, err := strconv.ParseFloat(strings.TrimSpace(p.Sz), 64); err == nil && sz != 0 {
			openN++
		}
	}
	got := Capital{
		Address:       user,
		Abstraction:   c.userAbstraction(user),
		SpotUSDC:      spot,
		SpotHold:      hold,
		SpotFree:      free,
		PerpEquity:    asFloat(sum.AccountValue),
		Withdrawable:  asFloat(sum.Withdrawable),
		MarginUsed:    asFloat(sum.TotalMarginUsed),
		TotalNtlPos:   asFloat(sum.TotalNtlPos),
		OpenPositions: openN,
		OpenOrders:    c.openOrderCount(user),
		Positions:     pos,
	}
	got.BuyingPower, got.PowerSource, got.Note = DecidePower(got.Abstraction, got.SpotFree, got.Withdrawable, got.PerpEquity)
	capitalMu.Lock()
	capitalCache[key] = capitalCacheEntry{at: time.Now(), got: got}
	capitalMu.Unlock()
	return got, nil
}
