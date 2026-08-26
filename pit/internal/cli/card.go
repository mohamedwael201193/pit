package cli

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func PreviewPath(dir string) string {
	return filepath.Join(dir, "preview.json")
}

func NewCloid() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(b[:]), nil
}

func SavePreview(dir string, p engine.Preview) (string, error) {
	bound, err := engine.BindPreview(p, map[string]any{"sz": 1e9, "side": "sell"})
	if err != nil {
		return "", err
	}
	h, err := engine.CanonicalHash(bound)
	if err != nil {
		return "", err
	}
	body, err := json.MarshalIndent(bound, "", "  ")
	if err != nil {
		return "", err
	}
	if err := RefuseSessionSecret(body); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(PreviewPath(dir), body, 0o600); err != nil {
		return "", err
	}
	return "0x" + h, nil
}

func LoadPreview(dir string) (engine.Preview, string, error) {
	b, err := os.ReadFile(PreviewPath(dir))
	if err != nil {
		return engine.Preview{}, "", fmt.Errorf("preview_required")
	}
	if err := RefuseSessionSecret(b); err != nil {
		return engine.Preview{}, "", err
	}
	var p engine.Preview
	if err := json.Unmarshal(b, &p); err != nil {
		return engine.Preview{}, "", err
	}
	h, err := engine.CanonicalHash(p)
	if err != nil {
		return engine.Preview{}, "", err
	}
	return p, "0x" + h, nil
}
