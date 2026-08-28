package auto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/policy"
)

const (
	ModeManual    = "manual"
	ModeResearch  = "research_only"
	ModeGuarded   = "guarded"
	EnableToken   = "ENABLE GUARDED AUTONOMY"
	StopToken     = "STOP AUTONOMY"
)

type Mission struct {
	Mode               string  `json:"mode"`
	Running            bool    `json:"running"`
	Hours              int     `json:"hours,omitempty"`
	DeadlineUnix       int64   `json:"deadline_unix,omitempty"`
	GuardedEnabledUnix int64   `json:"guarded_enabled_unix,omitempty"`
	GuardedUntilUnix   int64   `json:"guarded_until_unix,omitempty"`
	PolicyHash         string  `json:"policy_hash,omitempty"`
	MaxTrades          int     `json:"max_trades,omitempty"`
	StopLossUSD        float64 `json:"stop_loss_usd,omitempty"`
	MinLiquidityUSD    float64 `json:"min_liquidity_usd,omitempty"`
	PauseUncertain     bool    `json:"pause_uncertain,omitempty"`
	Assets             []string `json:"assets,omitempty"`
	TradesToday        int     `json:"trades_today"`
	TradesDay          string  `json:"trades_day,omitempty"`
	ConsecutiveLosses int     `json:"consecutive_losses"`
	LastAction         string  `json:"last_action,omitempty"`
	LastStop           string  `json:"last_stop,omitempty"`
	LastPreviewHash    string  `json:"last_preview_hash,omitempty"`
	LastOID            string  `json:"last_oid,omitempty"`
	BestCoin           string  `json:"best_coin,omitempty"`
	BestWhy            string  `json:"best_why,omitempty"`
	NextScanUnix       int64   `json:"next_scan_unix,omitempty"`
	CurrentPosition    string  `json:"current_position,omitempty"`
}

func DefaultMission() Mission {
	return Mission{Mode: ModeManual, Running: false}
}

func missionPath(dir string) string {
	return filepath.Join(dir, "mission.json")
}

func LoadMission(dir string) Mission {
	m := DefaultMission()
	if dir == "" {
		return m
	}
	raw, err := os.ReadFile(missionPath(dir))
	if err != nil {
		return m
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return m
	}
	_ = json.Unmarshal(raw, &m)
	m = normalizeMission(m)
	return m
}

func normalizeMission(m Mission) Mission {
	switch m.Mode {
	case ModeResearch, ModeGuarded, ModeManual:
	default:
		m.Mode = ModeManual
	}
	if m.Hours < 0 {
		m.Hours = 0
	}
	if m.Hours > 72 {
		m.Hours = 72
	}
	day := time.Now().UTC().Format("2006-01-02")
	if m.TradesDay != day {
		m.TradesToday = 0
		m.TradesDay = day
	}
	if m.Mode != ModeGuarded {
		m.Running = m.Mode == ModeResearch
		m.GuardedUntilUnix = 0
	}
	return m
}

func SaveMission(dir string, m Mission) error {
	if dir == "" {
		return fmt.Errorf("unbound")
	}
	m = normalizeMission(m)
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(missionPath(dir), raw, 0o600)
}

func ConfirmEnable(typed string) error {
	if strings.TrimSpace(typed) != EnableToken {
		return fmt.Errorf("need_ENABLE_GUARDED_AUTONOMY")
	}
	return nil
}

func EnableGuarded(dir, typed string, hours int, policyHash string) (Mission, error) {
	if err := ConfirmEnable(typed); err != nil {
		return LoadMission(dir), err
	}
	if hours <= 0 {
		hours = 24
	}
	if hours > 72 {
		hours = 72
	}
	now := time.Now().Unix()
	m := LoadMission(dir)
	m.Mode = ModeGuarded
	m.Running = true
	m.Hours = hours
	m.DeadlineUnix = now + int64(hours)*3600
	m.GuardedEnabledUnix = now
	m.GuardedUntilUnix = m.DeadlineUnix
	m.PolicyHash = policyHash
	m.LastStop = ""
	m.LastAction = "guarded_enabled"
	if err := SaveMission(dir, m); err != nil {
		return m, err
	}
	p := Load(dir)
	p.Watch = true
	p.Notify = true
	p.AutoResearch = true
	p.Execute = false
	_ = Save(dir, p)
	return LoadMission(dir), nil
}

func SetMode(dir, mode string) (Mission, error) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == ModeGuarded {
		return LoadMission(dir), fmt.Errorf("need_ENABLE_GUARDED_AUTONOMY")
	}
	m := LoadMission(dir)
	switch mode {
	case ModeResearch:
		m.Mode = ModeResearch
		m.Running = true
		m.GuardedUntilUnix = 0
		m.LastAction = "research_only"
		m.LastStop = ""
		p := Load(dir)
		p.Watch = true
		p.Notify = true
		p.AutoResearch = true
		p.Execute = false
		_ = Save(dir, p)
	default:
		m.Mode = ModeManual
		m.Running = false
		m.GuardedUntilUnix = 0
		m.LastAction = "manual"
		m.LastStop = ""
	}
	if err := SaveMission(dir, m); err != nil {
		return m, err
	}
	return LoadMission(dir), nil
}

func Stop(dir, why string) Mission {
	m := LoadMission(dir)
	if why == "" {
		why = "stopped"
	}
	m.Mode = ModeManual
	m.Running = false
	m.GuardedUntilUnix = 0
	m.LastStop = why
	m.LastAction = "stopped:" + why
	_ = SaveMission(dir, m)
	p := Load(dir)
	p.AutoResearch = false
	p.Execute = false
	_ = Save(dir, p)
	return LoadMission(dir)
}

