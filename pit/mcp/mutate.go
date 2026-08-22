package mcp

import "fmt"

func MayMutate(tool string) error {
	if TradeDenied(tool) {
		return fmt.Errorf("mcp_read_only")
	}
	return nil
}
