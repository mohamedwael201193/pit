package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Receipt struct {
	PreviewHash string
	StorageRoot string
	Network     string
	Workspace   string
}

func HashPreview(previewJSON []byte) string {
	sum := sha256.Sum256(previewJSON)
	return "0x" + hex.EncodeToString(sum[:])
}

func Check(r Receipt) error {
	if !strings.HasPrefix(r.PreviewHash, "0x") || len(r.PreviewHash) != 66 {
		return fmt.Errorf("bad_preview_hash")
	}
	if !strings.HasPrefix(r.StorageRoot, "0x") || len(r.StorageRoot) != 66 {
		return fmt.Errorf("bad_storage_root")
	}
	if r.Network == "" || r.Workspace == "" {
		return fmt.Errorf("receipt_unbound")
	}
	return nil
}
