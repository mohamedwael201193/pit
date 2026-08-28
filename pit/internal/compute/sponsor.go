package compute

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const sponsorQuotaPerDay = 8

type sponsorFile struct {
	Day string         `json:"day"`
	Use map[string]int `json:"use"`
}

func sponsorQuotaPath(dir string) string {
	return filepath.Join(dir, "sponsor-quota.json")
}

func DirectSponsorPath() string {
	if p := strings.TrimSpace(os.Getenv("PIT_DIRECT_SPONSOR_FILE")); p != "" {
		return p
	}
	if p := strings.TrimSpace(os.Getenv("PIT_DIRECT_AUTH_FILE")); p != "" {
		return p
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("PIT_KEYRING")), "file") {
		return ""
	}
	for _, p := range sponsorCandidates() {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func sponsorCandidates() []string {
	out := []string{}
	if exe, err := os.Executable(); err == nil {
		out = append(out, filepath.Join(filepath.Dir(exe), "direct-sponsor.json"))
	}
	if la := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); la != "" {
		out = append(out, filepath.Join(la, "PIT", "direct-sponsor.json"))
	}
	return out
}

func LoadSponsorAuthFile() (AuthFile, error) {
	p := DirectSponsorPath()
	if p == "" {
		return AuthFile{}, fmt.Errorf("direct_token_required")
	}
	return ReadAuthFile(p)
}

func ConsumeSponsorQuota(dir, workspace string) error {
	ws := strings.TrimSpace(workspace)
	if ws == "" {
		return fmt.Errorf("SPONSOR_QUOTA")
	}
	day := time.Now().UTC().Format("2006-01-02")
	p := sponsorQuotaPath(dir)
	var f sponsorFile
	if raw, err := os.ReadFile(p); err == nil {
		_ = json.Unmarshal(raw, &f)
	}
	if f.Day != day {
		f = sponsorFile{Day: day, Use: map[string]int{}}
	}
	if f.Use == nil {
		f.Use = map[string]int{}
	}
	if f.Use[ws] >= sponsorQuotaPerDay {
		return fmt.Errorf("SPONSOR_QUOTA")
	}
	f.Use[ws]++
	raw, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("SPONSOR_QUOTA")
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return fmt.Errorf("companion_leak")
	}
	return os.WriteFile(p, raw, 0o600)
}
