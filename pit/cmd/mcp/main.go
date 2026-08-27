package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/watch"
	"github.com/mohamedwael201193/pit/mcp"
)

func main() {
	if err := config.GuardFallbacks(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	mcp.LiveOpportunities = func() map[string]any {
		net := config.Mainnet
		view := watch.EmptyPublic(string(net))
		cands, err := watch.Live(hl.New(config.For(net)), watch.PolicyForWatch())
		if err == nil {
			view = watch.Public(cands, string(net))
		}
		return map[string]any{
			"count": view.Count,
			"copy":  view.Copy,
			"coins": view.Coins,
			"trade": false,
			"sign":  false,
		}
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
