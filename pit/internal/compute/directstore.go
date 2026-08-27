package compute

import (
	"fmt"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/keyring"
)

func DirectItemName(provider string) (string, error) {
	p, err := ChecksumAddress(provider)
	if err != nil {
		return "", err
	}
	return strings.ToLower(p), nil
}

func StoreDirect(store keyring.Store, network, workspace, provider, authorization string) error {
	if store == nil {
		return fmt.Errorf("keychain")
	}
	if err := RefuseRouterKey(authorization); err != nil {
		return err
	}
	ns, err := keyring.Namespace(network, workspace, "direct")
	if err != nil {
		return err
	}
	name, err := DirectItemName(provider)
	if err != nil {
		return err
	}
	return store.Put(ns, name, []byte(authorization))
}

func LoadDirect(store keyring.Store, network, workspace, provider string) (string, error) {
	if store == nil {
		return "", fmt.Errorf("direct_token_required")
	}
	ns, err := keyring.Namespace(network, workspace, "direct")
	if err != nil {
		return "", err
	}
	name, err := DirectItemName(provider)
	if err != nil {
		return "", err
	}
	b, err := store.Get(ns, name)
	if err != nil || len(b) == 0 {
		return "", fmt.Errorf("direct_token_required")
	}
	auth := string(b)
	if err := RefuseRouterKey(auth); err != nil {
		return "", err
	}
	return auth, nil
}

func DeleteDirect(store keyring.Store, network, workspace, provider string) error {
	if store == nil {
		return nil
	}
	ns, err := keyring.Namespace(network, workspace, "direct")
	if err != nil {
		return err
	}
	name, err := DirectItemName(provider)
	if err != nil {
		return err
	}
	_ = store.Delete(ns, name)
	return nil
}

func AuthFromKeychain(store keyring.Store, net config.Network, workspace string, now time.Time) (AuthFile, DirectMeta, error) {
	sku := ForNetwork(net)
	auth, err := LoadDirect(store, string(net), workspace, sku.Provider)
	if err != nil {
		return AuthFile{}, DirectMeta{}, err
	}
	tok, _, err := ParseBearer(auth)
	if err != nil {
		return AuthFile{}, DirectMeta{}, err
	}
	if TokenExpired(tok, now) {
		return AuthFile{}, DirectMeta{}, fmt.Errorf("direct_token_expired")
	}
	if err := RefuseRouterKey(auth); err != nil {
		return AuthFile{}, DirectMeta{}, err
	}
	if err := DenyRouter(sku.URL); err != nil {
		return AuthFile{}, DirectMeta{}, err
	}
	file := AuthFile{
		Provider:      sku.Provider,
		URL:           sku.URL,
		Model:         sku.Model,
		TeeSigner:     sku.TeeSigner,
		Verifiability: sku.Verifiability,
		Authorization: auth,
	}
	return file, PublicMeta(sku, tok, "keychain"), nil
}

func LoadEnvAuthFile() (AuthFile, error) {
	p := DirectAuthPath()
	if p == "" {
		return AuthFile{}, fmt.Errorf("direct_token_required")
	}
	return ReadAuthFile(p)
}
