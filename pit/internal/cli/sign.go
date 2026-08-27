package cli

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/keyring"
	"github.com/mohamedwael201193/pit/internal/session"
)

func SignBound(dir string, live session.Session, network string, raw []byte, nonce int64) (hl.Envelope, error) {
	if live.ID == "" || live.Workspace == "" {
		return hl.Envelope{}, fmt.Errorf("session_expired")
	}
	ring, err := keyring.Open(KeyringDir(dir))
	if err != nil {
		return hl.Envelope{}, err
	}
	key, err := session.LoadAgent(ring, live.Workspace, live.ID)
	if err != nil {
		return hl.Envelope{}, err
	}
	if !strings.EqualFold(key.Address, live.AgentAddr) {
		return hl.Envelope{}, fmt.Errorf("signer_mismatch")
	}
	var env hl.Envelope
	err = key.WithSecret(func(priv *ecdsa.PrivateKey) error {
		signed, e := hl.SignL1(priv, raw, nonce, network == "mainnet")
		if e != nil {
			return e
		}
		env = signed
		return nil
	})
	return env, err
}
