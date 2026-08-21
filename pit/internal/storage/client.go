package storage

import (
	"fmt"
	"strings"
)

func RejectUnofficialClient(cliPath string) error {
	low := strings.ToLower(cliPath)
	if cliPath == "" {
		return fmt.Errorf("official_go_client_required")
	}
	if strings.HasSuffix(low, ".ts") || strings.HasSuffix(low, ".js") || strings.Contains(low, "node_modules") {
		return fmt.Errorf("typescript_sdk_forbidden")
	}
	return nil
}

func DownloadMustProve(args []string) error {
	hasProof := false
	for _, a := range args {
		if a == "--proof" {
			hasProof = true
		}
	}
	if !hasProof {
		return fmt.Errorf("proof_flag_required")
	}
	return nil
}

func RejectBadRoot(root string) error {
	if !strings.HasPrefix(root, "0x") || len(root) < 10 {
		return fmt.Errorf("bad_root")
	}
	return nil
}
