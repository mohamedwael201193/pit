package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	out := []Check{
		checkVersion(),
		checkWallet(dir),
		checkNetwork(dir),
		checkKeychain(dir),
		checkMemoryEnv(),
		checkHyperliquid(dir),
		checkRPC(dir),
		checkCompanion(),
		checkSealer(),
		checkDirectAuth(),
		checkStorage(),
		checkRegistry(dir),
		checkSession(dir),
		checkPolicy(dir),
	}
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
	if _, err := hl.New(config.For(net)).PublicBook("ETH"); err != nil {
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
	client := &http.Client{Timeout: 8 * time.Second}
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

func checkDirectAuth() Check {
	p := strings.TrimSpace(os.Getenv("PIT_DIRECT_AUTH_FILE"))
	if p == "" {
		return Check{Name: "direct_auth", Detail: "PIT_DIRECT_AUTH_FILE unset"}
	}
	if _, err := os.Stat(p); err != nil {
		return Check{Name: "direct_auth", Detail: "auth file missing"}
	}
	return Check{Name: "direct_auth", OK: true, Detail: "file present"}
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
