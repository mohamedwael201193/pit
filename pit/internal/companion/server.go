package companion

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/obs"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/version"
	"github.com/mohamedwael201193/pit/internal/watch"
)

const DefaultAddr = "127.0.0.1:17373"

type device struct {
	Hash    string `json:"hash"`
	Created int64  `json:"created"`
	Label   string `json:"label"`
}

type Hub struct {
	Dir     string
	mu      sync.Mutex
	code    string
	codeExp time.Time
	devices []device

	researchMu sync.Mutex
	job        researchJob
}

func ListenAddr() (string, error) {
	a := strings.TrimSpace(os.Getenv("PIT_COMPANION_ADDR"))
	if a == "" {
		a = DefaultAddr
	}
	host, _, err := net.SplitHostPort(a)
	if err != nil {
		return "", fmt.Errorf("companion_loopback_only")
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return "", fmt.Errorf("companion_loopback_only")
	}
	return a, nil
}

func New(dir string) *Hub {
	h := &Hub{Dir: dir}
	h.rotateLocked(time.Now())
	h.researchMu.Lock()
	h.loadJobLocked()
	h.researchMu.Unlock()
	go h.autoLoop()
	return h
}

func (h *Hub) Code() (string, time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if time.Now().After(h.codeExp) || h.code == "" {
		h.rotateLocked(time.Now())
	}
	return h.code, h.codeExp
}

func (h *Hub) rotateLocked(now time.Time) {
	h.code = newCode()
	h.codeExp = now.Add(2 * time.Minute)
}

func newCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	var b [8]byte
	_, _ = rand.Read(b[:])
	out := make([]byte, 8)
	for i := range out {
		out[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(out)
}

func normalizeCode(s string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(s), "-", ""))
}

func (h *Hub) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/pair", h.pair)
	mux.HandleFunc("/status", h.status)
	mux.HandleFunc("/local/status", h.localStatus)
	mux.HandleFunc("/local/code", h.localCode)
	mux.HandleFunc("/local/doctor", h.localDoctor)
	mux.HandleFunc("/local/init", h.localInit)
	mux.HandleFunc("/local/session", h.localSession)
	mux.HandleFunc("/local/connection-preview", h.localConnectionPreview)
	mux.HandleFunc("/local/policy", h.localPolicy)
	mux.HandleFunc("/local/revoke-session", h.localRevokeSession)
	mux.HandleFunc("/local/direct-intent", h.localDirectIntent)
	mux.HandleFunc("/local/direct-complete", h.localDirectComplete)
	mux.HandleFunc("/local/direct-status", h.localDirectStatus)
	mux.HandleFunc("/local/research", h.localResearch)
	mux.HandleFunc("/local/research/start", h.localResearchStart)
	mux.HandleFunc("/local/research/status", h.localResearchStatus)
	mux.HandleFunc("/local/research/result", h.localResearchResult)
	mux.HandleFunc("/local/research/cancel", h.localResearchCancel)
	mux.HandleFunc("/local/authorize", h.localAuthorize)
	mux.HandleFunc("/local/cancel", h.localCancelOrder)
	mux.HandleFunc("/local/activity", h.localActivity)
	mux.HandleFunc("/local/positions", h.localPositions)
	mux.HandleFunc("/local/chat", h.localChat)
	mux.HandleFunc("/local/chat/log", h.localChatLog)
	mux.HandleFunc("/local/chat/threads", h.localChatThreads)
	mux.HandleFunc("/local/chat/thread", h.localChatThread)
	mux.HandleFunc("/local/memory/forget", h.localMemoryForget)
	mux.HandleFunc("/local/calibration", h.localCalibration)
	mux.HandleFunc("/local/security", h.localSecurity)
	mux.HandleFunc("/local/identity", h.localIdentity)
	mux.HandleFunc("/local/update", h.localUpdate)
	mux.HandleFunc("/local/explain", h.localExplain)
	mux.HandleFunc("/local/models", h.localModels)
	mux.HandleFunc("/local/automation", h.localAutomation)
	mux.HandleFunc("/local/mission", h.localMission)
	mux.HandleFunc("/local/kill", h.localKill)
	mux.HandleFunc("/direct/intent", h.deviceDirectIntent)
	mux.HandleFunc("/direct/complete", h.deviceDirectComplete)
	mux.HandleFunc("/bind", h.bind)
	mux.HandleFunc("/watch", h.watch)
	mux.HandleFunc("/devices", h.devicesList)
	mux.HandleFunc("/revoke", h.revoke)
	mux.HandleFunc("/authorize", deny)
	mux.HandleFunc("/export", deny)
	mux.HandleFunc("/session", deny)
	return httpx.Companion(mux)
}

