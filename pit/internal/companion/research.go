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
	ID          string           `json:"id"`
	PID         int              `json:"pid"`
	running     bool             `json:"-"`
	done        bool             `json:"-"`
	cancel      bool             `json:"-"`
	started     time.Time        `json:"-"`
	stage       string           `json:"stage"`
	coin        string           `json:"coin"`
	err         string           `json:"error,omitempty"`
	note        string           `json:"note,omitempty"`
	roles       []map[string]any `json:"roles,omitempty"`
	elapsedMS   int64            `json:"elapsed_ms"`
	preview     map[string]any   `json:"-"`
	previewHash string           `json:"-"`
	deny        string           `json:"-"`
	eligible    bool             `json:"-"`
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
	case "sealer_not_wired", "direct_provider_http", "DIRECT_PROVIDER_UNAVAILABLE", "sealer_runtime", "direct_no_chat_id":
		return "DIRECT_PROVIDER_UNAVAILABLE"
	case "direct_ledger":
		return "DIRECT_CREDIT_INSUFFICIENT"
	case "DIRECT_PROVIDER_TIMEOUT", "direct_signature_http":
		return "DIRECT_PROVIDER_TIMEOUT"
	case "TEE_VERIFY_FAIL":
		return "TEE_SIGNATURE_INVALID"
	case "TEE_OPEN_FAIL", "ROUTER_DOWNGRADE_DENIED", "NOT_TEEML":
		return "TEE_RESPONSE_INVALID"
	case "missing_tee_signer":
		return "TEE_SIGNER_MISMATCH"
	case "bad_role", "duplicate_role":
		return "RESEARCHER_FAILED"
	case "asset_not_allowed", "kill_switch":
		return "POLICY_REJECTED"
	case "committee_incomplete":
		return "COMMITTEE_INCOMPLETE"
	case "risk_killed", "challenger_killed", "no_side":
		return code
	default:
		return code
	}
}

func compactRolesFromDisk(dir string) []map[string]any {
	raw, err := os.ReadFile(filepath.Join(dir, "last-research.json"))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return nil
	}
	var body struct {
		Roles []map[string]any `json:"roles"`
	}
	if json.Unmarshal(raw, &body) != nil || len(body.Roles) == 0 {
		return nil
	}
	out := make([]map[string]any, 0, len(body.Roles))
	for _, rm := range body.Roles {
		role := strings.TrimSpace(fmtString(rm["role"]))
		if role == "" {
			continue
		}
		item := map[string]any{
			"role":          rm["role"],
			"verify_e2ee":   rm["verify_e2ee"],
			"pubkey_signer": rm["pubkey_signer"],
			"teeSigner":     rm["teeSigner"],
		}
		if v, ok := rm["proposed_side"]; ok {
			item["proposed_side"] = v
		}
		if v, ok := rm["survives"]; ok {
			item["survives"] = v
		}
		if v, ok := rm["kill"]; ok {
			item["kill"] = v
		}
		out = append(out, item)
	}
	return out
}

func verifiedRolesFromDisk(dir string) []map[string]any {
	out := compactRolesFromDisk(dir)
	ok := 0
	for _, rm := range out {
		if strings.EqualFold(strings.TrimSpace(fmtString(rm["verify_e2ee"])), "OK") {
			ok++
		}
	}
	if ok < 3 {
		return nil
	}
	return out
}

func fmtString(v any) string {
	s, _ := v.(string)
	return s
}

