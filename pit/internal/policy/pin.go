package policy

import (
	"fmt"
	"os"
	"path/filepath"
)

func PinFile(dir, workspaceID, hash string) (string, error) {
	if workspaceID == "" || hash == "" {
		return "", fmt.Errorf("pin_incomplete")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, workspaceID+".policy")
	if err := os.WriteFile(path, []byte(hash+"\n"), 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func ReadPin(dir, workspaceID string) (string, error) {
	b, err := os.ReadFile(filepath.Join(dir, workspaceID+".policy"))
	if err != nil {
		return "", err
	}
	return string(bytesTrim(b)), nil
}

func bytesTrim(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r' || b[len(b)-1] == ' ') {
		b = b[:len(b)-1]
	}
	return b
}
