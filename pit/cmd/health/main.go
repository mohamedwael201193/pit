package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/obs"
)

func main() {
	if err := config.GuardFallbacks(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}
	if err := config.RefuseSessionEnv(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodHead {
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":        true,
			"service":   "pit",
			"time":      time.Now().UTC().Format(time.RFC3339),
			"sign":      false,
			"requestId": obs.NewRequestID(),
		})
	})
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	_ = http.ListenAndServe(":"+addr, mux)
}
