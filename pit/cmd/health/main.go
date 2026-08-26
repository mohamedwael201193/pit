package main

import (
	"net/http"
	"os"

	"github.com/mohamedwael201193/pit/internal/config"
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
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	_ = http.ListenAndServe(":"+addr, newMux())
}
