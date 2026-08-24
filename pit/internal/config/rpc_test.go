package config

import "testing"

func TestRefuseMixedRPC(t *testing.T) {
	if err := RefuseMixedRPC(16661, "https://evmrpc.0g.ai"); err != nil {
		t.Fatal(err)
	}
	if err := RefuseMixedRPC(16661, "https://evmrpc-testnet.0g.ai"); err == nil {
		t.Fatal("mix")
	}
	if err := RefuseMixedRPC(16602, "https://evmrpc.0g.ai"); err == nil {
		t.Fatal("main")
	}
}
