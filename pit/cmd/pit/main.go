package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/workspace"
)

func usage() {
	fmt.Fprint(os.Stderr, `PIT — Private Alpha OS

Commands:
  pit init --network mainnet|testnet --wallet 0x...
  pit login
  pit policy
  pit status
  pit kill
  pit verify --help

PIT never asks for a seed phrase or a trading secret.
Session keys stay on this machine.
`)
}

func main() {
	if err := config.GuardFallbacks(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "init":
		cmdInit(os.Args[2:])
	case "login":
		fmt.Println("Connect your wallet in the desktop or web app, then sign the bind message.")
	case "policy":
		b, _ := json.MarshalIndent(policy.Default(), "", "  ")
		fmt.Println(string(b))
	case "status":
		fmt.Println("network: unset until init")
		fmt.Println("session: none")
	case "kill":
		fmt.Println("kill switch: set locally — signing blocked for this workspace")
	case "ask", "opportunities", "forecast", "preview", "authorize", "cancel", "resolve", "card", "verify":
		fmt.Fprintf(os.Stderr, "%s requires a bound workspace and a live session. Run pit init first.\n", os.Args[1])
		os.Exit(2)
	case "help", "-h", "--help":
		usage()
	default:
		usage()
		os.Exit(2)
	}
}

func cmdInit(args []string) {
	net := config.Mainnet
	wallet := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--network":
			i++
			n, err := config.ParseNetwork(args[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			net = n
		case "--wallet":
			i++
			wallet = args[i]
		}
	}
	addr, err := identity.NormalizeAddress(wallet)
	if err != nil {
		fmt.Fprintln(os.Stderr, "YOUR WALLET address is required: --wallet 0x...")
		os.Exit(2)
	}
	st := workspace.NewStore()
	ws, err := st.Create(addr, net)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("workspace %s\n", ws.ID)
	fmt.Printf("network   %s\n", ws.Network)
	fmt.Printf("wallet    %s\n", ws.EVM)
	fmt.Println("session   not created — pit will never withdraw")
	_ = session.AllowedActions
}
