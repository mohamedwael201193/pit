package compute

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
)

type SKU struct {
	Network       config.Network
	Model         string
	Provider      string
	TeeSigner     string
	URL           string
	Verifiability string
	ProvenE2EE    bool
}

func MainnetChat() SKU {
	return SKU{
		Network:       config.Mainnet,
		Model:         "glm-5.3",
		Provider:      "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D",
		TeeSigner:     "0x089EBc23206267FCD5ef46725c6196DF21bE45D7",
		URL:           "https://compute-network-19.integratenetwork.work",
		Verifiability: "TeeML",
		ProvenE2EE:    true,
	}
}

func TestnetChat() SKU {
	return SKU{
		Network:       config.Testnet,
		Model:         "qwen/qwen2.5-omni-7b",
		Provider:      "0xa48f01287233509FD694a22Bf840225062E67836",
		TeeSigner:     "0x83df4B8EbA7c0B3B740019b8c9a77ffF77D508cF",
		Verifiability: "TeeML",
		ProvenE2EE:    false, // PIT has not VerifyE2EE'd Omni
	}
}

func ForNetwork(n config.Network) SKU {
	if n == config.Testnet {
		return TestnetChat()
	}
	return MainnetChat()
}

func RefuseSKUCopy(n config.Network, model string) error {
	m := strings.ToLower(strings.TrimSpace(model))
	if n == config.Testnet && (strings.Contains(m, "glm-5.2") || strings.Contains(m, "glm-5.3")) {
		return fmt.Errorf("sku_copy_denied")
	}
	if n == config.Mainnet && strings.Contains(m, "qwen2.5-omni") {
		return fmt.Errorf("sku_copy_denied")
	}
	return nil
}

func SealedAskEnabled(s SKU) error {
	if s.Verifiability != "TeeML" {
		return fmt.Errorf("NOT_TEEML")
	}
	if s.Network == config.Testnet && !s.ProvenE2EE {
		return fmt.Errorf("galileo_e2ee_unproven")
	}
	return nil
}
