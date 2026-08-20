package compute

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestSKUIsolation(t *testing.T) {
	if err := RefuseSKUCopy(config.Testnet, "glm-5.2"); err == nil {
		t.Fatal("copied glm onto galileo")
	}
	if err := RefuseSKUCopy(config.Mainnet, "qwen/qwen2.5-omni-7b"); err == nil {
		t.Fatal("copied omni onto aristotle")
	}
	if err := SealedAskEnabled(TestnetChat()); err == nil {
		t.Fatal("omni not proven")
	}
	if err := SealedAskEnabled(MainnetChat()); err != nil {
		t.Fatal(err)
	}
	if ForNetwork(config.Testnet).Model == ForNetwork(config.Mainnet).Model {
		t.Fatal("same catalog")
	}
}
