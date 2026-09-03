package compute

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestSKUIsolation(t *testing.T) {
	if err := RefuseSKUCopy(config.Testnet, "glm-5.2"); err == nil {
		t.Fatal("copied glm onto galileo")
	}
	if err := RefuseSKUCopy(config.Testnet, "glm-5.3"); err == nil {
		t.Fatal("copied glm-5.3 onto galileo")
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
	m := MainnetChat()
	if m.Model != "glm-5.3" || m.Verifiability != "TeeML" || !m.ProvenE2EE {
		t.Fatalf("%+v", m)
	}
	if m.Provider != "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D" {
		t.Fatal("do not swap provider")
	}
	if m.URL != "https://compute-network-19.integratenetwork.work" {
		t.Fatal("do not swap url")
	}
	if m.TeeSigner != "0x089EBc23206267FCD5ef46725c6196DF21bE45D7" {
		t.Fatal("do not swap teeSigner")
	}
}