func (h *Hub) loadJobLocked() {
	raw, err := os.ReadFile(jobFile(h.Dir))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return
	}
	var p struct {
		ID          string           `json:"id"`
		PID         int              `json:"pid"`
		Running     bool             `json:"running"`
		Done        bool             `json:"done"`
		Stage       string           `json:"stage"`
		Coin        string           `json:"coin"`
		Error       string           `json:"error"`
		Note        string           `json:"note"`
		Roles       []map[string]any `json:"roles"`
		ElapsedMS   int64            `json:"elapsed_ms"`
		Started     int64            `json:"started_unix_ms"`
		Preview     map[string]any   `json:"preview"`
		PreviewHash string           `json:"preview_hash"`
		Deny        string           `json:"deny"`
		Eligible    bool             `json:"eligible"`
	}
	if json.Unmarshal(raw, &p) != nil {
		return
	}
	h.job.ID = p.ID
	h.job.PID = p.PID
	h.job.running = p.Running
	h.job.done = p.Done
	h.job.stage = p.Stage
	h.job.coin = p.Coin
	h.job.err = p.Error
	h.job.note = p.Note
	h.job.roles = p.Roles
	h.job.elapsedMS = p.ElapsedMS
	h.job.preview = p.Preview
	h.job.previewHash = p.PreviewHash
	h.job.deny = p.Deny
	h.job.eligible = p.Eligible
	if p.Started > 0 {
		h.job.started = time.UnixMilli(p.Started)
	}
	if h.job.PID == os.Getpid() && h.job.running {
		return
	}
	if h.job.running {
		if roles := verifiedRolesFromDisk(h.Dir); len(roles) > 0 {
			h.job.running = false
			h.job.done = true
			h.job.err = ""
			h.job.stage = "READY"
			h.job.roles = roles
			attachPreviewLocked(h)
			h.persistJobLocked()
			return
		}
		h.job.running = false
		h.job.done = true
		if h.job.err == "" {
			h.job.err = "COMPANION_NOT_RUNNING"
			h.job.stage = "FAILED"
		}
		h.persistJobLocked()
	}
	if !h.job.running && h.job.done && h.job.err == "" {
		if roles := verifiedRolesFromDisk(h.Dir); len(roles) >= 3 {
			if len(h.job.roles) == 0 {
				h.job.roles = roles
			}
			if h.job.stage == "" || h.job.stage == "POLICY" || h.job.stage == "STOPPED" {
				h.job.stage = "READY"
			}
			attachPreviewLocked(h)
		}
	}
}

func attachPreviewLocked(h *Hub) {
	if h.job.preview != nil && strings.TrimSpace(h.job.previewHash) != "" {
		return
	}
	rep, err := cli.CommitteeDecisionFromLastResearch(h.Dir, h.job.coin)
	if err == nil {
		h.job.deny = rep.Deny
		h.job.eligible = rep.Eligible
		if !rep.Eligible {
			h.job.preview = rep.Preview
			h.job.previewHash = ""
			return
		}
	}
	p, hash, err := cli.LoadPreview(h.Dir)
	if err != nil || !h.job.eligible {
		return
	}
	h.job.previewHash = hash
	h.job.preview = map[string]any{
		"eligible": true, "market": p.Market, "side": p.Side, "sz": p.Sz,
		"orderType": p.OrderType, "limitPx": p.LimitPx, "hash": hash, "cloid": p.Cloid,
		"expiryUnixMs": p.ExpiryUnixMs,
	}
}

