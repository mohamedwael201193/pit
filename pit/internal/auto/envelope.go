package auto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/policy"
)

// AutonomyEnvelope is created only by trusted host code at arm time.
// LLM, web, chat, MCP, and SDK cannot create it.
type AutonomyEnvelope struct {
	WorkspaceID           string   `json:"workspace_id,omitempty"`
	MissionID             string   `json:"mission_id,omitempty"`
	PolicyHash            string   `json:"policy_hash,omitempty"`
	PolicyVersion         string   `json:"policy_version,omitempty"`
	MaxTradeUSD           float64  `json:"max_trade_usd,omitempty"`
	MaxPositionUSD        float64  `json:"max_position_usd,omitempty"`
	MaxAutonomyNotional   float64  `json:"max_autonomy_notional,omitempty"`
	MaxAutonomyTrades     int      `json:"max_autonomy_trades,omitempty"`
	MaxDailyAutonomyLoss  float64  `json:"max_daily_autonomy_loss,omitempty"`
	MaxOpenPositions      int      `json:"max_open_positions,omitempty"`
	AllowedAssets         []string `json:"allowed_assets,omitempty"`
	SessionID             string   `json:"session_id,omitempty"`
	SessionExpiryUnix     int64    `json:"session_expiry_unix,omitempty"`
	Venue                 string   `json:"venue,omitempty"`
	AccountSnapshotHash   string   `json:"account_snapshot_hash,omitempty"`
	CapitalSnapshot       float64  `json:"capital_snapshot,omitempty"`
	PositionSnapshotHash  string   `json:"position_snapshot_hash,omitempty"`
	ResearchJobID         string   `json:"research_job_id,omitempty"`
	ResearchDigest        string   `json:"research_digest,omitempty"`
	TEEVerified            bool     `json:"tee_verified,omitempty"`
	TeeSigner             string   `json:"tee_signer,omitempty"`
	OnChainTeeSigner      string   `json:"on_chain_tee_signer,omitempty"`
	SkillID               string   `json:"skill_id,omitempty"`
	SkillVersion          string   `json:"skill_version,omitempty"`
	SkillStatsVersion     string   `json:"skill_stats_version,omitempty"`
	MemoryRoot            string   `json:"memory_root,omitempty"`
	MemoryVersion         string   `json:"memory_version,omitempty"`
	PreviewHash           string   `json:"preview_hash,omitempty"`
	ClientOrderID         string   `json:"client_order_id,omitempty"`
	MaxOpportunityAgeSec  int64    `json:"max_opportunity_age_sec,omitempty"`
	DataFreshnessUnix     int64    `json:"data_freshness_unix,omitempty"`
	ArmedAtUnix           int64    `json:"armed_at_unix,omitempty"`
	ExpiresAtUnix         int64    `json:"expires_at_unix,omitempty"`
}

func NewMissionID(now int64) string {
	if now <= 0 {
		now = time.Now().Unix()
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("pit-sleep-%d", now)))
	return "sm-" + hex.EncodeToString(sum[:8])
}

func SnapshotEnvelope(dir string, hours int, now int64, policyHash, sessionID, workspaceID string) AutonomyEnvelope {
	if now <= 0 {
		now = time.Now().Unix()
	}
	if hours <= 0 {
		hours = 24
	}
	if hours > 72 {
		hours = 72
	}
	pol := policy.Peek(dir)
	if policyHash == "" {
		policyHash, _ = pol.Hash()
	}
	venue := "hyperliquid"
	if len(pol.AllowedVenues) > 0 {
		venue = pol.AllowedVenues[0]
	}
	clip := pol.MaxClipUSD
	return AutonomyEnvelope{
		WorkspaceID:          workspaceID,
		MissionID:            NewMissionID(now),
		PolicyHash:           policyHash,
		PolicyVersion:        pol.Version,
		MaxTradeUSD:          clip,
		MaxPositionUSD:       clip,
		MaxAutonomyNotional:  clip,
		MaxAutonomyTrades:    1,
		MaxDailyAutonomyLoss: pol.DailyLossUSD,
		MaxOpenPositions:     pol.MaxOpenPositions,
		AllowedAssets:        append([]string{}, pol.AllowedAssets...),
		SessionID:            sessionID,
		SessionExpiryUnix:    now + pol.SessionTTLSeconds,
		Venue:                venue,
		MaxOpportunityAgeSec: 120,
		ArmedAtUnix:          now,
		ExpiresAtUnix:        now + int64(hours)*3600,
	}
}

func (e AutonomyEnvelope) Digest() string {
	raw, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return "0x" + hex.EncodeToString(sum[:])
}

