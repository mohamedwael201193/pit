package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/keyring"
	"github.com/mohamedwael201193/pit/internal/storage"
	"github.com/mohamedwael201193/pit/internal/version"
)

type Check struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

func Doctor(dir string) []Check {
	var hlC, rpcC, authC, creditC, agentC Check
	var wg sync.WaitGroup
	wg.Add(5)
	go func() { defer wg.Done(); hlC = checkHyperliquid(dir) }()
	go func() { defer wg.Done(); rpcC = checkRPC(dir) }()
	go func() { defer wg.Done(); authC = checkDirectAuth(dir) }()
	go func() { defer wg.Done(); creditC = checkDirectCredit(dir) }()
	go func() { defer wg.Done(); agentC = checkHLAgent(dir) }()
	out := []Check{
		checkVersion(),
		checkWallet(dir),
		checkNetwork(dir),
		checkKeychain(dir),
		checkMemoryEnv(),
	}
	wg.Wait()
	out = append(out,
		hlC,
		rpcC,
		checkCompanion(),
		checkSealer(),
		authC,
		creditC,
		checkTee(dir),
		checkStorage(),
		checkRegistry(dir),
		checkSession(dir),
		agentC,
		checkPolicy(dir),
	)
	return out
}

func checkVersion() Check {
	return Check{Name: "version", OK: true, Detail: version.String()}
}

func checkMemoryEnv() Check {
	if strings.TrimSpace(os.Getenv("PIT_MEMORY_KEY")) != "" {
		return Check{Name: "memory_key", Detail: "global PIT_MEMORY_KEY is set. Product refuses it. Use a generated --key-file per workspace."}
	}
	return Check{Name: "memory_key", OK: true, Detail: "no global memory key"}
}

func checkWallet(dir string) Check {
	st, err := Load(dir)
	if err != nil {
		return Check{Name: "wallet", Detail: "unbound until pit init"}
	}
	return Check{Name: "wallet", OK: true, Detail: st.Wallet}
}

func checkNetwork(dir string) Check {
	st, err := Load(dir)
	if err != nil {
		return Check{Name: "network", Detail: "unbound"}
	}
	if _, err := config.ParseNetwork(st.Network); err != nil {
		return Check{Name: "network", Detail: err.Error()}
	}
	return Check{Name: "network", OK: true, Detail: st.Network}
}

func checkKeychain(dir string) Check {
	store, err := keyring.OpenProduct(KeyringDir(dir))
	if err != nil {
		return Check{Name: "keychain", Detail: err.Error()}
	}
	if err := store.Put("doctor", "ping", []byte("ok")); err != nil {
		return Check{Name: "keychain", Detail: err.Error()}
	}
	_ = store.Delete("doctor", "ping")
	return Check{Name: "keychain", OK: true, Detail: keyring.BackendName()}
}

func checkHyperliquid(dir string) Check {
	net := config.Mainnet
	if st, err := Load(dir); err == nil {
		if n, err := config.ParseNetwork(st.Network); err == nil {
			net = n
		}
	}
	cl := hl.New(config.For(net))
	cl.HTTP = &http.Client{Timeout: 5 * time.Second}
	if _, err := cl.PublicBook("ETH"); err != nil {
		return Check{Name: "hyperliquid", Detail: err.Error()}
	}
	return Check{Name: "hyperliquid", OK: true, Detail: string(net) + " public book"}
}