func (h *Hub) persistJobLocked() {
	body := map[string]any{
		"id":              h.job.ID,
		"pid":             h.job.PID,
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
		"eligible":        h.job.eligible,
		"deny":            h.job.deny,
		"preview_hash":    h.job.previewHash,
	}
	if h.job.preview != nil {
		body["preview"] = h.job.preview
	}
	if len(h.job.roles) > 0 {
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
	if h.job.running {
		if roles := compactRolesFromDisk(h.Dir); len(roles) > 0 {
			h.job.roles = roles
		}
	}
	roles := make([]any, 0, len(h.job.roles))
	for _, role := range h.job.roles {
		roles = append(roles, role)
	}
	body := map[string]any{
		"ok":           h.job.err == "",
		"job_id":       h.job.ID,
		"running":      h.job.running,
		"done":         h.job.done,
		"stage":        h.job.stage,
		"coin":         h.job.coin,
		"elapsed_ms":   elapsed,
		"note":         h.job.note,
		"roles":        roles,
		"sign":         false,
		"trade":        false,
		"verify":       h.job.done && h.job.err == "" && len(h.job.roles) > 0,
		"eligible":     h.job.eligible,
		"deny":         h.job.deny,
		"preview":      h.job.preview,
		"preview_hash": h.job.previewHash,
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
					if m, ok := ev.(map[string]any); ok {
						if rr, ok := m["roles"].([]any); ok && len(rr) >= 3 {
							roles = roles[:0]
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
									"proposed_side": rm["proposed_side"],
									"survives":      rm["survives"],
									"kill":          rm["kill"],
								})
							}
							body["roles"] = roles
							body["verify"] = true
							if h.job.preview != nil && strings.TrimSpace(h.job.previewHash) != "" {
								body["preview"] = h.job.preview
								body["preview_hash"] = h.job.previewHash
								body["deny"] = h.job.deny
								body["eligible"] = h.job.eligible
							} else if rep, rerr := cli.CommitteeDecisionFromLastResearch(h.Dir, h.job.coin); rerr == nil {
								body["preview"] = rep.Preview
								body["preview_hash"] = rep.PreviewHash
								body["deny"] = rep.Deny
								body["eligible"] = rep.Eligible
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
	if roles := compactRolesFromDisk(h.Dir); len(roles) > 0 {
		h.job.roles = roles
	}
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
		PID:     os.Getpid(),
		running: true,
		started: time.Now(),
		stage:   "READING_MARKET",
		coin:    want,
	}
	_ = os.WriteFile(filepath.Join(h.Dir, "last-research.json"), []byte(`{"sign":false,"trade":false,"roles":[]}`), 0o600)
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
		h.job.stage = "CANCELED"
		h.persistJobLocked()
		return
	}
	h.job.note = rep.Note
	h.job.roles = rep.Roles
	h.job.preview = rep.Preview
	h.job.previewHash = rep.PreviewHash
	h.job.deny = rep.Deny
	h.job.eligible = rep.Eligible
	if err != nil {
		h.job.err = classifyResearch(err.Error())
		h.job.stage = "FAILED"
		if len(h.job.roles) == 0 {
			if roles := compactRolesFromDisk(h.Dir); len(roles) > 0 {
				h.job.roles = roles
			}
		}
		h.persistJobLocked()
		return
	}
	h.job.err = ""
	h.job.stage = "READY"
	h.persistJobLocked()
}

func (h *Hub) localResearchStart(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Coin       string `json:"coin"`
		Hypothesis string `json:"hypothesis"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if strings.TrimSpace(body.Hypothesis) != "" {
		if err := cli.SaveHypothesis(h.Dir, body.Hypothesis); err != nil {
			writeLocal(w, http.StatusBadRequest, map[string]any{"ok": false, "error": err.Error(), "sign": false, "trade": false})
			return
		}
	}
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

func (h *Hub) localAuthorize(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Typed string `json:"typed"`
		Hash  string `json:"hash"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	got := cli.ExecuteDeskOrder(h.Dir, body.Typed, body.Hash)
	out := map[string]any{
		"ok": got.OK, "posted": got.Posted, "oid": got.OID, "cloid": got.Cloid,
		"hash": got.Hash, "market": got.Market, "side": got.Side, "sz": got.Sz,
		"agent": got.Agent, "sign": false, "trade": false, "order": true, "cancel": true,
	}
	if got.Error != "" {
		out["ok"] = false
		out["error"] = got.Error
	}
	writeLocal(w, http.StatusOK, out)
}

func (h *Hub) localCancelOrder(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Typed string `json:"typed"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	got := cli.ExecuteDeskCancel(h.Dir, body.Typed)
	out := map[string]any{
		"ok": got.OK, "posted": got.Posted, "cloid": got.Cloid, "hash": got.Hash,
		"market": got.Market, "agent": got.Agent, "sign": false, "trade": false, "order": true, "cancel": true,
	}
	if got.Error != "" {
		out["ok"] = false
		out["error"] = got.Error
	}
	writeLocal(w, http.StatusOK, out)
}
