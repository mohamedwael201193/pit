package companion

import (
	"encoding/json"
	"net/http"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/httpx"
)

func challengeBody(ch compute.Challenge) map[string]any {
	return map[string]any{
		"ok":        true,
		"message":   ch.Message,
		"digest":    ch.Digest,
		"provider":  ch.Provider,
		"model":     ch.Model,
		"network":   ch.Network,
		"wallet":    ch.Wallet,
		"expiresAt": ch.ExpiresAt,
		"tokenId":   ch.TokenId,
		"explain":   ch.Explain,
		"sign":      false,
		"trade":     false,
	}
}

func (h *Hub) localDirectIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	ch, err := cli.IssueDirectChallenge(h.Dir)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, challengeBody(ch))
}

func (h *Hub) localDirectComplete(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Signature string `json:"signature"`
		Message   string `json:"message"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	meta, err := cli.CompleteDirect(h.Dir, body.Message, body.Signature)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":        true,
		"provider":  meta.Provider,
		"model":     meta.Model,
		"expiresAt": meta.ExpiresAt,
		"source":    meta.Source,
		"sign":      false,
		"trade":     false,
	})
}

func (h *Hub) localDirectStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	writeLocal(w, http.StatusOK, cli.DirectStatus(h.Dir))
}

func (h *Hub) localResearch(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Coin string `json:"coin"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	rep, err := cli.RunWorkspaceResearch(h.Dir, body.Coin)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	roles := make([]any, 0, len(rep.Roles))
	for _, role := range rep.Roles {
		roles = append(roles, role)
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":     true,
		"roles":  roles,
		"note":   rep.Note,
		"sign":   false,
		"trade":  false,
		"verify": true,
	})
}

func (h *Hub) deviceDirectIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.bearer(r) {
		http.Error(w, "device_required", http.StatusUnauthorized)
		return
	}
	ch, err := cli.IssueDirectChallenge(h.Dir)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, challengeBody(ch))
}

func (h *Hub) deviceDirectComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.bearer(r) {
		http.Error(w, "device_required", http.StatusUnauthorized)
		return
	}
	var body struct {
		Signature string `json:"signature"`
		Message   string `json:"message"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	meta, err := cli.CompleteDirect(h.Dir, body.Message, body.Signature)
	if err != nil {
		writeBindErr(w, err)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{
		"ok":        true,
		"provider":  meta.Provider,
		"model":     meta.Model,
		"expiresAt": meta.ExpiresAt,
		"source":    meta.Source,
		"sign":      false,
		"trade":     false,
	})
}
