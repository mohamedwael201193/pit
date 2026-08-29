package demo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const Label = "DEMO REHEARSAL — recorded real evidence, not live"

type View struct {
	OK       bool             `json:"ok"`
	Mode     string           `json:"mode"`
	Live     bool             `json:"live"`
	Label    string           `json:"label"`
	Note     string           `json:"note"`
	At       string           `json:"timestamp"`
	Sign     bool             `json:"sign"`
	Trade    bool             `json:"trade"`
	Count    int              `json:"count"`
	Receipts []map[string]any `json:"receipts"`
	Events   []map[string]any `json:"events"`
}

func Replay(dir string, receipts []map[string]any, events []map[string]any) View {
	if receipts == nil {
		receipts = []map[string]any{}
	}
	if events == nil {
		events = []map[string]any{}
	}
	return View{
		OK:       true,
		Mode:     "replay",
		Live:     false,
		Label:    Label,
		Note:     "These rows were recorded on this computer from real 0G/Hyperliquid evidence. They are not a live account, not a live fill, and not a simulated market.",
		At:       time.Now().UTC().Format(time.RFC3339),
		Sign:     false,
		Trade:    false,
		Count:    len(receipts),
		Receipts: receipts,
		Events:   events,
	}
}

func Live(dir string) View {
	return View{
		OK:    true,
		Mode:  "live",
		Live:  true,
		Label: "LIVE",
		Note:  "This is the live desk. Recorded evidence is listed under Activity. Replay is a separate, labeled rehearsal.",
		At:    time.Now().UTC().Format(time.RFC3339),
		Sign:  false,
		Trade: false,
	}
}

func PrefPath(dir string) string {
	return filepath.Join(dir, "demo-mode.json")
}

func LoadPref(dir string) string {
	raw, err := os.ReadFile(PrefPath(dir))
	if err != nil {
		return "live"
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return "live"
	}
	var body struct {
		Mode string `json:"mode"`
	}
	if json.Unmarshal(raw, &body) != nil {
		return "live"
	}
	if body.Mode == "replay" {
		return "replay"
	}
	return "live"
}

func SavePref(dir, mode string) error {
	if dir == "" {
		return nil
	}
	if mode != "replay" {
		mode = "live"
	}
	b, err := json.MarshalIndent(map[string]any{"mode": mode, "live": mode == "live"}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(PrefPath(dir), append(b, '\n'), 0o600)
}