func writeLocal(w http.ResponseWriter, status int, body map[string]any) {
	b, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "encode", http.StatusInternalServerError)
		return
	}
	low := strings.ToLower(string(b))
	for _, needle := range []string{"private_key", "mnemonic", "session_key", "app-sk-", "hl_secret"} {
		if strings.Contains(low, needle) {
			http.Error(w, "companion_leak", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(b, '\n'))
}

func deny(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "companion_denied", http.StatusForbidden)
}

func (h *Hub) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	listen, err := ListenAddr()
	if err != nil {
		listen = DefaultAddr
	}
	h.researchMu.Lock()
	running := h.job.running
	h.researchMu.Unlock()
	obs.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":               true,
		"service":          "pit-companion",
		"listen":           listen,
		"sign":             false,
		"trade":            false,
		"pairing":          true,
		"version":          version.Number,
		"research_running": running,
		"restart_allowed":  !running,
	})
}

func (h *Hub) pair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Code string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	h.mu.Lock()
	defer h.mu.Unlock()
	if time.Now().After(h.codeExp) || h.code == "" {
		h.rotateLocked(time.Now())
		http.Error(w, "pairing_expired", http.StatusForbidden)
		return
	}
	if normalizeCode(body.Code) != h.code {
		http.Error(w, "pairing_denied", http.StatusForbidden)
		return
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		http.Error(w, "pairing_denied", http.StatusInternalServerError)
		return
	}
	token := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	h.devices = append(h.devices, device{
		Hash:    hex.EncodeToString(sum[:]),
		Created: time.Now().Unix(),
		Label:   "browser",
	})
	h.rotateLocked(time.Now())
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":      true,
		"device":  token,
		"sign":    false,
		"trade":   false,
		"canSign": false,
	})
}

func (h *Hub) bearer(r *http.Request) bool {
	got := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	if got == "" {
		return false
	}
	sum := sha256.Sum256([]byte(strings.TrimSpace(got)))
	want := hex.EncodeToString(sum[:])
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, d := range h.devices {
		if d.Hash == want {
			return true
		}
	}
	return false
}

func (h *Hub) localStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	origin := r.Header.Get("Origin")
	if origin != "" && !httpx.CompanionOriginOK(origin) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	body := map[string]any{"sign": false, "trade": false, "sessionAlive": false, "surface": "desktop", "version": version.Number}
	st, err := cli.Load(h.Dir)
	if err == nil {
		body["network"] = st.Network
		body["wallet"] = st.Wallet
		body["workspace"] = st.WorkspaceID
		body["kill"] = st.Kill
		body["hypothesis"] = cli.LoadHypothesis(h.Dir)
		if sf, serr := cli.LoadSession(h.Dir); serr == nil {
			body["agent"] = sf.AgentAddr
			body["sessionExpires"] = sf.Expires
			now := time.Now().UnixMilli()
			body["sessionAlive"] = session.Alive(sf.Meta().Session(), now) && !st.Kill
			if name, nerr := session.AgentName(st.WorkspaceID); nerr == nil {
				body["agentName"] = name
			}
		}
		if last := cli.LoadLastOrder(h.Dir); last != nil {
			body["lastOrder"] = last
		}
	}
	ms := auto.LoadMission(h.Dir)
	body["mode"] = ms.Mode
	body["missionRunning"] = ms.Running && ms.Mode != auto.ModeManual
	body["missionStop"] = ms.LastStop
	writeLocal(w, http.StatusOK, body)
}

func (h *Hub) localCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	code, exp := h.Code()
	writeLocal(w, http.StatusOK, map[string]any{
		"code":    code,
		"expires": exp.UTC().Format(time.RFC3339),
		"sign":    false,
		"trade":   false,
	})
}

func (h *Hub) localDoctor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	checks := cli.Doctor(h.Dir)
	writeLocal(w, http.StatusOK, map[string]any{
		"checks":  checks,
		"sign":    false,
		"trade":   false,
		"version": version.Number,
	})
}

func desktopOnly(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return false
	}
	return true
}

func writeBindErr(w http.ResponseWriter, err error) {
	code := err.Error()
	status := http.StatusBadRequest
	switch code {
	case "workspace_owned", "network_switch_denied":
		status = http.StatusConflict
	case "unbound":
		status = http.StatusConflict
	}
	writeLocal(w, status, map[string]any{"ok": false, "error": code, "sign": false, "trade": false})
}

