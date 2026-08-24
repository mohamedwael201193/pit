package config

import (
	"fmt"
	"strings"
)

func RefuseMixedRPC(chainID int64, rpc string) error {
	u := strings.ToLower(strings.TrimSpace(rpc))
	if u == "" || chainID == 0 {
		return fmt.Errorf("rpc_unbound")
	}
	if chainID == 16661 && strings.Contains(u, "testnet") {
		return fmt.Errorf("rpc_mix")
	}
	if chainID == 16602 && strings.Contains(u, "evmrpc.0g.ai") && !strings.Contains(u, "testnet") {
		return fmt.Errorf("rpc_mix")
	}
	return nil
}
