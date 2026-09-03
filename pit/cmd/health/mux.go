package main

import (
	"net/http"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/httpx"
	"github.com/mohamedwael201193/pit/internal/obs"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Method == http.MethodHead {
			w.Header().Set("Content-Type", "application/json")
			return
		}
		obs.WriteJSON(w, http.StatusOK, obs.HealthBody(obs.NewRequestID()))
	})
	mux.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
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
	})
	mux.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Method == http.MethodHead {
			w.Header().Set("Content-Type", "application/json")
			return
		}
		rel, ok := cachedRelease()
		if !ok {
			obs.WriteJSON(w, http.StatusOK, releaseBody(publicRelease{Unsigned: true}))
			return
		}
		obs.WriteJSON(w, http.StatusOK, releaseBody(rel))
	})
	mux.HandleFunc("/windows", func(w http.ResponseWriter, r *http.Request) {
		rel, _ := cachedRelease()
		redirectLatestAsset(w, r, windowsAsset(rel), firstNonEmpty(rel.File, "PIT_x64-setup.exe"))
	})
	mux.HandleFunc("/checksums", func(w http.ResponseWriter, r *http.Request) {
		rel, _ := cachedRelease()
		redirectLatestAsset(w, r, checksumsAsset(rel), "SHA256SUMS.txt")
	})
	mux.HandleFunc("/.well-known/agent-card.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		if r.Method == http.MethodHead {
			return
		}
		_, _ = w.Write([]byte(agentCardJSON))
	})
	return httpx.Public(mux)
}

func firstNonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}
