package compute

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestBindSealMainnet(t *testing.T) {
	sku := MainnetChat()
	_, env, err := BindSeal(SealJob{
		Network:      config.Mainnet,
		Role:         Researcher,
		PublicMarket: []byte(`{"mark":1}`),
		PrivateBook:  []byte(`{"pos":1}`),
		Scheme:       SchemeE2EE,
		ProviderURL:  sku.URL,
		TeeSigner:    sku.TeeSigner,
		SealerBin:    "/opt/pit/sealer",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(env), "researcher") {
		t.Fatal(string(env))
	}
}

func TestBindSealRejectsRouter(t *testing.T) {
	sku := MainnetChat()
	_, _, err := BindSeal(SealJob{
		Network:      config.Mainnet,
		Role:         Researcher,
		PublicMarket: []byte("m"),
		PrivateBook:  []byte("b"),
		Scheme:       SchemeE2EE,
		ProviderURL:  "https://router-api.0g.ai/v1",
		TeeSigner:    sku.TeeSigner,
		SealerBin:    "/opt/pit/sealer",
	})
	if err == nil {
		t.Fatal("router")
	}
}

func TestBindSealGalileoUnproven(t *testing.T) {
	sku := TestnetChat()
	_, _, err := BindSeal(SealJob{
		Network:      config.Testnet,
		Role:         Researcher,
		PublicMarket: []byte("m"),
		PrivateBook:  []byte("b"),
		Scheme:       SchemeE2EE,
		ProviderURL:  "https://compute.example.test",
		TeeSigner:    sku.TeeSigner,
		SealerBin:    "/opt/pit/sealer",
	})
	if err == nil {
		t.Fatal("galileo must stay disabled until VerifyE2EE is proven")
	}
}

func TestBindSealWrongScheme(t *testing.T) {
	sku := MainnetChat()
	_, _, err := BindSeal(SealJob{
		Network:      config.Mainnet,
		Role:         Risk,
		PublicMarket: []byte("m"),
		PrivateBook:  []byte("b"),
		Scheme:       "plain",
		ProviderURL:  sku.URL,
		TeeSigner:    sku.TeeSigner,
		SealerBin:    "/opt/pit/sealer",
	})
	if err == nil {
		t.Fatal("scheme")
	}
}
