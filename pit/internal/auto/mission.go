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
	ModeManual   = "manual"
	ModeResearch = "research_only"
	ModeGuarded  = "guarded"
	EnableToken  = "ENABLE GUARDED AUTONOMY"
	StopToken    = "STOP AUTONOMY"
)

type Mission struct {
	Mode               string   `json:"mode"`
	Running            bool     `json:"running"`
	Hours              int      `json:"hours,omitempty"`
	DeadlineUnix       int64    `json:"deadline_unix,omitempty"`
	GuardedEnabledUnix int64    `json:"guarded_enabled_unix,omitempty"`
	GuardedUntilUnix   int64    `json:"guarded_until_unix,omitempty"`
	PolicyHash         string   `json:"policy_hash,omitempty"`
	MaxTrades          int      `json:"max_trades,omitempty"`
	StopLossUSD        float64  `json:"stop_loss_usd,omitempty"`
	MinLiquidityUSD    float64  `json:"min_liquidity_usd,omitempty"`
	PauseUncertain     bool     `json:"pause_uncertain,omitempty"`
	Assets             []string `json:"assets,omitempty"`
	TradesToday        int      `json:"trades_today"`
	TradesDay          string   `json:"trades_day,omitempty"`
	ConsecutiveLosses  int      `json:"consecutive_losses"`
	LastAction         string   `json:"last_action,omitempty"`
	LastStop           string   `json:"last_stop,omitempty"`
	LastPreviewHash    string   `json:"last_preview_hash,omitempty"`
	LastOID            string   `json:"last_oid,omitempty"`
	BestCoin           string   `json:"best_coin,omitempty"`
	BestWhy            string   `json:"best_why,omitempty"`
	NextScanUnix       int64    `json:"next_scan_unix,omitempty"`
	CurrentPosition    string   `json:"current_position,omitempty"`
	Stage              string   `json:"stage,omitempty"`
	BlockReason        string   `json:"block_reason,omitempty"`
	ScanCount          int      `json:"scan_count,omitempty"`
	Scanned            int      `json:"scanned,omitempty"`
	Eligible           int      `json:"eligible,omitempty"`
	LastResult         string   `json:"last_result,omitempty"`
	OpenPositions      int      `json:"open_positions,omitempty"`
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
	m.BlockReason = ""
	m.LastAction = "guarded_enabled"
	m.Stage = "starting"
	m.NextScanUnix = now
	m.ScanCount = 0
	m.LastResult = ""
	if err := SaveMission(dir, m); err != nil {
		return m, err
	}
	p := Load(dir)
	p.Watch = true
	p.Notify = true
	p.AutoResearch = true
	p.Execute = false
	p.LastScanUnix = 0
	p.LastResearchCoin = ""
	p.LastNotifyCoin = ""
	_ = Save(dir, p)
	_ = ResetAway(dir)
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
	m.Stage = "stopped"
	m.BlockReason = ""
	_ = SaveMission(dir, m)
	p := Load(dir)
	p.AutoResearch = false
	p.Execute = false
	_ = Save(dir, p)
	AppendAway(dir, AwayEvent{Kind: "rejected", Why: why, Human: HumanWhy(why)})
	return LoadMission(dir)
}

func StopReason(m Mission, now int64, kill bool, sessionAlive bool, openPositions int, realizedPnL float64, pol policy.Policy) string {
	return MissionHaltReason(m, now, kill, sessionAlive, realizedPnL, pol)
}

// MissionHaltReason ends Guarded Autonomy. A full book (max_open_positions) must not halt scan/research.
func MissionHaltReason(m Mission, now int64, kill bool, sessionAlive bool, realizedPnL float64, pol policy.Policy) string {
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
	if pol.MaxConsecutiveLosses > 0 && m.ConsecutiveLosses >= pol.MaxConsecutiveLosses {
		return "consecutive_loss_limit"
	}
	return ""
}

