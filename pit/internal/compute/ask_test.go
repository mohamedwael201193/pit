package compute

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestSkuURLMatchTrimsSlash(t *testing.T) {
	if !skuURLMatch("https://compute-network-19.integratenetwork.work", "https://compute-network-19.integratenetwork.work/") {
		t.Fatal("slash")
	}
	if skuURLMatch("https://compute-network-19.integratenetwork.work", "https://router-api.0g.ai") {
		t.Fatal("router")
	}
}

func TestProductAskRequiresDeskAndSealer(t *testing.T) {
	if err := ProductAsk(config.Mainnet, false, "/opt/pit/sealer"); err == nil {
		t.Fatal("desk")
	}
	if err := ProductAsk(config.Testnet, true, "/opt/pit/sealer"); err == nil {
		t.Fatal("galileo")
	}
	err := ProductAsk(config.Mainnet, true, "")
	if err == nil || err.Error() != "sealer_not_wired" {
		t.Fatalf("%v", err)
	}
}

func TestSavePublicEvidenceRedactsDirectToken(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "researcher.json")
	body := `{"verify_e2ee":"FAIL","post_err_clip":"Bearer app-sk-secret","verify_err":"ok","sanitized_output":"nope"}`
	if err := os.WriteFile(out, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "last-research.json")
	if err := SavePublicEvidence(dest, []DirectJob{{OutPath: out, Role: Researcher}}, fmt.Errorf("TEE_VERIFY_FAIL")); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		t.Fatal(string(raw))
	}
	if strings.Contains(string(raw), "sanitized_output") {
		t.Fatal(string(raw))
	}
	if !strings.Contains(string(raw), "TEE_VERIFY_FAIL") {
		t.Fatal(string(raw))
	}
	if !strings.Contains(string(raw), `"verify_e2ee"`) {
		t.Fatal("verify field")
	}
}

func TestSavePublicEvidenceKeepsProposedSide(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "researcher.json")
	body := `{"verify_e2ee":"OK","sanitized_output":"{\"proposed_side\":\"buy\"}"}`
	if err := os.WriteFile(out, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "last-research.json")
	if err := SavePublicEvidence(dest, []DirectJob{{OutPath: out, Role: Researcher}}, nil); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "sanitized_output") {
		t.Fatal(string(raw))
	}
	if !strings.Contains(string(raw), `"proposed_side"`) || !strings.Contains(string(raw), `"buy"`) {
		t.Fatal(string(raw))
	}
}
