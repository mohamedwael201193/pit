package sdk

import "github.com/mohamedwael201193/pit/internal/keyring"

func (c Client) KeyringOnWeb() error {
	return keyring.RefuseWeb("web")
}