func ClampMissionLimits(hours, maxTrades int, maxNotional, dailyLoss float64, maxOpen int, assets []string, pol policy.Policy) (int, int, float64, float64, int, []string) {
	if hours <= 0 {
		hours = 24
	}
	if hours > 72 {
		hours = 72
	}
	if maxTrades <= 0 {
		maxTrades = 1
	}
	if pol.MaxConsecutiveLosses > 0 && maxTrades > pol.MaxConsecutiveLosses {
		maxTrades = pol.MaxConsecutiveLosses
	}
	if maxNotional <= 0 || maxNotional > pol.MaxClipUSD {
		maxNotional = pol.MaxClipUSD
	}
	if dailyLoss <= 0 || (pol.DailyLossUSD > 0 && dailyLoss > pol.DailyLossUSD) {
		dailyLoss = pol.DailyLossUSD
	}
	if maxOpen <= 0 {
		maxOpen = pol.MaxOpenPositions
	}
	if pol.MaxOpenPositions > 0 && maxOpen > pol.MaxOpenPositions {
		maxOpen = pol.MaxOpenPositions
	}
	allowed := map[string]bool{}
	for _, a := range pol.AllowedAssets {
		allowed[strings.ToUpper(strings.TrimSpace(a))] = true
	}
	out := []string{}
	for _, a := range assets {
		u := strings.ToUpper(strings.TrimSpace(a))
		if u != "" && allowed[u] {
			out = append(out, u)
		}
	}
	if len(out) == 0 {
		out = append([]string{}, pol.AllowedAssets...)
	}
	return hours, maxTrades, maxNotional, dailyLoss, maxOpen, out
}

func namedRefuse(code string) error {
	return fmt.Errorf("%s", code)
}

func (e AutonomyEnvelope) RefuseIfStale(live policy.Policy, g ExecGate) error {
	if e.PolicyHash == "" {
		return nil
	}
	liveHash, _ := live.Hash()
	if e.PolicyHash != "" && liveHash != "" && e.PolicyHash != liveHash {
		return namedRefuse("policy_changed")
	}
	if e.ExpiresAtUnix > 0 && g.Now >= e.ExpiresAtUnix {
		return namedRefuse("autonomy_expired")
	}
	if e.SessionID != "" && g.SessionID != "" && e.SessionID != g.SessionID {
		return namedRefuse("session_expired")
	}
	if e.WorkspaceID != "" && g.WorkspaceID != "" && e.WorkspaceID != g.WorkspaceID {
		return namedRefuse("workspace_mismatch")
	}
	if e.Venue != "" && g.Venue != "" && !strings.EqualFold(e.Venue, g.Venue) {
		return namedRefuse("venue_mismatch")
	}
	if g.Kill {
		return namedRefuse("kill_switch")
	}
	if !g.SessionOK {
		return namedRefuse("session_expired")
	}
	if g.MarketUnix > 0 && e.MaxOpportunityAgeSec > 0 && g.Now-g.MarketUnix > e.MaxOpportunityAgeSec {
		return namedRefuse("stale_market_data")
	}
	if e.PreviewHash != "" && g.PreviewHash != "" && e.PreviewHash != g.PreviewHash {
		return namedRefuse("stale_preview")
	}
	if e.AccountSnapshotHash != "" && g.AccountHash != "" && e.AccountSnapshotHash != g.AccountHash {
		return namedRefuse("stale_capital")
	}
	if e.PositionSnapshotHash != "" && g.PositionHash != "" && e.PositionSnapshotHash != g.PositionHash {
		return namedRefuse("position_limit")
	}
	if e.MaxAutonomyTrades > 0 && g.TradesToday >= e.MaxAutonomyTrades {
		return namedRefuse("max_autonomy_trades")
	}
	if e.MaxDailyAutonomyLoss > 0 && g.RealizedPnL <= -e.MaxDailyAutonomyLoss {
		return namedRefuse("daily_loss_limit")
	}
	if e.MaxAutonomyNotional > 0 && g.RequestedUSD > e.MaxAutonomyNotional {
		return namedRefuse("max_autonomy_notional")
	}
	if g.ResearchRequired && !g.ResearchVerified {
		return namedRefuse("research_unverified")
	}
	if g.TEERequired && !g.TEEVerified {
		return namedRefuse("tee_unverified")
	}
	if g.SkillRequired && !g.SkillEligible {
		return namedRefuse("skill_ineligible")
	}
	if g.ExtraAgentsMissing {
		return namedRefuse("extra_agent_missing")
	}
	if g.BelowMin {
		return namedRefuse("below_min_notional")
	}
	if g.InsufficientMargin {
		return namedRefuse("insufficient_margin")
	}
	return nil
}
