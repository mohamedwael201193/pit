// Package evidence resolves the 0G account that pays gas when PIT publishes a
// public receipt. This key only funds storage log entries on 0G Chain. It never
// touches the venue, never signs an order, and is deliberately separate from the
// workspace memory key so a public filing can never be encrypted by accident.
package evidence

import (
	"fmt"
	"os"
	"strings"

	"github.com/mohamedwael201193/pit/internal/keyring"
	"github.com/mohamedwael201193/pit/internal/storage"
)

const keyName = "og-evidence-payer"

// PayerKey looks in the environment first so an operator can run the desk
// without storing anything, then in the OS keyring where the desktop keeps it.
func PayerKey(dir, network, workspaceID string) (string, error) {
	if raw := strings.TrimSpace(os.Getenv("PIT_OG_PAYER_KEY")); raw != "" {
		return storage.NormalizePayerKey(raw)
	}
	ns, err := keyring.Namespace(network, workspaceID, "evidence")
	if err != nil {
		return "", fmt.Errorf("payer_key_missing")
	}
	store, err := keyring.OpenProduct(dir)
	if err != nil {
		return "", fmt.Errorf("payer_key_missing")
	}
	secret, err := store.Get(ns, keyName)
	if err != nil || len(secret) == 0 {
		return "", fmt.Errorf("payer_key_missing")
	}
	return storage.NormalizePayerKey(strings.TrimSpace(string(secret)))
}

// SavePayerKey stores the key in the OS keyring for this workspace only.
func SavePayerKey(dir, network, workspaceID, raw string) error {
	norm, err := storage.NormalizePayerKey(raw)
	if err != nil {
		return err
	}
	ns, err := keyring.Namespace(network, workspaceID, "evidence")
	if err != nil {
		return err
	}
	store, err := keyring.OpenProduct(dir)
	if err != nil {
		return err
	}
	return store.Put(ns, keyName, []byte(norm))
}

// Present reports whether a payer key is reachable without returning it.
func Present(dir, network, workspaceID string) bool {
	_, err := PayerKey(dir, network, workspaceID)
	return err == nil
}
