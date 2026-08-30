package storage

import (
	"fmt"
	"path/filepath"
	"strings"
)

func RejectUnofficialClient(cliPath string) error {
	if strings.TrimSpace(cliPath) == "" {
		return fmt.Errorf("official_go_client_required")
	}
	low := strings.ToLower(cliPath)
	if strings.HasSuffix(low, ".ts") || strings.HasSuffix(low, ".js") || strings.Contains(low, "node_modules") {
		return fmt.Errorf("typescript_sdk_forbidden")
	}
	base := strings.ToLower(filepath.Base(cliPath))
	if base != "0g-storage-client" && base != "0g-storage-client.exe" {
		return fmt.Errorf("official_go_client_required")
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
