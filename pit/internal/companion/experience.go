package companion

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/experience"
	"github.com/mohamedwael201193/pit/internal/httpx"
)

func experiencePath(dir string) string {
	return filepath.Join(dir, "experience.enc")
}

func loadExperience(dir string) []experience.Case {
	raw, err := os.ReadFile(experiencePath(dir))
	if err != nil {
		return nil
	}
	plain, err := openBytes(dir, string(raw))
	if err != nil || strings.Contains(strings.ToLower(string(plain)), "app-sk-") {
		return nil
	}
	var rows []experience.Case
	if json.Unmarshal(plain, &rows) != nil {
		return nil
	}
	return rows
}

func appendExperience(dir string, row experience.Case) {
	if dir == "" || strings.TrimSpace(row.Coin) == "" {
		return
	}
	if row.Unix == 0 {
		row.Unix = time.Now().Unix()
	}
	rows := loadExperience(dir)
	rows = append(rows, row)
	if len(rows) > 400 {
		rows = rows[len(rows)-400:]
	}
	raw, err := json.Marshal(rows)
	if err != nil {
		return
	}
	sealed, err := sealBytes(dir, raw)
	if err != nil {
		return
	}
	_ = os.WriteFile(experiencePath(dir), []byte(sealed), 0o600)
}

func (h *Hub) recordResearchExperience(coin, deny, jobErr string, eligible bool, hash string) {
	appendExperience(h.Dir, experience.Case{
		Coin:        strings.ToUpper(strings.TrimSpace(coin)),
		Decision:    experience.Classify(eligible, deny, jobErr),
		Executable:  eligible,
		Interesting: true,
		PreviewHash: hash,
		Why:         deny,
	})
}

func (h *Hub) recordFillExperience(coin, oid, hash string) {
	if strings.TrimSpace(oid) == "" {
		return
	}
	appendExperience(h.Dir, experience.Case{
		Coin:        strings.ToUpper(strings.TrimSpace(coin)),
		Decision:    experience.DecisionFilled,
		Executable:  true,
		Interesting: true,
		PreviewHash: hash,
		OID:         oid,
		Why:         "venue_oid",
	})
}

func ExperienceSummary(dir, coin string) map[string]any {
	rows := loadExperience(dir)
	matched := experience.ForCoin(rows, coin)
	return map[string]any{
		"ok": true, "sign": false, "trade": false, "coin": strings.ToUpper(strings.TrimSpace(coin)),
		"n": len(matched), "min_samples": experience.MinSamples,
		"why": experience.WhyThisSetup(coin, matched),
	}
}

func (h *Hub) whyThisSetup(coin string) string {
	return experience.WhyThisSetup(coin, experience.ForCoin(loadExperience(h.Dir), coin))
}

func (h *Hub) localExperience(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	coin := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("coin")))
	writeLocal(w, http.StatusOK, ExperienceSummary(h.Dir, coin))
}
