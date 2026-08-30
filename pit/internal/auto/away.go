package auto

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AwayEvent struct {
	Unix  int64  `json:"unix"`
	Kind  string `json:"kind"`
	Coin  string `json:"coin,omitempty"`
	Why   string `json:"why,omitempty"`
	Human string `json:"human,omitempty"`
	OID   string `json:"oid,omitempty"`
}

type AwayLog struct {
	SinceUnix  int64       `json:"since_unix"`
	Detected   int         `json:"detected"`
	Researched int         `json:"researched"`
	Rejected   int         `json:"rejected"`
	Traded     int         `json:"traded"`
	Filled     int         `json:"filled"`
	Events     []AwayEvent `json:"events"`
}

func awayPath(dir string) string {
	return filepath.Join(dir, "away.json")
}

func HumanWhy(code string) string {
	switch strings.TrimSpace(code) {
	case "insufficient_margin":
		return "This account cannot clear that book's Hyperliquid floor. PIT will not invent size."
	case "below_min_notional", "below_minimum":
		return "This market's rounded Hyperliquid minimum is above this account's buying power. PIT will not invent size."
	case "policy_denied", "policy_fail", "asset_blocked":
		return "Pinned policy does not allow this market."
	case "asset_not_allowed":
		return "This asset is not on your allowlist."
	case "max_exposure", "max_open_positions":
		return "Open position ceiling is full. Scan continues. Positions are not flattened."
	case "cooldown":
		return "Cooldown is still running."
	case "daily_loss":
		return "Daily loss limit is reached. Positions are not flattened."
	case "consecutive_loss_limit":
		return "Consecutive loss limit is reached."
	case "liquidity_insufficient", "thin_book":
		return "Liquidity is below your floor."
	case "slippage_too_high":
		return "Estimated slippage is above your limit."
	case "committee_disagreement", "committee_incomplete":
		return "Committee did not agree. No order was placed."
	case "research_stood_down", "stood_down":
		return "Private research stood down. No order was placed."
	case "tee_verify_fail", "TEE_VERIFY_FAIL":
		return "TEE verification failed. No order was placed."
	case "session_expired":
		return "Hyperliquid session permissions no longer match this desk."
	case "kill_switch":
		return "Kill switch is on. New orders are halted."
	case "market_moved", "stale_quote", "opportunity_expired":
		return "The opportunity moved. PIT will not chase a stale quote."
	case "need_pin", "policy_changed":
		return "Your policy changed. Re-pin it before trading."
	case "user_stop", "chat_stop":
		return "You stopped Guarded Autonomy."
	case "deadline":
		return "The confirmed duration ended."
	case "max_trades":
		return "Today's autonomous trade ceiling was reached."
	case "duplicate_preview":
		return "That exact preview was already used."
	case "preview_before_guarded":
		return "That preview started before Guarded Autonomy was enabled."
	case "no_opportunity":
		return "Nothing in the live universe qualifies under your law right now."
	case "unpinned":
		return "Pin a trading policy on this computer before Guarded Autonomy."
	case "":
		return ""
	default:
		return "Host refused: " + strings.ReplaceAll(code, "_", " ") + "."
	}
}

func LoadAway(dir string) AwayLog {
	empty := AwayLog{Events: []AwayEvent{}}
	if dir == "" {
		return empty
	}
	raw, err := os.ReadFile(awayPath(dir))
	if err != nil {
		return empty
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return empty
	}
	var log AwayLog
	if json.Unmarshal(raw, &log) != nil {
		return empty
	}
	if log.Events == nil {
		log.Events = []AwayEvent{}
	}
	return log
}

func ResetAway(dir string) AwayLog {
	log := AwayLog{SinceUnix: time.Now().Unix(), Events: []AwayEvent{}}
	_ = saveAway(dir, log)
	return log
}

func AppendAway(dir string, ev AwayEvent) AwayLog {
	if dir == "" {
		return AwayLog{}
	}
	log := LoadAway(dir)
	if log.SinceUnix == 0 {
		log.SinceUnix = time.Now().Unix()
	}
	if ev.Unix == 0 {
		ev.Unix = time.Now().Unix()
	}
	if ev.Human == "" {
		ev.Human = HumanWhy(ev.Why)
	}
	if dupAway(log, ev) {
		return log
	}
	switch ev.Kind {
	case "detected":
		log.Detected++
	case "researched":
		log.Researched++
	case "rejected":
		log.Rejected++
	case "traded":
		log.Traded++
	case "filled":
		log.Filled++
	}
	log.Events = append(log.Events, ev)
	if len(log.Events) > 80 {
		log.Events = log.Events[len(log.Events)-80:]
	}
	_ = saveAway(dir, log)
	return log
}

func dupAway(log AwayLog, ev AwayEvent) bool {
	if len(log.Events) == 0 {
		return false
	}
	last := log.Events[len(log.Events)-1]
	if last.Kind != ev.Kind || last.Coin != ev.Coin || last.Why != ev.Why {
		return false
	}
	return ev.Unix-last.Unix < 90
}

func saveAway(dir string, log AwayLog) error {
	if dir == "" {
		return nil
	}
	raw, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(awayPath(dir), raw, 0o600)
}
