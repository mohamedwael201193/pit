package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/calib"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/companion"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/ledger"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/storage"
	"github.com/mohamedwael201193/pit/internal/verify"
	"github.com/mohamedwael201193/pit/internal/version"
	"github.com/mohamedwael201193/pit/internal/watch"
	"github.com/mohamedwael201193/pit/internal/workspace"
)

var asJSON bool

func usage() {
	fmt.Fprint(os.Stderr, `PIT — Private Alpha OS

Commands:
  pit init --network mainnet|testnet --wallet 0x...
  pit login
  pit wallet
  pit network
  pit policy
  pit session
  pit companion
  pit ask --market market.json --book book.json
  pit watch
  pit opportunities
  pit forecast
  pit preview --market ETH --side buy --forecast <id>
  pit authorize --i-understand
  pit orders
  pit cancel
  pit status
  pit resolve
  pit card
  pit verify --preview 0x... --root 0x... --network mainnet --workspace <id>
  pit proof --root 0x... --out file --key-file key.hex
  pit kill
  pit doctor
  pit logout [--forget]
  pit revoke
  pit version

Every command accepts --json. Exit 0 on success, 2 on usage, 1 on failed doctor.

PIT never asks for a seed phrase or a trading secret.
Session keys stay on this machine (OS keychain unless PIT_KEYRING=file).
authorize requires a TTY, the exact word AUTHORIZE, and a live session on this machine.
pit session creates a one-hour order/cancel agent. It never prints the key.
pit companion listens on 127.0.0.1 only. Pairing does not send the session key to the browser.
`)
}

