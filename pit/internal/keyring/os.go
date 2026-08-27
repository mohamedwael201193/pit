package keyring

import (
	"encoding/hex"
	"fmt"

	osk "github.com/zalando/go-keyring"
)

const osService = "os.pit.desktop"

// OSStore is Windows Credential Manager, macOS Keychain, or libsecret.
type OSStore struct{}

func item(namespace, name string) string {
	return namespace + "/" + name
}

func (OSStore) Put(namespace, name string, secret []byte) error {
	if namespace == "" || name == "" {
		return fmt.Errorf("empty_keyring_item")
	}
	return osk.Set(osService, item(namespace, name), hex.EncodeToString(secret))
}

func (OSStore) Get(namespace, name string) ([]byte, error) {
	s, err := osk.Get(osService, item(namespace, name))
	if err != nil {
		return nil, fmt.Errorf("not found")
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("corrupt_keyring")
	}
	return b, nil
}

func (OSStore) Delete(namespace, name string) error {
	if err := osk.Delete(osService, item(namespace, name)); err != nil {
		return fmt.Errorf("not found")
	}
	return nil
}
