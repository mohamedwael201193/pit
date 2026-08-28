package auto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Prefs struct {
	Watch            bool     `json:"watch"`
	AutoResearch     bool     `json:"auto_research"`
	Notify           bool     `json:"notify"`
	CadenceMinutes   int      `json:"cadence_minutes"`
	Trigger          string   `json:"trigger"`
	Markets          []string `json:"markets"`
	LastScanUnix     int64    `json:"last_scan_unix"`
	LastResearchCoin string   `json:"last_research_coin"`
	LastNotifyCoin   string   `json:"last_notify_coin"`
	Execute          bool     `json:"execute,omitempty"`
}

func Default() Prefs {
	return Prefs{
		Watch:          true,
		Notify:         true,
		AutoResearch:   false,
		CadenceMinutes: 15,
		Trigger:        "policy_pass",
		Markets:        nil,
	}
}

func path(dir string) string {
	return filepath.Join(dir, "automation.json")
}

func Load(dir string) Prefs {
	p := Default()
	if dir == "" {
		return p
	}
	raw, err := os.ReadFile(path(dir))
	if err != nil {
		return p
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return p
	}
	_ = json.Unmarshal(raw, &p)
	p.Execute = false
	if p.CadenceMinutes < 5 {
		p.CadenceMinutes = 5
	}
	if p.CadenceMinutes > 240 {
		p.CadenceMinutes = 240
	}
	if p.Trigger == "" {
		p.Trigger = "policy_pass"
	}
	return p
}

func Save(dir string, p Prefs) error {
	if dir == "" {
		return fmt.Errorf("unbound")
	}
	p.Execute = false
	if p.CadenceMinutes < 5 {
		p.CadenceMinutes = 5
	}
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(path(dir), raw, 0o600)
}

func RefuseExecute() error {
	return fmt.Errorf("automation_cannot_authorize")
}

func Matches(trigger, reason string) bool {
	switch strings.ToLower(strings.TrimSpace(trigger)) {
	case "mark_below_oracle":
		return reason == "mark_below_oracle"
	case "funding":
		return reason == "funding"
	default:
		return true
	}
}
