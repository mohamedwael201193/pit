package compute

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestRefuseRouterKey(t *testing.T) {
	if err := RefuseRouterKey("Bearer sk-dashboard"); err == nil {
		t.Fatal("sk")
	}
	if err := RefuseRouterKey("Bearer mk-dashboard"); err == nil {
		t.Fatal("mk")
	}
	if err := RefuseRouterKey("Bearer app-sk-direct"); err != nil {
		t.Fatal(err)
	}
}

func TestWriteAuthRefusesRouterURL(t *testing.T) {
	sku := MainnetChat()
	sku.URL = "https://router-api.0g.ai/v1"
	if err := WriteAuth(filepath.Join(t.TempDir(), "a.json"), sku, "Bearer app-sk-x"); err == nil {
		t.Fatal("router")
	}
}

func TestWriteAuthWritesTeeML(t *testing.T) {
	p := filepath.Join(t.TempDir(), "auth.json")
	if err := WriteAuth(p, MainnetChat(), "Bearer app-sk-x"); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty")
	}
}

func TestProductAskRequiresDirectAuthFile(t *testing.T) {
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	err := ProductAsk(config.Mainnet, true, "/opt/pit/sealer")
	if err == nil || err.Error() != "direct_token_required" {
		t.Fatalf("%v", err)
	}
}
