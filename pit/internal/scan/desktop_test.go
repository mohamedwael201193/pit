package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDesktopSourceHasNoSessionMaterial(t *testing.T) {
	root := filepath.Join("..", "..", "..", "apps", "desktop", "src")
	if err := DesktopSource(root); err != nil {
		t.Fatal(err)
	}
}

func TestOfficialLinkCatalog(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "apps", "desktop", "src", "links.ts"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	for _, want := range []string{
		"https://pit0g.vercel.app/pair",
		"https://pit0g.vercel.app/app",
		"https://pit0g.vercel.app/protect",
		"https://pc.0g.ai/sdk/dashboard/funds",
		"https://app.hyperliquid.xyz",
		"https://app.hyperliquid.xyz/API",
		"https://app.hyperliquid-testnet.xyz",
		"https://app.hyperliquid-testnet.xyz/API",
		"https://github.com/mohamedwael201193/pit/releases/latest",
		"https://chainscan.0g.ai",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %s", want)
		}
	}
	if strings.Contains(body, "extraAgents") {
		t.Fatal("protocol jargon in link catalog")
	}
}
