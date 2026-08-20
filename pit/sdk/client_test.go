package sdk

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestSDKNeverSigns(t *testing.T) {
	c := Client{Network: config.Mainnet}
	if c.Status().CanSign {
		t.Fatal("web/sdk cannot sign")
	}
	if !strings.Contains(c.Explorer("0xabc"), "chainscan.0g.ai") {
		t.Fatal("explorer")
	}
	c.Network = config.Testnet
	if strings.Contains(c.Explorer(""), "chainscan.0g.ai") && !strings.Contains(c.Explorer(""), "galileo") {
		t.Fatal("testnet explorer")
	}
}