// ExecBlockReason refuses a new order while the mission stays alive.
func ExecBlockReason(openPositions int, pol policy.Policy) string {
	if pol.MaxOpenPositions > 0 && openPositions >= pol.MaxOpenPositions {
		return "max_open_positions"
	}
	return ""
}

func Explain(code string) string {
	switch code {
	case "max_open_positions":
		return "A new order is refused because an open position already fills the host ceiling. Scan and private research continue. Existing positions are not flattened."
	case "kill_switch":
		return "Kill switch is on. Guarded Autonomy halted. No order was placed."
	case "session_expired":
		return "The scoped Hyperliquid session is not alive. Guarded Autonomy halted. No order was placed."
	case "deadline":
		return "The confirmed duration ended. Guarded Autonomy halted."
	case "policy_changed":
		return "Pinned policy hash changed. Guarded Autonomy halted until you enable it again."
	case "max_trades":
		return "Today's trade ceiling was reached. Guarded Autonomy halted."
	case "daily_loss":
		return "Daily loss ceiling was reached. Guarded Autonomy halted. Positions were not flattened."
	case "consecutive_loss_limit":
		return "Consecutive-loss ceiling was reached. Guarded Autonomy halted."
	case "user_stop", "chat_stop":
		return "You stopped Guarded Autonomy. PIT will not place further orders until you enable it again."
	case "duplicate_preview":
		return "This exact preview was already used. A new preview is required."
	case "preview_before_guarded":
		return "That preview started before Guarded Autonomy was enabled. It cannot be auto-executed."
	case "need_pin", "unpinned":
		return "Pin a trading policy on this computer before Guarded Autonomy."
	case "insufficient_margin":
		return HumanWhy("insufficient_margin")
	case "below_min_notional":
		return HumanWhy("below_min_notional")
	case "asset_not_allowed":
		return HumanWhy("asset_not_allowed")
	case "liquidity_insufficient":
		return HumanWhy("liquidity_insufficient")
	case "slippage_too_high":
		return HumanWhy("slippage_too_high")
	case "committee_disagreement":
		return HumanWhy("committee_disagreement")
	case "research_stood_down":
		return HumanWhy("research_stood_down")
	case "TEE_VERIFY_FAIL", "tee_verify_fail":
		return HumanWhy("TEE_VERIFY_FAIL")
	case "no_opportunity":
		return HumanWhy("no_opportunity")
	default:
		if code == "" {
			return ""
		}
		return HumanWhy(code)
	}
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
	if why := MissionHaltReason(m, g.Now, g.Kill, g.SessionOK, g.RealizedPnL, g.Policy); why != "" && why != "manual" {
		return fmt.Errorf("%s", why)
	}
	if why := ExecBlockReason(g.OpenCount, g.Policy); why != "" {
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

func RecordBlock(dir, reason, coin string) {
	m := LoadMission(dir)
	m.BlockReason = reason
	if m.Stage != "researching" && m.Stage != "scanning" && m.Stage != "ranked" {
		m.Stage = "execution-blocked"
	}
	m.LastAction = "exec_blocked:" + reason
	if coin != "" {
		m.BestCoin = coin
	}
	_ = SaveMission(dir, m)
}

func RecordStage(dir, stage, action, result, coin string) {
	m := LoadMission(dir)
	if stage != "" {
		m.Stage = stage
	}
	if action != "" {
		m.LastAction = action
	}
	if result != "" {
		m.LastResult = result
	}
	if coin != "" {
		m.BestCoin = coin
	}
	_ = SaveMission(dir, m)
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
		m.Stage = "executed"
		m.LastResult = "oid:" + oid
	}
	if stop != "" {
		switch stop {
		case "max_open_positions", "duplicate_preview", "preview_before_guarded":
			m.BlockReason = stop
		default:
			m.LastStop = stop
		}
	}
	_ = SaveMission(dir, m)
}

func Life(m Mission, kill bool, now int64) string {
	if kill {
		return "BLOCKED"
	}
	if m.Mode == ModeGuarded {
		if m.GuardedUntilUnix > 0 && now >= m.GuardedUntilUnix {
			return "STOPPED"
		}
		if m.DeadlineUnix > 0 && now >= m.DeadlineUnix {
			return "STOPPED"
		}
		if m.Running {
			return "ACTIVE"
		}
		return "STOPPED"
	}
	if m.Mode == ModeResearch && m.Running {
		return "ACTIVE"
	}
	if m.LastStop != "" {
		return "STOPPED"
	}
	return "READY"
}

func RemainingSeconds(m Mission, now int64) int64 {
	until := m.GuardedUntilUnix
	if until == 0 {
		until = m.DeadlineUnix
	}
	if until <= 0 {
		return 0
	}
	left := until - now
	if left < 0 {
		return 0
	}
	return left
}

func Public(dir string) map[string]any {
	m := LoadMission(dir)
	p := Load(dir)
	p.Execute = false
	pol := policy.Peek(dir)
	hash, _ := pol.Hash()
	now := time.Now().Unix()
	remainLosses := pol.MaxConsecutiveLosses - m.ConsecutiveLosses
	if remainLosses < 0 {
		remainLosses = 0
	}
	elapsed := int64(0)
	if m.GuardedEnabledUnix > 0 {
		elapsed = now - m.GuardedEnabledUnix
		if elapsed < 0 {
			elapsed = 0
		}
	}
	view := m
	if view.Running && view.NextScanUnix > 0 && view.NextScanUnix < now {
		view.NextScanUnix = now
	}
	return map[string]any{
		"ok": true, "mission": view, "prefs": p, "execute": false, "sign": false, "trade": false,
		"mode": view.Mode, "running": view.Running && view.Mode != ModeManual,
		"status":                       Life(view, pol.KillSwitch, now),
		"policy_hash":                  hash,
		"now":                          now,
		"elapsed_seconds":              elapsed,
		"next_scan_unix":               view.NextScanUnix,
		"remaining_seconds":            RemainingSeconds(m, now),
		"remaining_consecutive_losses": remainLosses,
		"remaining_risk_usd":           pol.DailyLossUSD,
		"block_reason":                 m.BlockReason,
		"stage":                        view.Stage,
		"explain":                      Explain(m.LastStop),
		"block_explain":                Explain(m.BlockReason),
		"limits": map[string]any{
			"allowed_assets":         pol.AllowedAssets,
			"allowed_venues":         pol.AllowedVenues,
			"max_clip_usd":           pol.MaxClipUSD,
			"max_trade_usd":          pol.MaxClipUSD,
			"max_position_usd":       pol.MaxClipUSD,
			"max_leverage":           pol.MaxLeverage,
			"daily_loss_usd":         pol.DailyLossUSD,
			"max_open_positions":     pol.MaxOpenPositions,
			"max_consecutive_losses": pol.MaxConsecutiveLosses,
			"cooldown_seconds":       pol.CooldownSeconds,
			"max_slippage_bps":       pol.MaxSlippageBps,
			"min_liquidity_usd":      pol.MinLiquidityUSD,
			"max_uncertainty":        pol.MaxUncertainty,
			"session_ttl_seconds":    pol.SessionTTLSeconds,
			"session_expiry":         pol.SessionTTLSeconds,
			"kill_switch":            pol.KillSwitch,
			"withdraw":               false,
			"transfer":               false,
			"policy_mutation":        false,
			"permission_escalation":  false,
		},
		"away":         LoadAway(dir),
		"why_not_code": m.BlockReason,
		"why_not":      HumanWhy(m.BlockReason),
		"note":         "Guarded Autonomy executes only after ENABLE GUARDED AUTONOMY is confirmed on this computer. Chat cannot enable it. The model cannot change these limits.",
	}
}