func main() {
	if err := config.GuardFallbacks(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	rest, args := []string{}, os.Args[1:]
	asJSON, rest = cli.WantJSON(args)
	if len(rest) < 1 {
		usage()
		os.Exit(2)
	}
	switch rest[0] {
	case "init":
		cmdInit(rest[1:])
	case "login":
		fmt.Println("Connect your wallet in the desktop or web app, then sign the bind message.")
		fmt.Println("PIT never asks for your private key.")
	case "wallet":
		cmdWallet()
	case "policy":
		cmdPolicy()
	case "session":
		cmdSession()
	case "companion":
		cmdCompanion()
	case "status":
		cmdStatus()
	case "kill":
		cmdKill()
	case "network":
		cmdNetwork()
	case "watch", "opportunities":
		cmdOpportunities()
	case "orders":
		cmdOrders()
	case "card":
		cmdCard()
	case "ask":
		cmdAsk(rest[1:])
	case "forecast":
		cmdForecast()
	case "preview":
		cmdPreview(rest[1:])
	case "cancel":
		cmdCancel()
	case "resolve":
		cmdResolve()
	case "authorize":
		cmdAuthorize(rest[1:])
	case "verify":
		cmdVerify(rest[1:])
	case "proof":
		cmdProof(rest[1:])
	case "doctor":
		cmdDoctor()
	case "logout":
		cmdLogout(rest[1:])
	case "revoke":
		cmdRevoke()
	case "version", "-v", "--version":
		cmdVersion()
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
		linked, linkErr := cli.LiveLinked(st.Network, st.Wallet, s.Workspace, s.AgentAddr, time.Now().UnixMilli())
		fmt.Println(cli.LinkCopy(linked, linkErr))
	} else {
		fmt.Println("session   none on this CLI until pit session")
	}
	fmt.Println("sign      never in the browser")
	fmt.Println("expired  ", cli.ExpiredCopy)
	fmt.Println("revoked  ", cli.RevokedCopy)
	if p, h, err := cli.LoadPreview(stateDir()); err == nil {
		fmt.Printf("preview   %s\n", h)
		if rec, err := cli.LookupAction(stateDir(), st.Network, st.WorkspaceID, p.Cloid); err == nil {
			fmt.Printf("ledger    %s\n", rec.Status)
		} else {
			fmt.Println("ledger    none")
		}
		found, vErr := cli.LiveOnVenue(st.Network, st.Wallet, p.Cloid)
		fmt.Println(cli.VenueCopy(found, vErr))
	}
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
	live, err := cli.LiveFromDisk(stateDir(), st.Kill, time.Now().UnixMilli())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	card, hash, err := cli.LoadPreview(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if card.WorkspaceID != live.Workspace {
		fmt.Fprintln(os.Stderr, "wrong_workspace")
		os.Exit(2)
	}
	rec, err := cli.LookupAction(stateDir(), st.Network, live.Workspace, card.Cloid)
	if err != nil || rec.Status != ledger.StatusAuthorized {
		fmt.Fprintln(os.Stderr, "cancel requires an authorized preview")
		os.Exit(2)
	}
	coin, err := pitexec.CoinFromMarket(card.Market)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	book, err := cli.LiveAsset(st.Network, coin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	raw, err := cli.CancelWire(book.Asset, card.Cloid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	found, qerr := cli.LiveOnVenue(st.Network, st.Wallet, card.Cloid)
	if qerr != nil {
		fmt.Fprintln(os.Stderr, qerr)
		if err := pitexec.QueryBeforeRetry(false, rec.OID, ""); err != nil {
			fmt.Println(err)
		}
	} else if err := pitexec.NeedOnVenue(found); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("on venue")
	}
	fmt.Println("cancel is order/cancel only. It cannot withdraw or change leverage.")
	fmt.Printf("cloid    %s\n", card.Cloid)
	fmt.Printf("asset    %d\n", book.Asset)
	fmt.Printf("preview  %s\n", hash)
	env, signErr := cli.SignBound(stateDir(), live, st.Network, raw, time.Now().UnixMilli())
	if signErr != nil || !env.Signed() {
		if err := pitexec.RefuseUnsigned(false); err != nil {
			fmt.Println(err)
		}
		fmt.Println("PIT did not send a cancel")
		return
	}
	fmt.Println("signed locally")
	linked, linkErr := cli.LiveLinked(st.Network, st.Wallet, live.Workspace, live.AgentAddr, time.Now().UnixMilli())
	if linkErr != nil {
		fmt.Fprintln(os.Stderr, linkErr)
		linked = false
	}
	if linked {
		fmt.Println("agent linked")
	}
	if err := pitexec.RefusePostUntilLinked(linked); err != nil {
		fmt.Println(err)
		fmt.Println("PIT did not send a cancel")
		return
	}
	if qerr != nil || !found {
		fmt.Println("PIT did not send a cancel")
		return
	}
	body, err := cli.PostLinked(st.Network, env, linked, hash)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("PIT did not send a cancel")
		os.Exit(2)
	}
	oid := pitexec.ReceiptOID(body)
	if err := cli.RememberPosted(stateDir(), st.Network, live.Workspace, card.Cloid, oid); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("posted")
	if oid != "" {
		fmt.Printf("oid      %s\n", oid)
	}
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
	h := pitexec.HashForAuthorize(hash)
	used := map[string]struct{}{}
	if err := pitexec.Prepare(pitexec.Intent{
		Action:    "order",
		Preview:   card,
		Hash:      h,
		Workspace: live.Workspace,
	}, time.Now().UnixMilli(), used); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	coin, err := pitexec.CoinFromMarket(card.Market)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	book, err := cli.LiveAsset(st.Network, coin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	raw, err := pitexec.WireFromPreview(card, book.Asset)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := cli.RememberAuthorized(stateDir(), st.Network, live.Workspace, card.Cloid, h); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("order allowed")
	fmt.Println("cancel allowed")
	fmt.Println("withdraw denied")
	fmt.Println("preview bound")
	fmt.Println("authorized")
	fmt.Println("ledger    recorded")
	fmt.Printf("asset    %d\n", book.Asset)
	env, signErr := cli.SignBound(stateDir(), live, st.Network, raw, time.Now().UnixMilli())
	if signErr != nil || !env.Signed() {
		if err := pitexec.RefuseUnsigned(false); err != nil {
			fmt.Println(err)
		}
		fmt.Println("PIT did not send an order")
		return
	}
	fmt.Println("signed locally")
	linked, linkErr := cli.LiveLinked(st.Network, st.Wallet, live.Workspace, live.AgentAddr, time.Now().UnixMilli())
	if linkErr != nil {
		fmt.Fprintln(os.Stderr, linkErr)
		linked = false
	}
	if linked {
		fmt.Println("agent linked")
	}
	if err := pitexec.RefusePostUntilLinked(linked); err != nil {
		fmt.Println(err)
		fmt.Println("PIT did not send an order")
		return
	}
	body, err := cli.PostLinked(st.Network, env, linked, h)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("PIT did not send an order")
		os.Exit(2)
	}
	oid := pitexec.ReceiptOID(body)
	if err := cli.RememberPosted(stateDir(), st.Network, live.Workspace, card.Cloid, oid); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("posted")
	if oid != "" {
		fmt.Printf("oid      %s\n", oid)
	}
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
	flags, err := cli.ParseProofFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := storage.RefuseMissingProof(storage.LookCLI()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	key, err := cli.LoadProofKey(flags.KeyFile, os.Getenv("PIT_MEMORY_KEY"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "proof requires pit init first")
		os.Exit(2)
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	ch := config.For(net)
	job, err := storage.ProofJob(storage.LookCLI(), ch.RPC, ch.StorageIndexer, key, flags.Root, flags.Out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	argv, err := storage.DownloadArgs(job)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("download", strings.Join(storage.RedactArgs(argv), " "))
	cmd := storage.Command(job, argv)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmdWallet() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "wallet unbound until pit init")
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"wallet": st.Wallet, "network": st.Network, "workspace": st.WorkspaceID, "sign": false,
		})
		return
	}
	fmt.Println("wallet   ", st.Wallet)
	fmt.Println("network  ", st.Network)
	fmt.Println("workspace", st.WorkspaceID)
	fmt.Println("PIT never asks for a seed phrase.")
}

func cmdOrders() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "orders require pit init first")
		os.Exit(2)
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	raw, err := hl.New(config.For(net)).OpenOrders(st.Wallet)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"orders": json.RawMessage(raw), "sign": false, "trade": false})
		return
	}
	fmt.Println("open orders for", st.Wallet)
	fmt.Println(string(raw))
}

