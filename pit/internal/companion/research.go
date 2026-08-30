package companion

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedwael201193/pit/internal/auto"
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
	updated     time.Time        `json:"-"`
	seq         int64            `json:"-"`
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
	source      string           `json:"-"`
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
	case "asset_not_allowed", "kill_switch", "policy_changed":
		return "POLICY_REJECTED"
	case "committee_incomplete":
		return "COMMITTEE_INCOMPLETE"
	case "risk_killed", "challenger_killed", "no_side":
		return code
	case "research_cancelled":
		return TermCanceledByUser
	case "429", "direct_rate_limited":
		return TermDirectRateLimited
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
		for _, k := range []string{"proposed_side", "survives", "kill", "elapsed_ms", "started_unix", "finished_unix"} {
			if v, ok := rm[k]; ok {
				item[k] = v
			}
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

func previewExpiredMap(p map[string]any, nowMs int64) bool {
	if p == nil {
		return false
	}
	var exp int64
	switch v := p["expiryUnixMs"].(type) {
	case float64:
		exp = int64(v)
	case int64:
		exp = v
	case json.Number:
		exp, _ = v.Int64()
	}
	return exp > 0 && nowMs > exp
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
		Updated     int64            `json:"updated_unix_ms"`
		Seq         int64            `json:"seq"`
		Preview     map[string]any   `json:"preview"`
		PreviewHash string           `json:"preview_hash"`
		Deny        string           `json:"deny"`
		Eligible    bool             `json:"eligible"`
		Source      string           `json:"source"`
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
	h.job.source = normalizeResearchSource(p.Source)
	h.job.seq = p.Seq
	if p.Started > 0 {
		h.job.started = time.UnixMilli(p.Started)
	}
	if p.Updated > 0 {
		h.job.updated = time.UnixMilli(p.Updated)
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
			h.job.err = "JOB_CRASHED"
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

func jobRetryable(code string) bool {
	switch code {
	case "DIRECT_PROVIDER_TIMEOUT", "DIRECT_PROVIDER_UNAVAILABLE", "JOB_CRASHED", "DIRECT_CREDIT_INSUFFICIENT", "DIRECT_NOT_AUTHORIZED":
		return true
	default:
		return false
	}
}

func rolesVerified(roles []map[string]any) bool {
	return namedRolesVerified(roles)
}

func (h *Hub) bumpLocked() {
	h.job.seq++
	h.job.updated = time.Now()
}

func (h *Hub) persistJobLocked() {
	if h.job.updated.IsZero() {
		h.job.updated = time.Now()
	}
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
		"updated_unix_ms": h.job.updated.UnixMilli(),
		"seq":             h.job.seq,
		"sign":            false,
		"trade":           false,
		"eligible":        h.job.eligible,
		"deny":            h.job.deny,
		"preview_hash":    h.job.previewHash,
		"source":          h.job.source,
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
	if h.job.updated.IsZero() {
		if !h.job.started.IsZero() {
			h.job.updated = h.job.started
		} else {
			h.job.updated = time.Now()
		}
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
	if !h.job.running {
		if disk := compactRolesFromDisk(h.Dir); len(disk) > 0 {
			h.job.roles = disk
		}
		if rolesVerified(h.job.roles) && len(h.job.roles) >= 3 {
			attachPreviewLocked(h)
		}
	}
	roles := make([]any, 0, len(h.job.roles))
	for _, role := range h.job.roles {
		roles = append(roles, role)
	}
	verified := namedRolesVerified(h.job.roles)
	if verified && !h.job.running {
		attachPreviewLocked(h)
	}
	now := time.Now()
	body := map[string]any{
		"ok":                (h.job.err == "" && h.job.running) || verified,
		"job_id":            h.job.ID,
		"workspace_id":      workspaceID(h.Dir),
		"running":           h.job.running,
		"done":              h.job.done,
		"terminal":          h.job.done && !h.job.running,
		"retryable":         jobRetryable(h.job.err) && !verified,
		"stage":             h.job.stage,
		"current_stage":     h.job.stage,
		"coin":              h.job.coin,
		"elapsed_ms":        elapsed,
		"created_at":        h.job.started.UnixMilli(),
		"updated_at":        h.job.updated.UnixMilli(),
		"heartbeat_unix_ms": now.UnixMilli(),
		"seq":               h.job.seq,
		"note":              h.job.note,
		"roles":             roles,
		"sign":              false,
		"trade":             false,
		"verify":            verified,
		"eligible":          h.job.eligible,
		"deny":              h.job.deny,
		"preview":           h.job.preview,
		"preview_hash":      h.job.previewHash,
	}
	kind := TerminalKind(h.job.running, h.job.err, h.job.deny, verified, h.job.eligible, h.job.roles)
	if kind != "" {
		body["terminal_kind"] = kind
		body["card_title"] = researchCardTitle(kind)
	}
	if h.job.err != "" && !verified {
		body["error"] = h.job.err
		body["ok"] = false
	}
	if verified && h.job.preview != nil && strings.TrimSpace(h.job.previewHash) != "" {
		body["preview"] = h.job.preview
		body["preview_hash"] = h.job.previewHash
		body["deny"] = h.job.deny
		body["eligible"] = h.job.eligible
	} else if verified && (h.job.preview == nil || strings.TrimSpace(h.job.previewHash) == "") {
		if rep, rerr := cli.CommitteeDecisionFromLastResearch(h.Dir, h.job.coin); rerr == nil {
			body["preview"] = rep.Preview
			body["preview_hash"] = rep.PreviewHash
			body["deny"] = rep.Deny
			body["eligible"] = rep.Eligible
		}
	}
	body["preview_expired"] = previewExpiredMap(h.job.preview, now.UnixMilli())
	if pv, ok := body["preview"].(map[string]any); ok && !body["preview_expired"].(bool) {
		body["preview_expired"] = previewExpiredMap(pv, now.UnixMilli())
	}
	return body
}

func (h *Hub) rememberResearchOutcome(coin, jobErr, deny string, eligible bool, hash string) {
	h.recordResearchExperience(coin, deny, jobErr, eligible, hash)
	p := auto.Load(h.Dir)
	kind, hold, latch := auto.ClassifyResearchSkip(jobErr, deny, eligible)
	if latch {
		if strings.TrimSpace(coin) != "" {
			p.LastResearchCoin = coin
		}
		_ = auto.Save(h.Dir, p)
		return
	}
	p.LastResearchCoin = ""
	p.LastScanUnix = 0
	p.RememberSkip(coin, deny, kind, hold)
	_ = auto.Save(h.Dir, p)
	m := auto.LoadMission(h.Dir)
	m.SearchNote = auto.SearchNote(coin, kind, "")
	m.LastResult = m.SearchNote
	m.Stage = "searching"
	_ = auto.SaveMission(h.Dir, m)
}

func (h *Hub) currentJobID() string {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	return h.job.ID
}

func workspaceID(dir string) string {
	st, err := cli.Load(dir)
	if err != nil {
		return ""
	}
	return st.WorkspaceID
}

func (h *Hub) setStage(stage string) {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	if h.job.cancel {
		return
	}
	h.job.stage = stage
	h.bumpLocked()
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

func normalizeResearchSource(src string) string {
	switch strings.ToLower(strings.TrimSpace(src)) {
	case "automation":
		return "automation"
	case "chat":
		return "chat"
	default:
		return "research_ui"
	}
}

func mayHostGuardedExecute(source string) bool {
	return normalizeResearchSource(source) == "automation"
}

func (h *Hub) beginResearch(coin, source string) {
	want := strings.ToUpper(strings.TrimSpace(coin))
	if want == "" {
		want = h.pickBestCoin()
	}
	if want == "" {
		return
	}
	h.researchMu.Lock()
	if h.job.running {
		h.researchMu.Unlock()
		return
	}
	now := time.Now()
	h.job = researchJob{
		ID:      uuid.NewString(),
		PID:     os.Getpid(),
		running: true,
		started: now,
		updated: now,
		seq:     1,
		stage:   "READING_MARKET",
		coin:    want,
		source:  normalizeResearchSource(source),
	}
	_ = os.WriteFile(filepath.Join(h.Dir, "last-research.json"), []byte(`{"sign":false,"trade":false,"roles":[]}`), 0o600)
	h.persistJobLocked()
	appendActivity(h.Dir, activityEvent{
		WorkspaceID: workspaceID(h.Dir), Kind: "research.started", Market: want,
		Action: "research", Status: "running", JobID: h.job.ID,
	})
	h.researchMu.Unlock()
	go h.execResearch(want)
}

func (h *Hub) execResearch(coin string) {
	rep, err := cli.RunWorkspaceResearchStage(h.Dir, coin, h.setStage, h.cancelled)
	h.researchMu.Lock()
	continueNext := false
	defer func() {
		h.researchMu.Unlock()
		if continueNext {
			go h.autoTick()
		}
	}()
	h.job.running = false
	h.job.done = true
	h.bumpLocked()
	if !h.job.started.IsZero() {
		h.job.elapsedMS = time.Since(h.job.started).Milliseconds()
	}
	if h.job.cancel && (err == nil || err.Error() == "research_cancelled") {
		h.job.err = "research_cancelled"
		h.job.stage = "CANCELED"
		h.persistJobLocked()
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "research.canceled", Market: h.job.coin,
			Action: "research", Status: TermCanceledByUser, JobID: h.job.ID, Reason: TermCanceledByUser,
		})
		auto.RecordStage(h.Dir, "research_canceled", "research_canceled", TermCanceledByUser, h.job.coin)
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
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "research.failed", Market: h.job.coin,
			Action: "research", Status: h.job.err, JobID: h.job.ID, Reason: h.job.err,
		})
		auto.RecordStage(h.Dir, "research_failed", "research_failed:"+h.job.err, h.job.err, h.job.coin)
		h.rememberResearchOutcome(h.job.coin, h.job.err, h.job.deny, false, "")
		continueNext = true
		return
	}
	h.job.err = ""
	h.job.stage = "READY"
	h.persistJobLocked()
	kind := TerminalKind(false, "", h.job.deny, namedRolesVerified(h.job.roles), h.job.eligible, h.job.roles)
	writeWorkingMemory(h.Dir, map[string]any{
		"coin": h.job.coin, "deny": h.job.deny, "eligible": h.job.eligible,
		"hash": h.job.previewHash, "kind": kind, "job_id": h.job.ID,
	})
	evKind := "research.verified"
	if kind == TermReadyStoodDown {
		evKind = "research.stood_down"
	}
	book := venueTradeLink(workspaceNetwork(h.Dir), h.job.coin)
	appendActivity(h.Dir, activityEvent{
		WorkspaceID: workspaceID(h.Dir), Kind: evKind, Market: h.job.coin,
		Action: "research", Status: kind, JobID: h.job.ID, PreviewHash: h.job.previewHash, Reason: h.job.deny, Link: book,
	})
	if namedRolesVerified(h.job.roles) {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "research.sealed", Market: h.job.coin,
			Action: "seal", Status: "ok", JobID: h.job.ID, PreviewHash: h.job.previewHash, Link: book,
		})
		for _, rm := range h.job.roles {
			role := strings.ToLower(strings.TrimSpace(fmtString(rm["role"])))
			ver := strings.TrimSpace(fmtString(rm["verify_e2ee"]))
			if role == "" || !strings.EqualFold(ver, "OK") {
				continue
			}
			appendActivity(h.Dir, activityEvent{
				WorkspaceID: workspaceID(h.Dir), Kind: role + ".verified", Market: h.job.coin,
				Action: role, Status: ver, JobID: h.job.ID, PreviewHash: h.job.previewHash, Link: book,
			})
		}
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "tee.verified", Market: h.job.coin,
			Action: "tee", Status: "VerifyE2EE", JobID: h.job.ID, PreviewHash: h.job.previewHash,
			Reason: committeeReason(h.job.roles), Link: book,
		})
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "committee.verified", Market: h.job.coin,
			Action: "committee", Status: kind, JobID: h.job.ID, PreviewHash: h.job.previewHash,
			Reason: committeeReason(h.job.roles), Link: book,
		})
	}
	if h.job.eligible && h.job.previewHash != "" {
		previewStatus := "awaiting_AUTHORIZE"
		if m := auto.LoadMission(h.Dir); m.Running && m.Mode == auto.ModeGuarded {
			previewStatus = "host_may_execute"
		}
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "preview.ready", Market: h.job.coin,
			Action: "preview", Status: previewStatus, JobID: h.job.ID, PreviewHash: h.job.previewHash, Link: book,
		})
	}
	auto.RecordStage(h.Dir, "researched", "research_done:"+kind, kind, h.job.coin)
	h.rememberResearchOutcome(h.job.coin, h.job.err, h.job.deny, h.job.eligible, h.job.previewHash)
	model, provider := researchModels()
	h.fileResearch(h.job.ID, h.job.coin, kind, h.job.deny, h.job.previewHash, h.job.roles, model, provider)
	if h.job.eligible && h.job.previewHash != "" && mayHostGuardedExecute(h.job.source) {
		hash := h.job.previewHash
		coin := h.job.coin
		started := h.job.started
		go h.maybeGuardedExecute(hash, coin, started)
		return
	}
	continueNext = true
}

