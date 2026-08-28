package main

import (
	"net/http"

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
	return httpx.Public(mux)
}
