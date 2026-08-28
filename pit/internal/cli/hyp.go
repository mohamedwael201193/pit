package cli

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mohamedwael201193/pit/internal/compute"
)

func hypothesisFile(dir string) string {
	return filepath.Join(dir, "hypothesis.json")
}

func SaveHypothesis(dir, side string) error {
	hyp, err := compute.ParseHypothesis(side)
	if err != nil {
		return err
	}
	b, err := json.Marshal(map[string]any{"side": hyp, "sign": false, "trade": false})
	if err != nil {
		return err
	}
	return os.WriteFile(hypothesisFile(dir), b, 0o600)
}

func LoadHypothesis(dir string) string {
	raw, err := os.ReadFile(hypothesisFile(dir))
	if err != nil {
		return "none"
	}
	var body struct {
		Side string `json:"side"`
	}
	if json.Unmarshal(raw, &body) != nil {
		return "none"
	}
	hyp, err := compute.ParseHypothesis(body.Side)
	if err != nil {
		return "none"
	}
	return hyp
}
