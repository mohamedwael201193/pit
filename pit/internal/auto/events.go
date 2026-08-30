package auto

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MissionEvent is a first-class persisted night-replay node. No-trade is success.
type MissionEvent struct {
	Unix     int64  `json:"unix"`
	Node     string  `json:"node"`
	Status   string  `json:"status"`
	Duration int64  `json:"duration_ms,omitempty"`
	Reason   string  `json:"reason,omitempty"`
	Human    string  `json:"human,omitempty"`
	JobID    string  `json:"job_id,omitempty"`
	OID      string  `json:"oid,omitempty"`
	Proof    string  `json:"proof,omitempty"`
	Coin     string  `json:"coin,omitempty"`
	NoTrade  bool    `json:"no_trade,omitempty"`
}

type MissionLog struct {
	MissionID    string         `json:"mission_id,omitempty"`
	SinceUnix    int64          `json:"since_unix"`
	Detected     int            `json:"opportunities_detected"`
	Researched   int            `json:"private_researches"`
	Challenger   int            `json:"challenger_rejects"`
	Risk         int            `json:"risk_rejects"`
	PolicyBlocks int            `json:"policy_blocks"`
	Executions   int            `json:"executions"`
	Fills        int            `json:"fills"`
	Receipts     int            `json:"receipts"`
	Lessons      int            `json:"lessons"`
	Skips        int            `json:"skips"`
	Proofs       int            `json:"proofs"`
	Events       []MissionEvent `json:"events"`
}

func eventsPath(dir string) string {
	return filepath.Join(dir, "mission-events.json")
}

func LoadEvents(dir string) MissionLog {
	empty := MissionLog{Events: []MissionEvent{}}
	if dir == "" {
		return empty
	}
	raw, err := os.ReadFile(eventsPath(dir))
	if err != nil {
		return empty
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return empty
	}
	var log MissionLog
	if json.Unmarshal(raw, &log) != nil {
		return empty
	}
	if log.Events == nil {
		log.Events = []MissionEvent{}
	}
	return log
}

func ResetEvents(dir, missionID string) MissionLog {
	log := MissionLog{MissionID: missionID, SinceUnix: time.Now().Unix(), Events: []MissionEvent{}}
	_ = saveEvents(dir, log)
	return log
}

func AppendEvent(dir string, ev MissionEvent) MissionLog {
	if dir == "" {
		return MissionLog{}
	}
	log := LoadEvents(dir)
	if log.SinceUnix == 0 {
		log.SinceUnix = time.Now().Unix()
	}
	if ev.Unix == 0 {
		ev.Unix = time.Now().Unix()
	}
	if ev.Human == "" {
		ev.Human = HumanWhy(ev.Reason)
	}
	node := strings.ToUpper(strings.TrimSpace(ev.Node))
	switch node {
	case "CANDIDATE":
		log.Detected++
	case "PRIVATE RESEARCH":
		log.Researched++
	case "CHALLENGER":
		if ev.NoTrade || strings.EqualFold(ev.Status, "NO-TRADE") {
			log.Challenger++
			log.Skips++
		}
	case "RISK":
		if ev.NoTrade || strings.EqualFold(ev.Status, "NO-TRADE") {
			log.Risk++
			log.Skips++
		}
	case "POLICY", "ENGINE":
		if ev.NoTrade {
			log.PolicyBlocks++
			log.Skips++
		}
	case "EXECUTION":
		log.Executions++
	case "FILL":
		log.Fills++
	case "0G PROOF":
		log.Proofs++
		log.Receipts++
	case "MEMORY", "OUTCOME":
		log.Lessons++
	}
	log.Events = append(log.Events, ev)
	if len(log.Events) > 200 {
		log.Events = log.Events[len(log.Events)-200:]
	}
	_ = saveEvents(dir, log)
	return log
}

func saveEvents(dir string, log MissionLog) error {
	if dir == "" {
		return nil
	}
	raw, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(eventsPath(dir), raw, 0o600)
}

func NightReplay(dir string) []MissionEvent {
	return LoadEvents(dir).Events
}
