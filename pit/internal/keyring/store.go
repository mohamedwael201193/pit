package keyring

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// Store is the secret backend. Web, Vercel, MCP, and Render never receive a Store.
type Store interface {
	Put(namespace, name string, secret []byte) error
	Get(namespace, name string) ([]byte, error)
	Delete(namespace, name string) error
}

// OpenProduct uses the OS keychain unless PIT_KEYRING=file (tests and recovery).
// Tests default to FileStore so CI does not write OS credentials.
func OpenProduct(root string) (Store, error) {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("PIT_KEYRING")))
	if testing.Testing() && mode == "" {
		mode = "file"
	}
	switch mode {
	case "", "os":
		return OSStore{}, nil
	case "file":
		return Open(root)
	default:
		return nil, fmt.Errorf("unknown_keyring")
	}
}

func BackendName() string {
	if testing.Testing() && strings.TrimSpace(os.Getenv("PIT_KEYRING")) == "" {
		return "file"
	}
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("PIT_KEYRING")))
	if mode == "file" {
		return "file"
	}
	return "os"
}
