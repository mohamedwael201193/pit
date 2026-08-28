package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/keyring"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func challengePath(dir, workspace string) string {
	return filepath.Join(dir, workspace+".direct-challenge.json")
}

func saveChallenge(dir, workspace string, ch compute.Challenge) error {
	b, err := json.MarshalIndent(ch, "", "  ")
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(string(b)), "app-sk-") {
		return fmt.Errorf("companion_leak")
	}
	return os.WriteFile(challengePath(dir, workspace), b, 0o600)
}

func loadChallenge(dir, workspace string) (compute.Challenge, error) {
	b, err := os.ReadFile(challengePath(dir, workspace))
	if err != nil {
		return compute.Challenge{}, fmt.Errorf("direct_challenge_required")
	}
	var ch compute.Challenge
	if err := json.Unmarshal(b, &ch); err != nil {
		return compute.Challenge{}, fmt.Errorf("direct_challenge_required")
	}
	return ch, nil
}

func IssueDirectChallenge(dir string) (compute.Challenge, error) {
	st, err := Load(dir)
	if err != nil {
		return compute.Challenge{}, fmt.Errorf("unbound")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return compute.Challenge{}, err
	}
	if net == config.Testnet {
		return compute.Challenge{}, fmt.Errorf("galileo_e2ee_unproven")
	}
	sku := compute.ForNetwork(net)
	if err := compute.SealedAskEnabled(sku); err != nil {
		return compute.Challenge{}, err
	}
	probe := compute.ProbeDirectAccount(config.For(net), st.Wallet, sku.Provider)
	_, ch, err := compute.NewChallenge(st.Wallet, sku.Provider, probe.Generation, time.Now())
	if err != nil {
		return compute.Challenge{}, err
	}
	ch.Model = sku.Model
	ch.Network = string(net)
	if !probe.Present || !probe.Acknowledged {
		ch.Explain = ch.Explain + " If the provider later rejects the token, fund Direct at pc.0g.ai Advanced with this same wallet. PIT does not ask for a private key."
	}
	if err := saveChallenge(dir, st.WorkspaceID, ch); err != nil {
		return compute.Challenge{}, err
	}
	return ch, nil
}

func CompleteDirect(dir, message, signature string) (compute.DirectMeta, error) {
	st, err := Load(dir)
	if err != nil {
		return compute.DirectMeta{}, fmt.Errorf("unbound")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return compute.DirectMeta{}, err
	}
	sku := compute.ForNetwork(net)
	if strings.TrimSpace(message) == "" {
		ch, err := loadChallenge(dir, st.WorkspaceID)
		if err != nil {
			return compute.DirectMeta{}, err
		}
		message = ch.Message
	}
	auth, tok, err := compute.AcceptDirectSignature(message, signature, st.Wallet, time.Now())
	if err != nil {
		return compute.DirectMeta{}, err
	}
	want, err := compute.ChecksumAddress(sku.Provider)
	if err != nil {
		return compute.DirectMeta{}, err
	}
	if !strings.EqualFold(tok.Provider, want) {
		return compute.DirectMeta{}, fmt.Errorf("provider_spoof")
	}
	store, err := keyring.OpenProduct(KeyringDir(dir))
	if err != nil {
		return compute.DirectMeta{}, err
	}
	if err := compute.StoreDirect(store, string(net), st.WorkspaceID, sku.Provider, auth); err != nil {
		return compute.DirectMeta{}, err
	}
	_ = os.Remove(challengePath(dir, st.WorkspaceID))
	return compute.PublicMeta(sku, tok, "keychain"), nil
}

