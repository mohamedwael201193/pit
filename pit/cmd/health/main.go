package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
)

func main() {
	if err := config.GuardFallbacks(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":      true,
			"service": "pit",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"sign":    false,
		})
	})
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	_ = http.ListenAndServe(":"+addr, mux)
}
