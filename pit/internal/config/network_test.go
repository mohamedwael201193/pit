package config

import (
	"os"
	"testing"
)

func TestParseNetwork(t *testing.T) {
	n, err := ParseNetwork("mainnet")
	if err != nil || n != Mainnet {
		t.Fatalf("mainnet: %v %v", n, err)
	}
	n, err = ParseNetwork("16602")
	if err != nil || n != Testnet {
		t.Fatalf("testnet: %v %v", n, err)
	}
	if _, err := ParseNetwork("ethereum"); err == nil {
		t.Fatal("expected error")
	}
}

func TestChainsDoNotShareIDs(t *testing.T) {
	m, g := MainnetChain(), TestnetChain()
	if m.ChainID == g.ChainID {
		t.Fatal("networks must not share chain id")
	}
	if m.RPC == g.RPC || m.HLInfo == g.HLInfo {
		t.Fatal("networks must not share RPC or book endpoints")
	}
	if m.Serving == g.Serving {
		t.Fatal("compute contracts must be chain-specific")
	}
}

func TestGuardFallbacks(t *testing.T) {
	t.Setenv("PIT_ALLOW_FALLBACKS", "false")
	if err := GuardFallbacks(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PIT_ALLOW_FALLBACKS", "true")
	if err := GuardFallbacks(); err == nil {
		t.Fatal("mocks must fail closed")
	}
	os.Unsetenv("PIT_ALLOW_FALLBACKS")
	if err := GuardFallbacks(); err != nil {
		t.Fatal(err)
	}
}
