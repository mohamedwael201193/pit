package storage

import (
	"fmt"
	"strings"
)

func RequireFlow(addr string) error {
	a := strings.TrimSpace(addr)
	if !strings.HasPrefix(a, "0x") || len(a) != 42 {
		return fmt.Errorf("flow_address_required")
	}
	return nil
}

func UploadMustEncrypt(args []string) error {
	for i, a := range args {
		if a == "--encryption-key" && i+1 < len(args) {
			return RequireHexKey(args[i+1])
		}
	}
	return fmt.Errorf("encryption_key_required")
}
