package policy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Policy struct {
	Version             string   `json:"version"`
	MaxClipUSD          float64  `json:"maxClipUsd"`
	DailyLossUSD        float64  `json:"dailyLossUsd"`
	MaxLeverage         int      `json:"maxLeverage"`
	AllowedAssets       []string `json:"allowedAssets"`
	AllowedMarketTypes  []string `json:"allowedMarketTypes"`
	AllowedVenues       []string `json:"allowedVenues"`
	MinSkillCalibration float64  `json:"minSkillCalibration"`
	CooldownSeconds     int64    `json:"cooldownSeconds"`
	SessionTTLSeconds   int64    `json:"sessionTtlSeconds"`
	KillSwitch          bool     `json:"killSwitch"`
	MaxUncertainty      float64  `json:"maxUncertainty"`
	MaxSlippageBps      int      `json:"maxSlippageBps"`
	MinLiquidityUSD     float64  `json:"minLiquidityUsd"`
	MaxOpenPositions    int      `json:"maxOpenPositions,omitempty"`
	MaxConsecutiveLosses int    `json:"maxConsecutiveLosses,omitempty"`
}

type Context struct {
	RequestedUSD       float64
	RequestedLev       int
	Coin               string
	MarketType         string
	Venue              string
	SkillCalib         float64
	Uncertainty        float64
	SlippageBps        int
	ImpactUSD          float64
	LastFillUnix       int64
	NowUnix            int64
	RealizedPnLUSD     float64
	SessionAlive       bool
	OpenPositions      int
	ConsecutiveLosses  int
}

func Default() Policy {
	return Policy{
		Version:              "v1",
		MaxClipUSD:           10,
		DailyLossUSD:         50,
		MaxLeverage:          1,
		AllowedAssets:        []string{"ETH", "BTC", "SOL", "HYPE", "DOGE", "AVAX"},
		AllowedMarketTypes:   []string{"perp"},
		AllowedVenues:        []string{"hyperliquid"},
		MinSkillCalibration:  0,
		CooldownSeconds:      0,
		SessionTTLSeconds:    3600,
		KillSwitch:           false,
		MaxUncertainty:       1,
		MaxSlippageBps:       80,
		MinLiquidityUSD:      0,
		MaxOpenPositions:     1,
		MaxConsecutiveLosses: 3,
	}
}

func (p Policy) Hash() (string, error) {
	// Pin identity stays on the v1 wire so adding host safety fields does not invalidate an existing pin.
	wire := struct {
		Version             string   `json:"version"`
		MaxClipUSD          float64  `json:"maxClipUsd"`
		DailyLossUSD        float64  `json:"dailyLossUsd"`
		MaxLeverage         int      `json:"maxLeverage"`
		AllowedAssets       []string `json:"allowedAssets"`
		AllowedMarketTypes  []string `json:"allowedMarketTypes"`
		AllowedVenues       []string `json:"allowedVenues"`
		MinSkillCalibration float64  `json:"minSkillCalibration"`
		CooldownSeconds     int64    `json:"cooldownSeconds"`
		SessionTTLSeconds   int64    `json:"sessionTtlSeconds"`
		KillSwitch          bool     `json:"killSwitch"`
		MaxUncertainty      float64  `json:"maxUncertainty"`
		MaxSlippageBps      int      `json:"maxSlippageBps"`
		MinLiquidityUSD     float64  `json:"minLiquidityUsd"`
	}{
		Version:             p.Version,
		MaxClipUSD:          p.MaxClipUSD,
		DailyLossUSD:        p.DailyLossUSD,
		MaxLeverage:          p.MaxLeverage,
		AllowedAssets:        p.AllowedAssets,
		AllowedMarketTypes:   p.AllowedMarketTypes,
		AllowedVenues:        p.AllowedVenues,
		MinSkillCalibration:  p.MinSkillCalibration,
		CooldownSeconds:      p.CooldownSeconds,
		SessionTTLSeconds:    p.SessionTTLSeconds,
		KillSwitch:           p.KillSwitch,
		MaxUncertainty:       p.MaxUncertainty,
		MaxSlippageBps:       p.MaxSlippageBps,
		MinLiquidityUSD:      p.MinLiquidityUSD,
	}
	b, err := json.Marshal(wire)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func contains(list []string, want string) bool {
	for _, x := range list {
		if strings.EqualFold(x, want) {
			return true
		}
	}
	return false
}

func Check(p Policy, c Context) error {
	if p.KillSwitch {
		return fmt.Errorf("kill_switch")
	}
	if !c.SessionAlive {
		return fmt.Errorf("session_expired")
	}
	if c.RequestedUSD > p.MaxClipUSD+1e-9 {
		return fmt.Errorf("max_trade")
	}
	if p.DailyLossUSD > 0 && c.RealizedPnLUSD <= -p.DailyLossUSD {
		return fmt.Errorf("daily_loss_halt")
	}
	if c.RequestedLev > 0 && p.MaxLeverage > 0 && c.RequestedLev > p.MaxLeverage {
		return fmt.Errorf("leverage")
	}
	if !contains(p.AllowedAssets, c.Coin) {
		return fmt.Errorf("asset_not_allowed")
	}
	if len(p.AllowedMarketTypes) > 0 && !contains(p.AllowedMarketTypes, c.MarketType) {
		return fmt.Errorf("market_type")
	}
	if len(p.AllowedVenues) > 0 && !contains(p.AllowedVenues, c.Venue) {
		return fmt.Errorf("venue")
	}
	if p.MinSkillCalibration > 0 && c.SkillCalib+1e-12 < p.MinSkillCalibration {
		return fmt.Errorf("calibration_floor")
	}
	if p.CooldownSeconds > 0 && c.LastFillUnix > 0 && c.NowUnix-c.LastFillUnix < p.CooldownSeconds {
		return fmt.Errorf("cooldown")
	}
	if p.MaxUncertainty > 0 && c.Uncertainty > p.MaxUncertainty {
		return fmt.Errorf("uncertainty")
	}
	if p.MaxSlippageBps > 0 && c.SlippageBps > p.MaxSlippageBps {
		return fmt.Errorf("slippage")
	}
	if p.MinLiquidityUSD > 0 && c.ImpactUSD+1e-9 < p.MinLiquidityUSD {
		return fmt.Errorf("liquidity")
	}
	if p.MaxOpenPositions > 0 && c.OpenPositions >= p.MaxOpenPositions {
		return fmt.Errorf("max_open_positions")
	}
	if p.MaxConsecutiveLosses > 0 && c.ConsecutiveLosses >= p.MaxConsecutiveLosses {
		return fmt.Errorf("consecutive_loss_limit")
	}
	if p.SessionTTLSeconds > 0 && c.NowUnix == 0 {
		c.NowUnix = time.Now().Unix()
	}
	return nil
}
