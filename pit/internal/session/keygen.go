package session

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/keyring"
)

type AgentKey struct {
	Address string
	// secret is never exported via JSON.
	secret *ecdsa.PrivateKey
}

func GenerateAgent(ttlHours int) (AgentKey, int64, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return AgentKey{}, 0, err
	}
	addr := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
	exp := TTL(time.Now().UTC(), ttlHours)
	return AgentKey{Address: addr, secret: key}, exp, nil
}

func (k AgentKey) Name(workspaceID string) string {
	short := workspaceID
	if len(short) > 8 {
		short = short[:8]
	}
	return "PIT-" + short
}

func (k AgentKey) Store(ring keyring.Store, workspaceID, sessionID string) error {
	if k.secret == nil {
		return fmt.Errorf("empty_session_key")
	}
	raw := hex.EncodeToString(crypto.FromECDSA(k.secret))
	return ring.Put(workspaceID+"/session", sessionID, []byte(raw))
}

func LoadAgent(ring keyring.Store, workspaceID, sessionID string) (AgentKey, error) {
	b, err := ring.Get(workspaceID+"/session", sessionID)
	if err != nil {
		return AgentKey{}, err
	}
	raw, err := hex.DecodeString(strings.TrimSpace(string(b)))
	if err != nil {
		return AgentKey{}, fmt.Errorf("corrupt_session")
	}
	key, err := crypto.ToECDSA(raw)
	if err != nil {
		return AgentKey{}, err
	}
	addr := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
	return AgentKey{Address: addr, secret: key}, nil
}

func (k AgentKey) ExportJSON() error {
	return fmt.Errorf("session_export_denied")
}

func (k AgentKey) WithSecret(fn func(*ecdsa.PrivateKey) error) error {
	if k.secret == nil {
		return fmt.Errorf("empty_session_key")
	}
	return fn(k.secret)
}

type Permissions struct {
	Order    bool
	Cancel   bool
	Withdraw bool
	Leverage bool
}

func Card() Permissions {
	return Permissions{Order: true, Cancel: true, Withdraw: false, Leverage: false}
}
