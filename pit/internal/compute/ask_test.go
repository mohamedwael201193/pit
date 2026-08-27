package compute

import (
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
