package companion

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/httpx"
)

type researchJob struct {
	running   bool
	done      bool
	cancel    bool
	started   time.Time
	stage     string
	coin      string
	err       string
	note      string
	roles     []map[string]any
	elapsedMS int64
}

func (h *Hub) snapshotResearch() map[string]any {
	h.researchMu.Lock()
	defer h.researchMu.Unlock()
	elapsed := h.job.elapsedMS
	if h.job.running && !h.job.started.IsZero() {
		elapsed = time.Since(h.job.started).Milliseconds()
	}
	roles := make([]any, 0, len(h.job.roles))
	for _, role := range h.job.roles {
		roles = append(roles, role)
	}
	body := map[string]any{
		"ok":         h.job.err == "",
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
		running: true,
		started: time.Now(),
		stage:   "READING_MARKET",
		coin:    want,
	}
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
		return
	}
	if err != nil {
		h.job.err = err.Error()
		h.job.stage = "STOPPED"
		return
	}
	h.job.note = rep.Note
	h.job.roles = rep.Roles
	h.job.err = ""
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