func StopReason(m Mission, now int64, kill bool, sessionAlive bool, openPositions int, realizedPnL float64, pol policy.Policy) string {
	if m.Mode != ModeGuarded {
		if m.Mode == ModeResearch {
			return ""
		}
		return "manual"
	}
	if kill || pol.KillSwitch {
		return "kill_switch"
	}
	if !sessionAlive {
		return "session_expired"
	}
	if m.GuardedUntilUnix > 0 && now >= m.GuardedUntilUnix {
		return "deadline"
	}
	if m.DeadlineUnix > 0 && now >= m.DeadlineUnix {
		return "deadline"
	}
	hash, _ := pol.Hash()
	if m.PolicyHash != "" && hash != "" && m.PolicyHash != hash {
		return "policy_changed"
	}
	if m.MaxTrades > 0 && m.TradesToday >= m.MaxTrades {
		return "max_trades"
	}
	if m.StopLossUSD > 0 && realizedPnL <= -m.StopLossUSD {
		return "daily_loss"
	}
	if pol.DailyLossUSD > 0 && realizedPnL <= -pol.DailyLossUSD {
		return "daily_loss"
	}
	if pol.MaxOpenPositions > 0 && openPositions >= pol.MaxOpenPositions {
		return "max_open_positions"
	}
	if pol.MaxConsecutiveLosses > 0 && m.ConsecutiveLosses >= pol.MaxConsecutiveLosses {
		return "consecutive_loss_limit"
	}
	return ""
}

type ExecGate struct {
	PreviewHash string
	StartedUnix int64
	OpenCount   int
	RealizedPnL float64
	SessionOK   bool
	Kill        bool
	Now         int64
	Policy      policy.Policy
	Coin        string
}

func AllowHostExecute(dir string, g ExecGate) error {
	m := LoadMission(dir)
	if m.Mode != ModeGuarded {
		return fmt.Errorf("automation_cannot_authorize")
	}
	if g.Now == 0 {
		g.Now = time.Now().Unix()
	}
	if why := StopReason(m, g.Now, g.Kill, g.SessionOK, g.OpenCount, g.RealizedPnL, g.Policy); why != "" && why != "manual" {
		return fmt.Errorf("%s", why)
	}
	if m.GuardedEnabledUnix > 0 && g.StartedUnix > 0 && g.StartedUnix < m.GuardedEnabledUnix {
		return fmt.Errorf("preview_before_guarded")
	}
	if g.PreviewHash != "" && m.LastPreviewHash != "" && g.PreviewHash == m.LastPreviewHash {
		return fmt.Errorf("duplicate_preview")
	}
	if !g.SessionOK {
		return fmt.Errorf("session_expired")
	}
	ctx := policy.Context{
		RequestedUSD:      g.Policy.MaxClipUSD,
		RequestedLev:      1,
		Coin:              "ETH",
		MarketType:        "perp",
		Venue:             "hyperliquid",
		SessionAlive:      g.SessionOK,
		NowUnix:           g.Now,
		RealizedPnLUSD:    g.RealizedPnL,
		OpenPositions:     g.OpenCount,
		ConsecutiveLosses: m.ConsecutiveLosses,
	}
	if g.Coin != "" {
		ctx.Coin = g.Coin
	} else if m.BestCoin != "" {
		ctx.Coin = m.BestCoin
	}
	if err := policy.Check(g.Policy, ctx); err != nil {
		return err
	}
	return nil
}

func RecordAction(dir, action, coin, preview, oid, stop string) {
	m := LoadMission(dir)
	if action != "" {
		m.LastAction = action
	}
	if coin != "" {
		m.BestCoin = coin
	}
	if preview != "" {
		m.LastPreviewHash = preview
	}
	if oid != "" {
		m.LastOID = oid
		m.TradesToday++
	}
	if stop != "" {
		m.LastStop = stop
	}
	_ = SaveMission(dir, m)
}

func Public(dir string) map[string]any {
	m := LoadMission(dir)
	p := Load(dir)
	p.Execute = false
	pol := policy.Default()
	hash, _ := pol.Hash()
	now := time.Now().Unix()
	return map[string]any{
		"ok": true, "mission": m, "prefs": p, "execute": false, "sign": false, "trade": false,
		"mode": m.Mode, "running": m.Running && m.Mode != ModeManual,
		"policy_hash": hash,
		"limits": map[string]any{
			"allowed_assets":         pol.AllowedAssets,
			"allowed_venues":         pol.AllowedVenues,
			"max_clip_usd":           pol.MaxClipUSD,
			"max_leverage":           pol.MaxLeverage,
			"daily_loss_usd":         pol.DailyLossUSD,
			"max_open_positions":     pol.MaxOpenPositions,
			"max_consecutive_losses": pol.MaxConsecutiveLosses,
			"cooldown_seconds":       pol.CooldownSeconds,
			"max_slippage_bps":       pol.MaxSlippageBps,
			"min_liquidity_usd":      pol.MinLiquidityUSD,
			"max_uncertainty":        pol.MaxUncertainty,
			"session_ttl_seconds":    pol.SessionTTLSeconds,
			"withdraw":               false,
			"transfer":              false,
			"policy_mutation":        false,
			"permission_escalation": false,
		},
		"now": now,
		"note": "Guarded Autonomy executes only after ENABLE GUARDED AUTONOMY on this computer. Chat cannot enable it. The model cannot change these limits.",
	}
}
