package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
	"github.com/mohamedwael201193/pit/internal/calib"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/companion"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskcmd"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/ledger"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
	"github.com/mohamedwael201193/pit/internal/storage"
	"github.com/mohamedwael201193/pit/internal/verify"
	"github.com/mohamedwael201193/pit/internal/version"
	"github.com/mohamedwael201193/pit/internal/watch"
	"github.com/mohamedwael201193/pit/mcp"
)

var asJSON bool

func usage() {
	fmt.Fprint(os.Stderr, `PIT — Private Alpha OS

Workspace
  pit init --network mainnet|testnet --wallet 0x...
  pit login
  pit wallet
  pit pair
  pit network
  pit logout [--forget]
  pit revoke

Policy
  pit policy
  pit kill

Session
  pit session
  pit companion
  pit approve
  pit hyperliquid
  pit agent
  pit compute

Research
  pit watch
  pit scan
  pit mission
  pit opportunities
  pit chat "what is happening?"
  pit positions
  pit health
  pit activity
  pit receipt
  pit direct
  pit direct --sig 0x...
  pit ask --market market.json --book book.json
  pit research [ETH|BTC] [--hypothesis none|long|short]
  pit research --market market.json --book book.json
  pit forecast
  pit calibration
  pit preview --market ETH --side buy --forecast <id>
  pit authorize --i-understand
  pit execute --i-understand

Orders
  pit orders
  pit cancel
  pit status
  pit resolve
  pit card

Proof
  pit verify --preview 0x... --root 0x... --network mainnet --workspace <id>
  pit proof --root 0x... --out file --key-file key.hex

System
  pit doctor
  pit security
  pit activity
  pit receipt
  pit mcp
  pit update
  pit version

Every command accepts --json. Exit 0 on success, 2 on usage, 1 on failed doctor.

PIT never asks for a seed phrase or a trading secret.
Session keys stay on this machine (OS keychain unless PIT_KEYRING=file).
authorize requires a TTY, the exact word AUTHORIZE, and a live session on this machine.
pit session creates a 24-hour order/cancel agent. If Hyperliquid still lists that PIT agent, PIT reuses it instead of minting a new address.
pit companion listens on 127.0.0.1 only. Pairing does not send the session key to the browser.
The desktop can bind, pin policy, and mint a session without a terminal.
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
	case "hyperliquid", "agent":
		cmdHyperliquid(rest[1:])
	case "compute":
		cmdCompute()
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
	case "scan":
		cmdScan()
	case "mission", "mode":
		cmdMission()
	case "chat":
		cmdChat(rest[1:])
	case "positions":
		cmdPositions()
	case "health":
		cmdHealth()
	case "orders":
		cmdOrders()
	case "card":
		cmdCard()
	case "ask":
		cmdAsk(rest[1:])
	case "research":
		cmdResearch(rest[1:])
	case "pair":
		cmdPair()
	case "approve":
		cmdApprove()
	case "execute":
		cmdAuthorize(rest[1:])
	case "calibration":
		cmdForecast()
	case "direct":
		cmdDirect(rest[1:])
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
	case "evidence":
		cmdEvidence(rest[1:])
	case "doctor", "security":
		cmdDoctor()
	case "mcp":
		cmdMCP()
	case "activity":
		cmdStatus()
	case "receipt":
		cmdStatus()
	case "update":
		cmdUpdate()
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
	st, err := cli.Bind(stateDir(), string(net), wallet)
	if err != nil {
		if err.Error() == "wallet_required" {
			fmt.Fprintln(os.Stderr, "YOUR WALLET address is required: --wallet 0x...")
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(cli.PublicBind(st))
		return
	}
	fmt.Printf("workspace %s\n", st.WorkspaceID)
	fmt.Printf("network   %s\n", st.Network)
	fmt.Printf("wallet    %s\n", st.Wallet)
	fmt.Println("session   not created — pit session then your wallet approveAgent")
	_ = session.AllowedActions
}

func cmdSession() {
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "session requires pit init first")
		os.Exit(2)
	}
	sf, _, err := cli.EnsureLocalSession(stateDir())
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
	fmt.Println("ttl       24h (reused while Hyperliquid lists this PIT agent)")
	fmt.Println("order     allowed")
	fmt.Println("cancel    allowed")
	fmt.Println("withdraw  denied")
	fmt.Println("your wallet must approveAgent this address. PIT never prints the session key.")
}

func cmdMCP() {
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"allowed": mcp.AllowedTools, "forbidden": mcp.ForbiddenTools,
			"sign": false, "trade": false, "authorize": false, "withdraw": false,
		})
		return
	}
	fmt.Println("MCP is read-only. Trade tools are denied.")
	fmt.Println(mcp.Schema())
	fmt.Println("Use the pit-mcp binary. PIT Desktop remains the signer.")
}

func cmdStatus() {
	st, err := cli.Load(stateDir())
	if err != nil {
		if asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": false, "error": "unbound", "sign": false, "trade": false})
			return
		}
		fmt.Println("network: unset until init")
		fmt.Println("session: none")
		fmt.Println("desk: isAuthorized must be true before sealed inference")
		return
	}
	if asJSON {
		body := map[string]any{"workspace": st.WorkspaceID, "network": st.Network, "wallet": st.Wallet, "kill": st.Kill, "sign": false, "trade": false}
		if last := cli.LoadLastOrder(stateDir()); last != nil {
			body["lastOrder"] = last
		}
		_ = json.NewEncoder(os.Stdout).Encode(body)
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
	if last := cli.LoadLastOrder(stateDir()); last != nil {
		fmt.Printf("last oid  %v\n", last["oid"])
		fmt.Printf("status    %v\n", last["status"])
		if last["cancelled"] == true {
			fmt.Println("cancel    recorded")
		}
		if last["status"] == "filled" {
			fmt.Println("filled orders are positions. cancel does not apply. flatten only with a reduce-only close that YOU authorize.")
		}
	} else {
		fmt.Println("last oid  none until Hyperliquid accepts an order after AUTHORIZE")
	}
	if net, nerr := config.ParseNetwork(st.Network); nerr == nil {
		if pos, perr := hl.New(config.For(net)).Positions(st.Wallet); perr == nil {
			if len(pos) == 0 {
				fmt.Println("positions none")
			}
			for _, p := range pos {
				fmt.Printf("position  %s sz %s entry %s\n", p.Coin, p.Sz, p.EntryPx)
			}
		}
	}
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
	all, err := watch.LiveUniverse(hl.New(config.For(net)), p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	view := watch.Public(all, string(net))
	fmt.Println(view.Copy)
	fmt.Printf("scanned %d live Hyperliquid perps\n", view.Scanned)
	if view.Best != nil {
		fmt.Printf("best %s mark=%g policy=%s\n", view.Best.Coin, view.Best.Mark, view.Best.PolicyFit)
		fmt.Println(view.BestWhy)
	}
	fmt.Println("Markets does not place orders.")
	shown := 0
	for _, c := range all {
		if !c.Eligible {
			continue
		}
		fmt.Printf("%s  %s  mark=%g  fit=PASS\n", c.Coin, c.Reason, c.Book.MarkPx)
		shown++
		if shown >= 12 {
			break
		}
	}
}

func cmdScan() {
	cmdOpportunities()
}

func cmdMission() {
	dir := stateDir()
	m := auto.LoadMission(dir)
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(auto.Public(dir))
		return
	}
	fmt.Println("mode", m.Mode)
	fmt.Println("running", m.Running)
	fmt.Println("best", m.BestCoin)
	fmt.Println("last_action", m.LastAction)
	fmt.Println("last_stop", m.LastStop)
	fmt.Println("trades_today", m.TradesToday)
	fmt.Println("Chat cannot enable Guarded Autonomy.")
}

func cmdChat(args []string) {
	text := strings.TrimSpace(strings.Join(args, " "))
	if text == "" {
		fmt.Fprintln(os.Stderr, `pit chat "what is happening?"`)
		os.Exit(2)
	}
	payload, _ := json.Marshal(map[string]string{"text": text})
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:17373/local/chat", strings.NewReader(string(payload)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r := deskcmd.Parse(text)
		if asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
				"ok": true, "reply": r.Reply, "tool": r.Tool, "coin": r.Coin,
				"start_research": r.StartResearch, "execute": false, "sign": false, "trade": false,
				"companion": false,
			})
			return
		}
		fmt.Println(r.Reply)
		fmt.Println("Companion not running. This is host parse only, not live desk state.")
		return
	}
	defer resp.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(body)
		return
	}
	if reply, _ := body["reply"].(string); reply != "" {
		fmt.Println(reply)
		return
	}
	fmt.Println(body["error"])
}

func cmdPositions() {
	st, err := cli.Load(stateDir())
	if err != nil {
		if asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": false, "error": "unbound", "positions": []any{}, "sign": false, "trade": false})
			return
		}
		fmt.Println("unbound")
		return
	}
	net, nerr := config.ParseNetwork(st.Network)
	if nerr != nil {
		fmt.Fprintln(os.Stderr, nerr)
		os.Exit(2)
	}
	rows, acct, perr := hl.New(config.For(net)).Clearinghouse(st.Wallet)
	if perr != nil {
		fmt.Fprintln(os.Stderr, "HYPERLIQUID_OUTAGE")
		os.Exit(1)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"ok": true, "account": st.Wallet, "equity": acct.AccountValue, "available": acct.Withdrawable,
			"exposure": acct.TotalNtlPos, "positions": rows, "sign": false, "trade": false,
		})
		return
	}
	fmt.Println("account", st.Wallet, "(master, not PIT agent)")
	fmt.Println("equity", acct.AccountValue)
	fmt.Println("available", acct.Withdrawable)
	fmt.Println("exposure", acct.TotalNtlPos)
	if len(rows) == 0 {
		fmt.Println("positions none")
		return
	}
	for _, p := range rows {
		fmt.Printf("position  %s sz %s entry %s uPnL %s\n", p.Coin, p.Sz, p.EntryPx, p.UnrealizedPnl)
	}
}

func cmdHealth() {
	resp, err := http.Get("http://127.0.0.1:17373/health")
	if err != nil {
		if asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": false, "error": "companion_down", "version": version.Number, "sign": false, "trade": false})
			os.Exit(1)
		}
		fmt.Println("companion down")
		fmt.Println("cli", version.String())
		os.Exit(1)
	}
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(body)
		return
	}
	fmt.Println("companion", body["version"], "research_running", body["research_running"])
}

func cmdResearch(args []string) {
	coin := "ETH"
	hyp := ""
	hasFiles := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--market" || a == "--book" {
			hasFiles = true
			break
		}
		if a == "--hypothesis" && i+1 < len(args) {
			hyp = args[i+1]
			i++
			continue
		}
		if !strings.HasPrefix(a, "-") && a != "" {
			coin = a
		}
	}
	if hasFiles {
		cmdAsk(args)
		return
	}
	if hyp != "" {
		if err := cli.SaveHypothesis(stateDir(), hyp); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
	rep, err := cli.RunWorkspaceResearchStage(stateDir(), coin, func(s string) {
		if !asJSON {
			fmt.Fprintln(os.Stderr, s)
		}
	}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"ok":           true,
			"roles":        rep.Roles,
			"note":         rep.Note,
			"sign":         false,
			"trade":        false,
			"verify":       true,
			"eligible":     rep.Eligible,
			"deny":         rep.Deny,
			"preview":      rep.Preview,
			"preview_hash": rep.PreviewHash,
			"hypothesis":   cli.LoadHypothesis(stateDir()),
		})
		return
	}
	fmt.Println(rep.Note)
	for _, role := range rep.Roles {
		fmt.Printf("%v  %v  %v  %v\n", role["role"], role["verify_e2ee"], role["proposed_side"], role["pubkey_signer"])
	}
	if rep.Eligible {
		fmt.Println("preview ready. type AUTHORIZE on the desktop or pit authorize --i-understand")
		return
	}
	if rep.Deny == "no_side" {
		fmt.Println("stand down. committee proposed no side. this is a verified result, not a crash.")
	}
}

func cmdPair() {
	fmt.Println("Launch PIT Desktop on this computer.")
	fmt.Println("Type the one-time pairing code at https://pit0g.vercel.app/pair")
	fmt.Println("The website never receives a session key.")
}

func cmdApprove() {
	cmdHyperliquid(nil)
}

func cmdHyperliquid(args []string) {
	prepare := false
	for _, a := range args {
		if a == "--prepare" {
			prepare = true
		}
	}
	st, err := cli.Load(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, "hyperliquid requires pit init first")
		os.Exit(2)
	}
	sf, serr := cli.LoadSession(stateDir())
	name, _ := session.AgentName(st.WorkspaceID)
	now := time.Now().UnixMilli()
	linked := false
	var until int64
	var linkErr error
	if serr == nil {
		linked, until, linkErr = cli.LookupAgent(st.Network, st.Wallet, st.WorkspaceID, sf.AgentAddr, now)
	}
	if prepare {
		p, hash, err := cli.BindConnectionPreview(stateDir(), "ETH")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(cli.ConnectionPreviewPublic(p, hash))
			return
		}
		fmt.Println("kind      connection_test")
		fmt.Printf("market    %s\n", p.Market)
		fmt.Printf("side      %s\n", p.Side)
		fmt.Printf("sz        %v\n", p.Sz)
		fmt.Printf("limit     %s\n", p.LimitPx)
		fmt.Printf("hash      %s\n", hash)
		fmt.Println("this is not a research recommendation. type AUTHORIZE on the desktop to send it.")
		return
	}
	if asJSON {
		body := map[string]any{
			"ok": true, "sign": false, "trade": false, "withdraw": false,
			"order": true, "cancel": true, "api": "https://app.hyperliquid.xyz/API",
		}
		if serr == nil {
			body["agent"] = sf.AgentAddr
			body["expires"] = sf.Expires
			body["sessionAlive"] = session.Alive(sf.Meta().Session(), now)
		}
		if name != "" {
			body["name"] = name
		}
		body["linked"] = linked
		if until > 0 {
			body["approvedUntil"] = until
		}
		if linkErr != nil {
			body["error"] = "agent_list_query_failed"
		}
		_ = json.NewEncoder(os.Stdout).Encode(body)
		return
	}
	fmt.Println("Open https://app.hyperliquid.xyz/API")
	fmt.Println("API wallet name must be less than 17 characters.")
	if name != "" {
		fmt.Println("name", name)
	}
	if serr == nil {
		fmt.Println("address", sf.AgentAddr)
		fmt.Println("expires", sf.Expires)
	}
	fmt.Println(cli.LinkCopy(linked, linkErr))
	if linked {
		fmt.Println("next      session is approved. AUTHORIZE stays on this computer.")
	} else if serr == nil {
		fmt.Println("next      paste the address into Hyperliquid API and click Authorize API Wallet.")
	} else {
		fmt.Println("next      create a local session on PIT Desktop, then approve that agent.")
	}
	fmt.Println("Click Authorize API Wallet. PIT still cannot withdraw.")
}

func cmdCompute() {
	checks := cli.Doctor(stateDir())
	if asJSON {
		cli.PrintDoctor(os.Stdout, checks, true)
		return
	}
	for _, c := range checks {
		if c.Name == "direct_auth" || c.Name == "direct_sealer" || c.Name == "direct_credit" || c.Name == "tee" {
			mark := "fail"
			if c.OK {
				mark = "ok"
			}
			fmt.Printf("%-16s %s  %s\n", c.Name, mark, c.Detail)
		}
	}
	fmt.Println("Open https://pc.0g.ai/sdk/dashboard/funds")
	fmt.Println("That page is provider credit, not a Hyperliquid balance.")
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
	auth, _, err := cli.ResolveWorkspaceAuth(stateDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := compute.ProductAskAuth(net, true, compute.LookBin(), market, book, auth); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("sealed ask submitted")
}

func cmdDirect(args []string) {
	dir := stateDir()
	var sig string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--sig", "--signature":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "direct --sig requires a hex signature")
				os.Exit(2)
			}
			sig = args[i+1]
			i++
		}
	}
	if sig == "" {
		ch, err := cli.IssueDirectChallenge(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
				"message":   ch.Message,
				"digest":    ch.Digest,
				"provider":  ch.Provider,
				"model":     ch.Model,
				"explain":   ch.Explain,
				"expiresAt": ch.ExpiresAt,
				"sign":      false,
				"trade":     false,
			})
			return
		}
		fmt.Println(ch.Explain)
		fmt.Println("Digest", ch.Digest)
		fmt.Println("Sign this digest with the bound wallet (raw 32-byte personal_sign). Then run pit direct --sig 0x...")
		fmt.Println("PIT never asks for a seed or a private key.")
		return
	}
	meta, err := cli.CompleteDirect(dir, "", sig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"ok":        true,
			"provider":  meta.Provider,
			"model":     meta.Model,
			"expiresAt": meta.ExpiresAt,
			"source":    meta.Source,
			"sign":      false,
			"trade":     false,
		})
		return
	}
	fmt.Println("Direct token stored in the OS keychain. The token is not printed.")
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

func cmdUpdate() {
	body := map[string]any{
		"ok":       true,
		"sign":     false,
		"trade":    false,
		"url":      "https://github.com/mohamedwael201193/pit/releases/latest",
		"unsigned": true,
		"note":     "Windows installers are unsigned until Authenticode is obtained. PIT does not claim a signature it does not have.",
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(body)
		return
	}
	fmt.Println(body["url"])
	fmt.Println(body["note"])
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
	_ = os.Setenv("PIT_COMPANION", "1")
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
