package keyring

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileStore is the test and PIT_KEYRING=file recovery backend.
// Production CLI and desktop use OSStore (Windows Credential Manager / macOS Keychain / libsecret).
type FileStore struct {
	root string
	mu   sync.Mutex
}

func Open(root string) (*FileStore, error) {
	if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, err
	}
	return &FileStore{root: root}, nil
}

func (s *FileStore) path(namespace, name string) string {
	return filepath.Join(s.root, namespace, name)
}

func (s *FileStore) Put(namespace, name string, secret []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	dir := filepath.Join(s.root, namespace)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.WriteFile(s.path(namespace, name), secret, 0o600)
}

func (s *FileStore) Get(namespace, name string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := os.ReadFile(s.path(namespace, name))
	if err != nil {
		return nil, fmt.Errorf("not found")
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out, nil
}

func (s *FileStore) Delete(namespace, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.Remove(s.path(namespace, name)); err != nil {
		return fmt.Errorf("not found")
	}
	return nil
}

func NewMemoryKey() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(b[:]), nil
}
