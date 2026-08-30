package companion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/demo"
	"github.com/mohamedwael201193/pit/internal/deskcmd"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/proof"
)

func chatModelPath(dir string) string {
	return filepath.Join(dir, "chat-model.json")
}

func loadChatModel(dir string) string {
	raw, err := os.ReadFile(chatModelPath(dir))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return "host-parsed"
	}
	var body struct {
		Model string `json:"model"`
	}
	if json.Unmarshal(raw, &body) != nil || strings.TrimSpace(body.Model) == "" {
		return "host-parsed"
	}
	return strings.TrimSpace(body.Model)
}

func saveChatModel(dir, model string) error {
	if dir == "" {
		return fmt.Errorf("unbound")
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = "host-parsed"
	}
	b, err := json.MarshalIndent(map[string]any{"model": model, "sign": false, "trade": false}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(chatModelPath(dir), append(b, '\n'), 0o600)
}

func (h *Hub) chatModelHonesty(model string) string {
	m := strings.TrimSpace(model)
	if m == "" || m == "host-parsed" {
		return ""
	}
	net := config.Mainnet
	if st, err := cli.Load(h.Dir); err == nil && strings.TrimSpace(st.Network) != "" {
		if n, nerr := config.ParseNetwork(st.Network); nerr == nil {
			net = n
		}
	}
	ok, why := compute.CatalogUsableForChat(m, net)
	if ok && why == "direct_teeml_research_sku" {
		return m + " is the Direct TeeML research SKU. This thread stayed host-parsed. Chat cannot AUTHORIZE, pin policy, or place an order. PIT did not switch to Router."
	}
	return m + " is listed on the official catalog but is not Direct on this workspace (" + why + "). Desk commands stayed host-parsed. PIT will not call Router and will not send the private book there."
}

func (h *Hub) localPolicyPreview(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	before := cli.ActivePolicy(h.Dir)
	draft := decodePolicyBody(r, before)
	draft = policy.Clamp(draft)
	allow, refuse := policy.AllowedRefused(draft)
	lines := append(policy.Diff(before, draft), h.opportunityConsequences(draft)...)
	body := h.policyPublic(lines)
	body["allowed"] = allow
	body["refused"] = refuse
	body["pinned"] = false
	body["previewOnly"] = true
	writeLocal(w, http.StatusOK, body)
}

func (h *Hub) localDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !desktopOnly(w, r) {
			return
		}
		var body struct {
			Mode string `json:"mode"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if err := demo.SavePref(h.Dir, body.Mode); err != nil {
			writeBindErr(w, err)
			return
		}
		writeLocal(w, http.StatusOK, map[string]any{"ok": true, "mode": demo.LoadPref(h.Dir), "live": demo.LoadPref(h.Dir) == "live", "sign": false, "trade": false})
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	rows, _ := proof.Index(h.proofDir(), 20)
	receipts := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		b, err := json.Marshal(row)
		if err != nil {
			continue
		}
		var m map[string]any
		if json.Unmarshal(b, &m) != nil {
			continue
		}
		receipts = append(receipts, m)
	}
	events := make([]map[string]any, 0)
	for _, ev := range readActivity(h.Dir, 40) {
		b, err := json.Marshal(ev)
		if err != nil {
			continue
		}
		var m map[string]any
		if json.Unmarshal(b, &m) == nil {
			events = append(events, m)
		}
	}
	view := demo.Live(h.Dir)
	pref := demo.LoadPref(h.Dir)
	if pref == "replay" {
		view = demo.Replay(h.Dir, receipts, events)
	}
	out, _ := json.Marshal(view)
	var shaped map[string]any
	_ = json.Unmarshal(out, &shaped)
	shaped["pref"] = pref
	writeLocal(w, http.StatusOK, shaped)
}

func (h *Hub) localChatStream(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Text   string `json:"text"`
		Thread string `json:"thread"`
		Model  string `json:"model"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if secretful(body.Text) {
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "secret_refused", "sign": false, "trade": false})
		return
	}
	parsed := h.decorateChat(deskcmd.Parse(body.Text))
	thread := strings.TrimSpace(body.Thread)
	if thread == "" {
		thread = "desk"
	}
	appendChatThread(h.Dir, "user", body.Text, "", thread)
	if parsed.Agent == nil {
		if note := h.chatModelHonesty(firstNonEmpty(body.Model, loadChatModel(h.Dir))); note != "" {
			parsed.Reply = parsed.Reply + "\n\n" + note
		}
	}
	if parsed.StartResearch && parsed.Coin == "" {
		parsed.Coin = h.pickBestCoin()
	}
	if parsed.Hypothesis != "" {
		_ = cli.SaveHypothesis(h.Dir, parsed.Hypothesis)
	}
	appendChatThread(h.Dir, "pit", parsed.Reply, parsed.Tool, thread)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, _ := w.(http.Flusher)
	runes := []rune(parsed.Reply)
	for i := 0; i < len(runes); i += 12 {
		end := i + 12
		if end > len(runes) {
			end = len(runes)
		}
		chunk, _ := json.Marshal(map[string]any{"delta": string(runes[i:end]), "sign": false, "trade": false})
		_, _ = fmt.Fprintf(w, "data: %s\n\n", chunk)
		if flusher != nil {
			flusher.Flush()
		}
	}
	done, _ := json.Marshal(chatDonePayload(parsed, thread, "host-parsed", true))
	_, _ = fmt.Fprintf(w, "data: %s\n\n", done)
	if flusher != nil {
		flusher.Flush()
	}
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}
