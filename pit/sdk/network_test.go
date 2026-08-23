package sdk

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestBindNetwork(t *testing.T) {
	c := Client{Network: config.Mainnet}
	if err := c.BindNetwork(config.Mainnet); err != nil {
		t.Fatal(err)
	}
	if err := c.BindNetwork(config.Testnet); err == nil {
		t.Fatal("mix")
	}
}
