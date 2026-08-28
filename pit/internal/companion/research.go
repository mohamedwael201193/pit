package companion

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/httpx"
)

type researchJob struct {
	ID        string           `json:"id"`
	running   bool             `json:"-"`
	done      bool             `json:"-"`
	cancel    bool             `json:"-"`
	started   time.Time        `json:"-"`
	stage     string           `json:"stage"`
	coin      string           `json:"coin"`
	err       string           `json:"error,omitempty"`
	note      string           `json:"note,omitempty"`
	roles     []map[string]any `json:"roles,omitempty"`
	elapsedMS int64            `json:"elapsed_ms"`
}

func jobFile(dir string) string {
	return filepath.Join(dir, "research-job.json")
}

func classifyResearch(code string) string {
	switch code {
	case "unbound":
		return "WORKSPACE_NOT_BOUND"
	case "direct_token_required", "direct_token_expired":
		return "DIRECT_NOT_AUTHORIZED"
	case "empty_envelope":
		return "HL_MARKET_UNAVAILABLE"
	case "SPONSOR_QUOTA":
		return "SPONSOR_QUOTA"
	case "sealer_not_wired", "direct_ledger", "direct_provider_http":
		return "DIRECT_PROVIDER_UNAVAILABLE"
	case "TEE_VERIFY_FAIL":
		return "TEE_SIGNATURE_INVALID"
	case "TEE_OPEN_FAIL", "ROUTER_DOWNGRADE_DENIED", "NOT_TEEML":
		return "TEE_RESPONSE_INVALID"
	case "missing_tee_signer":
		return "TEE_SIGNER_MISMATCH"
	case "committee_incomplete", "bad_role", "duplicate_role":
		return "RESEARCHER_FAILED"
	case "asset_not_allowed", "kill_switch":
		return "POLICY_REJECTED"
	default:
		return code
	}
}

func (h *Hub) loadJobLocked() {
	raw, err := os.ReadFile(jobFile(h.Dir))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return
	}
	var p struct {
		ID        string           `json:"id"`
		Running   bool             `json:"running"`
		Done      bool             `json:"done"`
		Stage     string           `json:"stage"`
		Coin      string           `json:"coin"`
		Error     string           `json:"error"`
		Note      string           `json:"note"`
		Roles     []map[string]any `json:"roles"`
		ElapsedMS int64            `json:"elapsed_ms"`
		Started   int64            `json:"started_unix_ms"`
	}
	if json.Unmarshal(raw, &p) != nil {
		return
	}
	h.job.ID = p.ID
	h.job.running = p.Running
	h.job.done = p.Done
	h.job.stage = p.Stage
	h.job.coin = p.Coin
	h.job.err = p.Error
	h.job.note = p.Note
	h.job.roles = p.Roles
	h.job.elapsedMS = p.ElapsedMS
	if p.Started > 0 {
		h.job.started = time.UnixMilli(p.Started)
	}
	if h.job.running {
		h.job.running = false
		h.job.done = true
		if h.job.err == "" {
			h.job.err = "COMPANION_NOT_RUNNING"
			h.job.stage = "STOPPED"
		}
		h.persistJobLocked()
	}
}

func (h *Hub) persistJobLocked() {
	body := map[string]any{
		"id":              h.job.ID,
		"running":         h.job.running,
		"done":            h.job.done,
		"stage":           h.job.stage,
		"coin":            h.job.coin,
		"error":           h.job.err,
		"note":            h.job.note,
		"elapsed_ms":      h.job.elapsedMS,
		"started_unix_ms": h.job.started.UnixMilli(),
		"sign":            false,
		"trade":           false,
	}
	if !h.job.running && len(h.job.roles) > 0 {
		roles := make([]any, 0, len(h.job.roles))
		for _, r := range h.job.roles {
			roles = append(roles, r)
		}
		body["roles"] = roles
	}
	raw, err := json.Marshal(body)
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return
	}
	_ = os.WriteFile(jobFile(h.Dir), raw, 0o600)
}

