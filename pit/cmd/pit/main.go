package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/calib"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/storage"
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
  pit session
  pit ask --market market.json --book book.json
  pit opportunities
  pit forecast
  pit preview --market ETH --side buy --forecast <id>
  pit authorize --i-understand
  pit cancel
  pit status
  pit resolve
  pit card
  pit verify --preview 0x... --root 0x... --network mainnet --workspace <id>
  pit proof --root 0x... --out file
  pit kill

PIT never asks for a seed phrase or a trading secret.
Session keys stay on this machine.
authorize requires a TTY, the exact word AUTHORIZE, and a live session on this machine.
pit session creates a one-hour order/cancel agent in the local keyring. It never prints the key.
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
	case "session":
		cmdSession()
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
		cmdPreview(os.Args[2:])
	case "cancel":
		cmdCancel()
	case "resolve":
		cmdResolve()
	case "authorize":
		cmdAuthorize(os.Args[2:])
	case "verify":
		cmdVerify(os.Args[2:])
	case "proof":
		cmdProof(os.Args[2:])
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
	fmt.Println("session   not created — pit session then your wallet approveAgent")
	_ = session.AllowedActions
}

func cmdSession() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "session requires pit init first")
		os.Exit(2)
	}
	sf, err := cli.CreateLocalSession(stateDir(), st.WorkspaceID, st.Network, "v1")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	name, err := session.AgentName(st.WorkspaceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("agent     %s\n", sf.AgentAddr)
	fmt.Printf("name      %s\n", name)
	fmt.Printf("expires   %d\n", sf.Expires)
	fmt.Println("ttl       1h")
	fmt.Println("order     allowed")
	fmt.Println("cancel    allowed")
	fmt.Println("withdraw  denied")
	fmt.Println("your wallet must approveAgent this address. PIT never prints the session key.")
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
	if s, err := cli.LiveFromDisk(stateDir(), st.Kill, time.Now().UnixMilli()); err == nil {
		fmt.Printf("session   live %s\n", s.AgentAddr)
		fmt.Printf("expires   %d\n", s.Expires)
		fmt.Println("perms     order yes  cancel yes  withdraw no")
	} else {
		fmt.Println("session   none on this CLI until pit session")
	}
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

func cmdPreview(args []string) {
	coin, side, forecast, err := cli.ParsePreviewFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "preview requires --market --side --forecast")
		os.Exit(2)
	}
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "preview requires pit init first")
		os.Exit(2)
	}
	live, err := cli.LiveFromDisk(stateDir(), st.Kill, time.Now().UnixMilli())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	book, err := hl.New(config.For(net)).PublicBook(strings.ToUpper(coin))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cloid, err := cli.NewCloid()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	p, err := cli.HostPreview(coin, side, forecast, book, policy.Default(), live, time.Now().UTC(), cloid, time.Now().UnixMilli())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	hash, err := cli.SavePreview(stateDir(), p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(cli.PreviewCopy)
	fmt.Println(cli.MutationInvalidates())
	fmt.Printf("hash    %s\n", hash)
	fmt.Printf("market  %s\n", p.Market)
	fmt.Printf("side    %s\n", p.Side)
	fmt.Printf("sz      %g\n", p.Sz)
	fmt.Printf("cloid   %s\n", p.Cloid)
	fmt.Println("model size was ignored")
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
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "authorize requires pit init first")
		os.Exit(2)
	}
	live, liveErr := cli.LiveFromDisk(stateDir(), st.Kill, time.Now().UnixMilli())
	if err := cli.RunAuthorize(true, typed, true, live, time.Now().UnixMilli()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if liveErr != nil {
		fmt.Fprintln(os.Stderr, liveErr)
		os.Exit(2)
	}
	card, hash, err := cli.LoadPreview(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if card.WorkspaceID != live.Workspace || card.SessionID != live.ID {
		fmt.Fprintln(os.Stderr, "wrong_workspace")
		os.Exit(2)
	}
	if err := pitexec.RequirePreview(card, hash, time.Now().UnixMilli()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("order allowed")
	fmt.Println("cancel allowed")
	fmt.Println("withdraw denied")
	fmt.Println("preview bound")
	fmt.Println("PIT did not send an order")
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

func cmdProof(args []string) {
	if err := storage.RefuseMissingProof(storage.LookCLI()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "proof download needs --root, --out, and a workspace key file. PIT does not use a global memory key.")
	os.Exit(2)
}
