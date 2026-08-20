package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/mcp"
)

func main() {
	if err := config.GuardFallbacks(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var req mcp.Request
		if err := dec.Decode(&req); err != nil {
			return
		}
		if err := enc.Encode(mcp.Handle(req)); err != nil {
			return
		}
	}
}