func (h *Hub) localResearchStart(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Coin       string `json:"coin"`
		Hypothesis string `json:"hypothesis"`
		Source     string `json:"source"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if strings.TrimSpace(body.Hypothesis) != "" {
		if err := cli.SaveHypothesis(h.Dir, body.Hypothesis); err != nil {
			writeLocal(w, http.StatusBadRequest, map[string]any{"ok": false, "error": err.Error(), "sign": false, "trade": false})
			return
		}
	}
	h.beginResearch(body.Coin, body.Source)
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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	body := h.snapshotResearch()
	raw, err := os.ReadFile(filepath.Join(h.Dir, "last-research.json"))
	if err == nil && !strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		var ev any
		if json.Unmarshal(raw, &ev) == nil {
			body["evidence"] = ev
		}
	}
	writeLocal(w, http.StatusOK, body)
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
		"status": got.Status, "agent": got.Agent, "sign": false, "trade": false, "order": true, "cancel": true,
	}
	if got.Error != "" {
		out["ok"] = false
		out["error"] = got.Error
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "approval.rejected", PreviewHash: body.Hash,
			Action: "authorize", Status: "failed", Reason: got.Error,
		})
	} else {
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "approval.accepted", Market: got.Market,
			Action: "authorize", Status: "accepted", OID: got.OID, PreviewHash: got.Hash,
			Link: venueTradeLink(workspaceNetwork(h.Dir), got.Market),
		})
		appendActivity(h.Dir, activityEvent{
			WorkspaceID: workspaceID(h.Dir), Kind: "order.submitted", Market: got.Market,
			Action: "authorize", Status: got.Status, OID: got.OID, PreviewHash: got.Hash,
			Link: venueTradeLink(workspaceNetwork(h.Dir), got.Market),
		})
		h.recordPostedOrder(got, "authorize", h.currentJobID())
		h.fileOrder(got, h.currentJobID())
		if last := cli.LoadLastOrder(h.Dir); last != nil {
			out["lifecycle"] = last["lifecycle"]
			out["reconcile"] = last["reconcile"]
			if s, ok := last["status"].(string); ok && s != "" {
				out["status"] = s
			}
		}
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