func (h *Hub) snapshotResearch() map[string]any {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	if h.job.ID == "" {
		h.loadJobLocked()
	}
	elapsed := h.job.elapsedMS
	if h.job.running && !h.job.started.IsZero() {
		elapsed = time.Since(h.job.started).Milliseconds()
	}
	roles := make([]any, 0, len(h.job.roles))
	if !h.job.running {
		for _, role := range h.job.roles {
			roles = append(roles, role)
		}
	}
	body := map[string]any{
		"ok":         h.job.err == "",
		"job_id":     h.job.ID,
		"running":    h.job.running,
		"done":       h.job.done,
		"stage":      h.job.stage,
		"coin":       h.job.coin,
		"elapsed_ms": elapsed,
		"note":       h.job.note,
		"roles":      roles,
		"sign":       false,
		"trade":      false,
		"verify":     h.job.done && h.job.err == "" && len(h.job.roles) > 0,
	}
	if h.job.err != "" {
		body["error"] = h.job.err
		body["ok"] = false
	}
	if !h.job.running {
		if raw, err := os.ReadFile(filepath.Join(h.Dir, "last-research.json")); err == nil {
			if !strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
				var ev any
				if json.Unmarshal(raw, &ev) == nil {
					body["evidence"] = ev
					if len(roles) == 0 {
						if m, ok := ev.(map[string]any); ok {
							if rr, ok := m["roles"].([]any); ok {
								for _, item := range rr {
									rm, ok := item.(map[string]any)
									if !ok {
										continue
									}
									roles = append(roles, map[string]any{
										"role":          rm["role"],
										"verify_e2ee":   rm["verify_e2ee"],
										"pubkey_signer": rm["pubkey_signer"],
										"teeSigner":     rm["teeSigner"],
									})
								}
								body["roles"] = roles
							}
						}
					}
				}
			}
		}
	}
	return body
}

func (h *Hub) setStage(stage string) {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	if h.job.cancel {
		return
	}
	h.job.stage = stage
	if !h.job.started.IsZero() {
		h.job.elapsedMS = time.Since(h.job.started).Milliseconds()
	}
	h.persistJobLocked()
}

func (h *Hub) cancelled() bool {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	return h.job.cancel
}

func (h *Hub) beginResearch(coin string) {
	want := strings.ToUpper(strings.TrimSpace(coin))
	if want == "" {
		want = "ETH"
	}
	h.researchMu.Lock()
	if h.job.running {
		h.researchMu.Unlock()
		return
	}
	h.job = researchJob{
		ID:      uuid.NewString(),
		running: true,
		started: time.Now(),
		stage:   "READING_MARKET",
		coin:    want,
	}
	h.persistJobLocked()
	h.researchMu.Unlock()
	go h.execResearch(want)
}

func (h *Hub) execResearch(coin string) {
	rep, err := cli.RunWorkspaceResearchStage(h.Dir, coin, h.setStage, h.cancelled)
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	h.job.running = false
	h.job.done = true
	if !h.job.started.IsZero() {
		h.job.elapsedMS = time.Since(h.job.started).Milliseconds()
	}
	if h.job.cancel && (err == nil || err.Error() == "research_cancelled") {
		h.job.err = "research_cancelled"
		h.job.stage = "STOPPED"
		h.persistJobLocked()
		return
	}
	h.job.note = rep.Note
	h.job.roles = rep.Roles
	if err != nil {
		h.job.err = classifyResearch(err.Error())
		h.job.stage = "STOPPED"
		h.persistJobLocked()
		return
	}
	h.job.err = ""
	h.persistJobLocked()
}

func (h *Hub) localResearchStart(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Coin string `json:"coin"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	h.beginResearch(body.Coin)
	writeLocal(w, http.StatusOK, h.snapshotResearch())
}

func (h *Hub) localResearchStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	writeLocal(w, http.StatusOK, h.snapshotResearch())
}

func (h *Hub) localResearchResult(w http.ResponseWriter, r *http.Request) {
	h.localResearchStatus(w, r)
}

func (h *Hub) localResearchCancel(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	h.researchMu.Lock()
	if h.job.running {
		h.job.cancel = true
	}
	h.researchMu.Unlock()
	writeLocal(w, http.StatusOK, h.snapshotResearch())
}
