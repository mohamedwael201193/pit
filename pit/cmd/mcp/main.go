package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamedwael201193/pit/internal/cli"
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
		dir := ""
		pol := watch.PolicyForWatch()
		if d, err := cli.DefaultDir(); err == nil {
			dir = d
			if st, lerr := cli.Load(dir); lerr == nil {
				if n, perr := config.ParseNetwork(st.Network); perr == nil {
					net = n
				}
				pol = cli.ActivePolicy(dir)
			}
		}
		view := watch.EmptyPublic(string(net))
		cands, err := watch.LiveUniverse(hl.New(config.For(net)), pol)
		if err == nil {
			view = watch.Public(cands, string(net))
			if dir != "" {
				view = cli.SizeWatch(dir, view, pol)
			}
		}
		return map[string]any{
			"count":             view.Count,
			"copy":              view.Copy,
			"coins":             view.Coins,
			"execGate":          view.ExecGate,
			"execWhy":           view.ExecWhy,
			"executionFeasible": view.ExecFeasibleN,
			"routes":            view.Routes,
			"trade":             false,
			"sign":              false,
			"authorize":         false,
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