func checkRPC(dir string) Check {
	net := config.Mainnet
	if st, err := Load(dir); err == nil {
		if n, err := config.ParseNetwork(st.Network); err == nil {
			net = n
		}
	}
	ch := config.For(net)
	if err := config.RefuseMixedRPC(ch.ChainID, ch.RPC); err != nil {
		return Check{Name: "0g_rpc", Detail: err.Error()}
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(ch.RPC, "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`))
	if err != nil {
		return Check{Name: "0g_rpc", Detail: err.Error()}
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if !strings.Contains(strings.ToLower(string(b)), "result") {
		return Check{Name: "0g_rpc", Detail: "rpc_unreadable"}
	}
	return Check{Name: "0g_rpc", OK: true, Detail: fmt.Sprintf("%s %d", net, ch.ChainID)}
}

func checkCompanion() Check {
	if os.Getenv("PIT_COMPANION") == "1" {
		return Check{Name: "desktop", OK: true, Detail: "this process"}
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:17373/health")
	if err != nil {
		return Check{Name: "desktop", Detail: "companion not listening on 127.0.0.1:17373"}
	}
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["sign"] == true {
		return Check{Name: "desktop", Detail: "companion_must_not_sign"}
	}
	return Check{Name: "desktop", OK: true, Detail: "loopback companion"}
}

func checkSealer() Check {
	bin := compute.LookBin()
	if bin == "" {
		return Check{Name: "direct_sealer", Detail: "PIT_COMMITTEE_BIN empty and sealer/pit-sealer missing"}
	}
	return Check{Name: "direct_sealer", OK: true, Detail: "binary present"}
}

func checkDirectAuth(dir string) Check {
	_, meta, err := ResolveWorkspaceAuth(dir)
	if err == nil {
		src := meta.Source
		if src == "" {
			src = "keychain"
		}
		if src == "sponsor" {
			return Check{Name: "direct_auth", OK: true, Detail: "sponsored Direct on this computer. Trading credentials stay separate."}
		}
		return Check{Name: "direct_auth", OK: true, Detail: "wallet-signed Direct token in " + src}
	}
	switch err.Error() {
	case "unbound":
		return Check{Name: "direct_auth", Detail: "unbound until wallet is bound"}
	case "direct_token_expired":
		return Check{Name: "direct_auth", Detail: "Direct token expired. Sign again in the paired browser."}
	case "galileo_e2ee_unproven":
		return Check{Name: "direct_auth", Detail: "Galileo TeeML is unproven. Switch to MAINNET for sealed research."}
	default:
		return Check{Name: "direct_auth", Detail: "Direct token missing. Pair the browser and sign the sealed-path message."}
	}
}

func checkDirectCredit(dir string) Check {
	st, err := Load(dir)
	if err != nil {
		return Check{Name: "direct_credit", Detail: "unbound"}
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return Check{Name: "direct_credit", Detail: err.Error()}
	}
	sku := compute.ForNetwork(net)
	probe := compute.ProbeDirectAccount(config.For(net), st.Wallet, sku.Provider)
	teeOK := checkTee(dir).OK
	if !probe.Present {
		if teeOK {
			return Check{Name: "direct_credit", OK: true, Detail: "last committee verified. Provider ledger unread this pass — not treated as 0 0G."}
		}
		detail := "could not read the provider ledger. Open pc.0g.ai Advanced. This is unread, not zero."
		if probe.Err != "" {
			detail = "could not read the provider ledger (" + probe.Err + "). Open pc.0g.ai Advanced. This is unread, not zero."
		}
		return Check{Name: "direct_credit", Detail: detail}
	}
	if probe.EnoughForCommittee() {
		return Check{Name: "direct_credit", OK: true, Detail: "provider credit " + probe.BalanceOG() + " 0G"}
	}
	if _, err := compute.LoadSponsorAuthFile(); err == nil {
		return Check{Name: "direct_credit", OK: true, Detail: "your provider credit is " + probe.BalanceOG() + " 0G. PIT can sponsor sealed research within a daily workspace cap. Open pc.0g.ai Advanced to fund your own sub-account."}
	}
	return Check{Name: "direct_credit", Detail: "provider credit " + probe.BalanceOG() + " 0G. Three sealed roles need about 3 0G locked. Open pc.0g.ai Advanced with this wallet. PIT does not ask for a private key."}
}

func checkTee(dir string) Check {
	raw, err := os.ReadFile(filepath.Join(dir, "last-research.json"))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return Check{Name: "tee", Detail: "No sealed research has been verified on this machine yet."}
	}
	var body struct {
		Roles []struct {
			Verify string `json:"verify_e2ee"`
			Role   string `json:"role"`
		} `json:"roles"`
	}
	if json.Unmarshal(raw, &body) != nil {
		return Check{Name: "tee", Detail: "No sealed research has been verified on this machine yet."}
	}
	ok := 0
	for _, r := range body.Roles {
		if strings.EqualFold(strings.TrimSpace(r.Verify), "OK") {
			ok++
		}
	}
	if ok == 0 {
		return Check{Name: "tee", Detail: "Last sealed run did not verify a TeeML signature."}
	}
	return Check{Name: "tee", OK: true, Detail: fmt.Sprintf("%d role(s) VerifyE2EE OK on this machine", ok)}
}

func checkStorage() Check {
	if err := storage.RefuseMissingProof(storage.LookCLI()); err != nil {
		return Check{Name: "storage", Detail: err.Error()}
	}
	return Check{Name: "storage", OK: true, Detail: "official Go client"}
}

func checkRegistry(dir string) Check {
	net := config.Mainnet
	if st, err := Load(dir); err == nil {
		if n, err := config.ParseNetwork(st.Network); err == nil {
			net = n
		}
	}
	ch := config.For(net)
	if ch.Identity8004 == "" || ch.Reputation8004 == "" {
		return Check{Name: "registry", Detail: "8004 addresses missing"}
	}
	return Check{Name: "registry", OK: true, Detail: "erc-8004 addresses pinned"}
}

func checkSession(dir string) Check {
	st, err := Load(dir)
	if err != nil {
		return Check{Name: "session", Detail: "none"}
	}
	if _, err := LiveFromDisk(dir, st.Kill, time.Now().UnixMilli()); err != nil {
		return Check{Name: "session", Detail: err.Error()}
	}
	return Check{Name: "session", OK: true, Detail: "live order/cancel only"}
}

func checkHLAgent(dir string) Check {
	st, err := Load(dir)
	if err != nil {
		return Check{Name: "hl_agent", Detail: "unbound"}
	}
	live, err := LiveFromDisk(dir, st.Kill, time.Now().UnixMilli())
	if err != nil {
		return Check{Name: "hl_agent", Detail: "Create a local session, then approve it on Hyperliquid. PIT cannot withdraw."}
	}
	linked, err := LiveLinked(st.Network, st.Wallet, live.Workspace, live.AgentAddr, time.Now().UnixMilli())
	if err != nil {
		return Check{Name: "hl_agent", Detail: "extraAgents query failed. AUTHORIZE stays off until Hyperliquid lists this agent."}
	}
	if !linked {
		return Check{Name: "hl_agent", Detail: "Approve this agent on Hyperliquid. extraAgents must list it. PIT never signs withdraw, transfer, leverage, or account admin."}
	}
	return Check{Name: "hl_agent", OK: true, Detail: "extraAgents lists this session. PIT still refuses withdraw, transfer, leverage, and account admin."}
}

func checkPolicy(dir string) Check {
	st, err := Load(dir)
	if err != nil {
		return Check{Name: "policy", Detail: "unbound"}
	}
	if _, err := os.Stat(filepath.Join(dir, st.WorkspaceID+".policy")); err != nil {
		return Check{Name: "policy", Detail: "not pinned"}
	}
	return Check{Name: "policy", OK: true, Detail: "pinned"}
}

func DoctorFailed(checks []Check) bool {
	required := map[string]bool{"wallet": true, "network": true, "keychain": true, "hyperliquid": true, "0g_rpc": true, "memory_key": true}
	for _, c := range checks {
		if required[c.Name] && !c.OK {
			return true
		}
	}
	return false
}

func PrintDoctor(w io.Writer, checks []Check, asJSON bool) {
	if asJSON {
		_ = json.NewEncoder(w).Encode(map[string]any{"checks": checks, "sign": false, "version": version.Number})
		return
	}
	for _, c := range checks {
		mark := "fail"
		if c.OK {
			mark = "ok"
		}
		fmt.Fprintf(w, "%-16s %s  %s\n", c.Name, mark, c.Detail)
	}
}