func cmdDoctor() {
	checks := cli.Doctor(stateDir())
	cli.PrintDoctor(os.Stdout, checks, asJSON)
	if cli.DoctorFailed(checks) {
		os.Exit(1)
	}
}

func cmdVersion() {
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "version": version.Number, "name": version.Name, "sign": false})
		return
	}
	fmt.Println(version.String())
}

func cmdLogout(args []string) {
	forget := false
	for _, a := range args {
		if a == "--forget" {
			forget = true
		}
	}
	if err := cli.Logout(stateDir(), forget); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "forget": forget, "sign": false})
		return
	}
	if forget {
		fmt.Println("workspace unbound. Session secrets deleted.")
		return
	}
	fmt.Println("session deleted. Workspace bind remains. Re-run pit session to mint a new agent.")
}

func cmdRevoke() {
	if err := cli.Logout(stateDir(), false); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"ok":    true,
			"local": "session_deleted",
			"sign":  false,
			"trade": false,
			"next":  "remove the PIT agent from your Hyperliquid account",
		})
		return
	}
	fmt.Println("local session deleted. The key is gone from this machine.")
	fmt.Println("Remove the PIT agent from your Hyperliquid account. PIT cannot withdraw and cannot call approveAgent from a session.")
}

func cmdCompanion() {
	addr, err := companion.ListenAddr()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	h := companion.New(stateDir())
	code, exp := h.Code()
	pretty := code[:4] + "-" + code[4:]
	fmt.Printf("PIT companion on %s\n", addr)
	fmt.Printf("Pairing code %s (expires %s)\n", pretty, exp.Format(time.RFC3339))
	fmt.Println("Type this code at https://pit0g.vercel.app/pair")
	fmt.Println("Session keys stay on this machine. The browser never receives them.")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := http.Serve(ln, h.Handler()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
