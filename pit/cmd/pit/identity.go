package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mohamedwael201193/pit/internal/chain8004"
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskid"
	"github.com/mohamedwael201193/pit/internal/evidence"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func cmdIdentity(args []string) {
	sub := "verify"
	if len(args) > 0 {
		sub = args[0]
	}
	switch sub {
	case "verify", "":
		identityVerify()
	case "calldata":
		identityCalldata(args[1:])
	case "apply":
		identityApply(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "usage: pit identity [verify|calldata|apply] ...")
		os.Exit(2)
	}
}

func identityChain() (config.Chain, cli.DiskState) {
	st, err := cli.Load(stateDir())
	net := config.Mainnet
	if err == nil {
		if n, e := config.ParseNetwork(st.Network); e == nil {
			net = n
		}
	} else {
		st = cli.DiskState{Wallet: chain8004.RecordedOwner, Network: string(config.Mainnet)}
	}
	return config.For(net), st
}

func identityVerify() {
	ch, _ := identityChain()
	desk, err := deskid.Snapshot(ch, deskid.MainnetTokenID, deskid.RecordedAgent)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	id, err := chain8004.IdentitySnapshot(ch, chain8004.MainnetAgentID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	reg, _ := chain8004.IdentityRegistry(ch)
	impl, _ := chain8004.Implementation(ch.RPC, ch.Reputation8004)
	body := map[string]any{
		"desk":               ch.DeskID,
		"desk_token":         deskid.MainnetTokenID,
		"desk_owner":         desk.Owner,
		"desk_uri":           desk.URI,
		"agent_authorized":   desk.UserAuthorized,
		"ids":                []string{deskid.ID7857, deskid.ID7857Authorize, deskid.ID7857Cloneable, deskid.ID721},
		"supports_7857":      desk.Supports7857,
		"supports_authorize": desk.SupportsAuthorize,
		"identity":           ch.Identity8004,
		"agent_id":           chain8004.MainnetAgentID,
		"identity_owner":     id.Owner,
		"agent_uri":          id.URI,
		"reputation":         ch.Reputation8004,
		"identity_registry":  reg,
		"reputation_impl":    impl,
		"sign":               false,
		"trade":              false,
	}
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(body)
		return
	}
	fmt.Printf("desk     %s token %d owner %s\n", ch.DeskID, deskid.MainnetTokenID, desk.Owner)
	fmt.Printf("agent    authorized=%v\n", desk.UserAuthorized)
	fmt.Printf("8004     agent %d owner %s\n", chain8004.MainnetAgentID, id.Owner)
	fmt.Printf("uri      %s\n", id.URI)
	fmt.Printf("rep impl %s\n", impl)
}

func identityCalldata(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: pit identity calldata authorize|revoke|set-uri ...")
		os.Exit(2)
	}
	switch args[0] {
	case "authorize":
		user := deskid.RecordedAgent
		if len(args) > 1 {
			user = args[1]
		}
		data, err := deskid.EncodeAuthorizeUsage(deskid.MainnetTokenID, user)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(deskid.CalldataHex(data))
	case "revoke":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: pit identity calldata revoke 0x...")
			os.Exit(2)
		}
		data, err := deskid.EncodeRevokeAuthorization(deskid.MainnetTokenID, args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(deskid.CalldataHex(data))
	case "set-uri":
		uri := chain8004.AgentCardURL
		if len(args) > 1 {
			uri = args[1]
		}
		data, err := chain8004.EncodeSetAgentURI(chain8004.MainnetAgentID, uri)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(chain8004.CalldataHex(data))
	default:
		fmt.Fprintln(os.Stderr, "unknown calldata")
		os.Exit(2)
	}
}

func identityApply(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: pit identity apply authorize-agent|set-uri|feedback --confirm")
		os.Exit(2)
	}
	confirm := false
	kind := args[0]
	for _, a := range args[1:] {
		if a == "--confirm" {
			confirm = true
		}
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "refusing unsigned apply: pass --confirm")
		os.Exit(2)
	}
	ch, st := identityChain()
	switch kind {
	case "authorize-agent":
		key, err := evidence.PayerKey(stateDir(), st.Network, st.WorkspaceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		hash, err := deskid.AuthorizeUsageFromOwner(ch, deskid.MainnetTokenID, key, deskid.RecordedAgent)
		if err != nil {
			if strings.Contains(err.Error(), "already_authorized") {
				fmt.Println("already_authorized")
				return
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(hash)
	case "set-uri":
		key, err := evidence.PayerKey(stateDir(), st.Network, st.WorkspaceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		hash, err := chain8004.SetURIFromOwner(ch, key, chain8004.MainnetAgentID, chain8004.AgentCardURL)
		if err != nil {
			if strings.Contains(err.Error(), "uri_unchanged") {
				fmt.Println("uri_unchanged")
				return
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(hash)
	case "feedback":
		key := chain8004.ReporterKey()
		if key == "" {
			fmt.Fprintln(os.Stderr, "set PIT_8004_REPORTER_KEY or PIT_DEPLOYER_KEY")
			os.Exit(1)
		}
		c := hl.New(ch)
		wallet := st.Wallet
		if wallet == "" {
			wallet = chain8004.RecordedOwner
		}
		fills, err := c.UserFills(wallet)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		hash, err := chain8004.SubmitFillFeedback(ch, key, chain8004.MainnetAgentID, chain8004.RecordedOID, fills)
		if err != nil {
			if strings.Contains(err.Error(), "already_reported") {
				fmt.Println("already_reported")
				return
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(hash)
	default:
		fmt.Fprintln(os.Stderr, "unknown apply")
		os.Exit(2)
	}
}
