package compute

import (
	"path/filepath"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestProductAskMissingBinaryAfterAuth(t *testing.T) {
	p := filepath.Join(t.TempDir(), "auth.json")
	if err := WriteAuth(p, MainnetChat(), "Bearer app-sk-x"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PIT_DIRECT_AUTH_FILE", p)
	err := ProductAskEnvelope(config.Mainnet, true, filepath.Join(t.TempDir(), "no-such-sealer"), []byte("mkt"), []byte("book"))
	if err == nil || err.Error() != "sealer_not_wired" {
		t.Fatalf("%v", err)
	}
}

func TestProductAskEmptyEnvelope(t *testing.T) {
	p := filepath.Join(t.TempDir(), "auth.json")
	if err := WriteAuth(p, MainnetChat(), "Bearer app-sk-x"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PIT_DIRECT_AUTH_FILE", p)
	err := ProductAsk(config.Mainnet, true, "/opt/pit/sealer")
	if err == nil || err.Error() != "empty_envelope" {
		t.Fatalf("%v", err)
	}
}
