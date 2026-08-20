package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/verify"
	"github.com/mohamedwael201193/pit/internal/workspace"
)

func usage() {
	fmt.Fprint(os.Stderr, `PIT — Private Alpha OS

Commands:
  pit init --network mainnet|testnet --wallet 0x...
  pit login
  pit policy
  pit ask
  pit opportunities
  pit forecast
  pit preview
  pit authorize --i-understand
  pit cancel
  pit status
  pit resolve
  pit card
  pit verify --preview 0x... --root 0x... --network mainnet --workspace <id>
  pit kill

PIT never asks for a seed phrase or a trading secret.
Session keys stay on this machine.
authorize requires a TTY and the exact word AUTHORIZE.
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
		fmt.Println("PIT never asks for your private key.")
	case "policy":
		b, _ := json.MarshalIndent(policy.Default(), "", "  ")
		fmt.Println(string(b))
	case "status":
		fmt.Println("network: unset until init")
		fmt.Println("session: none")
		fmt.Println("desk: isAuthorized must be true before sealed inference")
	case "kill":
		fmt.Println("kill switch: set locally — signing blocked for this workspace")
	case "ask", "opportunities", "forecast", "preview", "cancel", "resolve", "card":
		fmt.Fprintf(os.Stderr, "%s requires a bound workspace and a live session. Run pit init first.\n", os.Args[1])
		os.Exit(2)
	case "authorize":
		cmdAuthorize(os.Args[2:])
	case "verify":
		cmdVerify(os.Args[2:])
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

func stdinTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func cmdAuthorize(args []string) {
	iUnderstand := false
	for _, a := range args {
		if a == "--i-understand" {
			iUnderstand = true
		}
	}
	if !stdinTTY() {
		fmt.Fprintln(os.Stderr, "authorize refused: piped confirmation is not allowed")
		os.Exit(2)
	}
	if !iUnderstand {
		fmt.Fprintln(os.Stderr, "authorize refused: pass --i-understand and type AUTHORIZE")
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "Type AUTHORIZE to sign the exact preview (order or cancel only):")
	var typed string
	_, _ = fmt.Fscanln(os.Stdin, &typed)
	if err := cli.ConfirmAuthorize(true, typed, true); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "authorize requires a bound workspace and a live session. Run pit init first.")
	os.Exit(2)
}

func cmdVerify(args []string) {
	var preview, root, network, workspaceID string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--preview":
			i++
			if i < len(args) {
				preview = args[i]
			}
		case "--root":
			i++
			if i < len(args) {
				root = args[i]
			}
		case "--network":
			i++
			if i < len(args) {
				network = args[i]
			}
		case "--workspace":
			i++
			if i < len(args) {
				workspaceID = args[i]
			}
		}
	}
	if err := verify.Check(verify.Receipt{
		PreviewHash: preview,
		StorageRoot: root,
		Network:     network,
		Workspace:   workspaceID,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("receipt fields well-formed")
	fmt.Println("recompute the storage proof with the official Go client --proof")
}