func (h *Hub) localInit(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Wallet  string `json:"wallet"`
		Network string `json:"network"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	st, err := cli.Bind(h.Dir, body.Network, body.Wallet)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, cli.PublicBind(st))
}

func (h *Hub) bind(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.bearer(r) {
		http.Error(w, "device_required", http.StatusUnauthorized)
		return
	}
	var body struct {
		Wallet  string `json:"wallet"`
		Network string `json:"network"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	st, err := cli.Bind(h.Dir, body.Network, body.Wallet)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, cli.PublicBind(st))
}

func (h *Hub) localSession(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	sf, _, err := cli.EnsureLocalSession(h.Dir)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, cli.SessionPublic(sf))
}

func (h *Hub) localConnectionPreview(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Coin       string `json:"coin"`
		ReduceOnly bool   `json:"reduceOnly"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.ReduceOnly {
		p, hash, err := cli.BindReduceOnlyClose(h.Dir, body.Coin)
		if err != nil {
			writeBindErr(w, err)
			return
		}
		writeLocal(w, http.StatusOK, cli.ReduceOnlyPublic(p, hash))
		return
	}
	p, hash, err := cli.BindConnectionPreview(h.Dir, body.Coin)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, cli.ConnectionPreviewPublic(p, hash))
}

func (h *Hub) localPolicy(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	st, err := cli.Load(h.Dir)
	if err != nil {
		writeBindErr(w, fmt.Errorf("unbound"))
		return
	}
	p := policy.Default()
	path, err := cli.PinWorkspace(h.Dir, st.WorkspaceID, p)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	if err := cli.CheckPinned(h.Dir, st.WorkspaceID, p); err != nil {
		writeBindErr(w, err)
		return
	}
	hash, err := p.Hash()
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":        true,
		"pinned":    path,
		"hash":      hash,
		"version":   p.Version,
		"workspace": st.WorkspaceID,
		"sign":      false,
		"trade":     false,
	})
}

func (h *Hub) localRevokeSession(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	if err := cli.Logout(h.Dir, false); err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":     true,
		"local":  "session_deleted",
		"sign":   false,
		"trade":  false,
		"order":  false,
		"cancel": false,
	})
}

func (h *Hub) localKill(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		On bool `json:"on"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if err := cli.SetKill(h.Dir, body.On); err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":     true,
		"kill":   body.On,
		"sign":   false,
		"trade":  false,
		"order":  false,
		"cancel": false,
	})
}

func (h *Hub) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.bearer(r) {
		http.Error(w, "device_required", http.StatusUnauthorized)
		return
	}
	body := map[string]any{"sign": false, "trade": false, "paired": true, "sessionAlive": false}
	st, err := cli.Load(h.Dir)
	if err == nil {
		body["workspace"] = st.WorkspaceID
		body["network"] = st.Network
		body["account"] = st.Wallet
		body["kill"] = st.Kill
		if s, err := cli.LiveFromDisk(h.Dir, st.Kill, time.Now().UnixMilli()); err == nil {
			body["sessionAlive"] = session.Alive(s, time.Now().UnixMilli())
			body["agent"] = s.AgentAddr
		}
	}
	writeLocal(w, http.StatusOK, body)
}

func (h *Hub) watch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	net := config.Mainnet
	if n, err := config.ParseNetwork(r.URL.Query().Get("network")); err == nil {
		net = n
	}
	view := watch.EmptyPublic(string(net))
	cands, err := watch.LiveUniverse(hl.New(config.For(net)), watch.PolicyForWatch())
	if err == nil {
		view = watch.Public(cands, string(net))
	}
	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		return
	}
	obs.WriteJSON(w, http.StatusOK, obs.WatchBody(view))
}

func (h *Hub) devicesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.bearer(r) {
		http.Error(w, "device_required", http.StatusUnauthorized)
		return
	}
	h.mu.Lock()
	n := len(h.devices)
	h.mu.Unlock()
	writeLocal(w, http.StatusOK, map[string]any{"count": n, "sign": false})
}

func (h *Hub) revoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Device string `json:"device"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	sum := sha256.Sum256([]byte(strings.TrimSpace(body.Device)))
	want := hex.EncodeToString(sum[:])
	h.mu.Lock()
	defer h.mu.Unlock()
	out := h.devices[:0]
	for _, d := range h.devices {
		if d.Hash != want {
			out = append(out, d)
		}
	}
	h.devices = out
	writeLocal(w, http.StatusOK, map[string]any{"ok": true, "sign": false})
}
