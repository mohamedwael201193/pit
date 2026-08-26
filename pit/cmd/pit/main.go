package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamedwael201193/pit/internal/calib"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/verify"
	"github.com/mohamedwael201193/pit/internal/watch"
	"github.com/mohamedwael201193/pit/internal/workspace"
)

func usage() {
	fmt.Fprint(os.Stderr, `PIT — Private Alpha OS

Commands:
  pit init --network mainnet|testnet --wallet 0x...
  pit login
  pit network
  pit policy
  pit ask --market market.json --book book.json
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
		cmdPolicy()
	case "status":
		cmdStatus()
	case "kill":
		cmdKill()
	case "network":
		cmdNetwork()
	case "opportunities":
		cmdOpportunities()
	case "card":
		cmdCard()
	case "ask":
		cmdAsk(os.Args[2:])
	case "forecast":
		cmdForecast()
	case "preview":
		cmdPreview()
	case "cancel":
		cmdCancel()
	case "resolve":
		cmdResolve()
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

func stateDir() string {
	d, err := cli.DefaultDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return d
}

func cmdNetwork() {
	mn := config.For(config.Mainnet)
	tn := config.For(config.Testnet)
	fmt.Println("MAINNET production")
	fmt.Printf("  chain    %d %s\n", mn.ChainID, mn.RPC)
	fmt.Printf("  hl       %s\n", mn.HLInfo)
	fmt.Printf("  desk     %s\n", mn.DeskID)
	fmt.Println("TESTNET laboratory")
	fmt.Printf("  chain    %d %s\n", tn.ChainID, tn.RPC)
	fmt.Printf("  hl       %s\n", tn.HLInfo)
	fmt.Println("one workspace binds one row. never mix.")
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
	if prev, err := cli.Load(stateDir()); err == nil {
		if err := cli.RefuseNetworkSwitch(prev.Network, string(net)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
	st := workspace.NewStore()
	ws, err := st.Create(addr, net)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := cli.Save(stateDir(), cli.DiskState{
		WorkspaceID: ws.ID,
		Network:     string(ws.Network),
		Wallet:      string(ws.EVM),
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("workspace %s\n", ws.ID)
	fmt.Printf("network   %s\n", ws.Network)
	fmt.Printf("wallet    %s\n", ws.EVM)
	fmt.Println("session   not created — pit will never withdraw")
	_ = session.AllowedActions
}

func cmdStatus() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Println("network: unset until init")
		fmt.Println("session: none")
		fmt.Println("desk: isAuthorized must be true before sealed inference")
		return
	}
	fmt.Printf("workspace %s\n", st.WorkspaceID)
	fmt.Printf("network   %s\n", st.Network)
	fmt.Printf("wallet    %s\n", st.Wallet)
	fmt.Printf("kill      %v\n", st.Kill)
	fmt.Println("session   none on this CLI until desktop or keychain bind")
	fmt.Println("sign      never in the browser")
	fmt.Println("expired  ", cli.ExpiredCopy)
	fmt.Println("revoked  ", cli.RevokedCopy)
}

func cmdKill() {
	if err := cli.SetKill(stateDir(), true); err != nil {
		fmt.Fprintln(os.Stderr, "kill switch requires pit init first")
		os.Exit(2)
	}
	fmt.Println("kill switch: on — signing blocked for this workspace")
}

func cmdPolicy() {
	p := policy.Default()
	for _, c := range policy.Cards(p) {
		fmt.Printf("%s\n  %s\n  %s\n", c.Title, c.Value, c.Law)
	}
	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
	if st, err := cli.Load(stateDir()); err == nil {
		path, err := cli.PinWorkspace(stateDir(), st.WorkspaceID, p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err := cli.CheckPinned(stateDir(), st.WorkspaceID, p); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Println("pinned", path)
	}
}

func cmdOpportunities() {
	p := policy.Default()
	net := config.Mainnet
	if st, err := cli.Load(stateDir()); err == nil {
		if n, err := config.ParseNetwork(st.Network); err == nil {
			net = n
		}
	}
	cands, err := watch.Live(hl.New(config.For(net)), p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(watch.Attention(len(cands)))
	fmt.Println("Watch does not place orders.")
	for _, c := range cands {
		fmt.Printf("%s  %s  mark=%g\n", c.Coin, c.Reason, c.Book.MarkPx)
	}
}

func cmdAsk(args []string) {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "ask requires pit init first")
		os.Exit(2)
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	marketPath, bookPath, err := cli.ParseAskFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ask requires --market and --book files")
		os.Exit(2)
	}
	market, book, err := compute.LoadEnvelope(marketPath, bookPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := compute.ProductAskEnvelope(net, false, compute.LookBin(), market, book); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("sealed ask submitted")
}

func cmdForecast() {
	fmt.Println("host scores forecasts; model size is ignored")
	if _, err := cli.Load(stateDir()); err != nil {
		fmt.Fprintln(os.Stderr, "forecast requires pit init first")
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "forecast requires a live sealed ask and a bound session")
	os.Exit(2)
}

func cmdPreview() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "preview requires pit init first")
		os.Exit(2)
	}
	fmt.Println(cli.PreviewCopy)
	fmt.Println(cli.MutationInvalidates())
	fmt.Fprintf(os.Stderr, "preview requires a live session for workspace %s\n", st.WorkspaceID)
	os.Exit(2)
}

func cmdCancel() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "cancel requires pit init first")
		os.Exit(2)
	}
	fmt.Println("cancel is order/cancel only. It cannot withdraw or change leverage.")
	fmt.Fprintf(os.Stderr, "cancel requires a live session and a bound clientOrderId for workspace %s\n", st.WorkspaceID)
	os.Exit(2)
}

func cmdResolve() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "resolve requires pit init first")
		os.Exit(2)
	}
	fmt.Println("resolve scores host probability. It does not invent a win rate.")
	fmt.Fprintf(os.Stderr, "resolve requires a stored forecast for workspace %s\n", st.WorkspaceID)
	os.Exit(2)
}

func cmdCard() {
	h := calib.Card(nil, 30)
	fmt.Println(h.Copy)
	fmt.Println("Not enough resolved forecasts.")
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