func ResolveWorkspaceAuth(dir string) (compute.AuthFile, compute.DirectMeta, error) {
	st, err := Load(dir)
	if err != nil {
		return compute.AuthFile{}, compute.DirectMeta{}, fmt.Errorf("unbound")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return compute.AuthFile{}, compute.DirectMeta{}, err
	}
	store, err := keyring.OpenProduct(KeyringDir(dir))
	if err != nil {
		file, err := compute.LoadEnvAuthFile()
		if err != nil {
			return compute.AuthFile{}, compute.DirectMeta{}, fmt.Errorf("direct_token_required")
		}
		tok, _, perr := compute.ParseBearer(file.Authorization)
		meta := compute.DirectMeta{Provider: file.Provider, Model: file.Model, Source: "operator_file"}
		if perr == nil {
			meta = compute.PublicMeta(compute.ForNetwork(net), tok, "operator_file")
		}
		return file, meta, nil
	}
	var recovery keyring.Store
	if keyring.BackendName() != "file" {
		recovery, _ = keyring.Open(KeyringDir(dir))
	}
	file, meta, err := compute.AuthFromStores(store, recovery, net, st.WorkspaceID, time.Now())
	if err == nil {
		return file, meta, nil
	}
	if err.Error() == "direct_token_expired" {
		return compute.AuthFile{}, compute.DirectMeta{}, err
	}
	op, err := compute.LoadEnvAuthFile()
	if err != nil {
		return compute.AuthFile{}, compute.DirectMeta{}, fmt.Errorf("direct_token_required")
	}
	tok, _, perr := compute.ParseBearer(op.Authorization)
	meta = compute.DirectMeta{Provider: op.Provider, Model: op.Model, Source: "operator_file"}
	if perr == nil {
		meta = compute.PublicMeta(compute.ForNetwork(net), tok, "operator_file")
	}
	return op, meta, nil
}

func DirectStatus(dir string) map[string]any {
	body := map[string]any{"sign": false, "trade": false, "ok": false, "source": ""}
	_, meta, err := ResolveWorkspaceAuth(dir)
	if err != nil {
		body["error"] = err.Error()
		return body
	}
	body["ok"] = true
	body["source"] = meta.Source
	body["provider"] = meta.Provider
	body["model"] = meta.Model
	body["expiresAt"] = meta.ExpiresAt
	body["tokenId"] = meta.TokenId
	body["wallet"] = meta.Wallet
	return body
}

func ForgetDirect(dir string) {
	st, err := Load(dir)
	if err != nil {
		return
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return
	}
	sku := compute.ForNetwork(net)
	if ring, err := keyring.OpenProduct(KeyringDir(dir)); err == nil {
		_ = compute.DeleteDirect(ring, string(net), st.WorkspaceID, sku.Provider)
	}
	_ = os.Remove(challengePath(dir, st.WorkspaceID))
}

func RunWorkspaceResearch(dir, coin string) (compute.AskReport, error) {
	return RunWorkspaceResearchStage(dir, coin, nil, nil)
}

func RunWorkspaceResearchStage(dir, coin string, stage compute.StageFn, stop func() bool) (compute.AskReport, error) {
	if stage != nil {
		stage("READING_MARKET")
	}
	st, err := Load(dir)
	if err != nil {
		return compute.AskReport{}, fmt.Errorf("unbound")
	}
	net, err := config.ParseNetwork(st.Network)
	if err != nil {
		return compute.AskReport{}, err
	}
	auth, _, err := ResolveWorkspaceAuth(dir)
	if err != nil {
		return compute.AskReport{}, err
	}
	p := policy.Default()
	_ = CheckPinned(dir, st.WorkspaceID, p)
	if p.KillSwitch {
		return compute.AskReport{}, fmt.Errorf("kill_switch")
	}
	hash, err := p.Hash()
	if err != nil {
		return compute.AskReport{}, err
	}
	book, err := compute.BuildPrivateBook(st.Wallet, st.WorkspaceID, st.Network, hash)
	if err != nil {
		return compute.AskReport{}, err
	}
	want := strings.ToUpper(strings.TrimSpace(coin))
	if want == "" {
		want = "ETH"
	}
	allowed := false
	for _, a := range p.AllowedAssets {
		if strings.EqualFold(a, want) {
			allowed = true
			break
		}
	}
	if !allowed {
		return compute.AskReport{}, fmt.Errorf("asset_not_allowed")
	}
	snap, err := hl.New(config.For(net)).PublicBook(want)
	if err != nil || snap.MarkPx <= 0 {
		return compute.AskReport{}, fmt.Errorf("empty_envelope")
	}
	market, err := compute.MarketJSON(snap)
	if err != nil {
		return compute.AskReport{}, err
	}
	last := filepath.Join(dir, "last-research.json")
	rep, err := compute.ProductAskReportStage(net, true, compute.LookBin(), market, book, auth, last, stage, stop)
	if err != nil {
		return rep, err
	}
	if stage != nil {
		stage("POLICY")
	}
	return rep, nil
}
