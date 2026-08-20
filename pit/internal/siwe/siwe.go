package siwe

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/identity"
)

type Message struct {
	Domain    string
	Address   identity.Address
	URI       string
	Version   string
	ChainID   int64
	Nonce     string
	IssuedAt  time.Time
	Statement string
}

func Build(m Message) (string, error) {
	if m.Domain == "" || m.URI == "" || m.Nonce == "" {
		return "", fmt.Errorf("incomplete siwe")
	}
	if m.Version == "" {
		m.Version = "1"
	}
	if m.Statement == "" {
		m.Statement = "Bind this wallet to your PIT workspace. PIT never asks for a seed phrase."
	}
	issued := m.IssuedAt.UTC().Format(time.RFC3339)
	return fmt.Sprintf(
		"%s wants you to sign in with your Ethereum account:\n%s\n\n%s\n\nURI: %s\nVersion: %s\nChain ID: %d\nNonce: %s\nIssued At: %s",
		m.Domain, strings.ToLower(string(m.Address)), m.Statement, m.URI, m.Version, m.ChainID, m.Nonce, issued,
	), nil
}

func Recover(message, sigHex string) (identity.Address, error) {
	sigHex = strings.TrimPrefix(strings.TrimSpace(sigHex), "0x")
	sig, err := decodeHex(sigHex)
	if err != nil {
		return "", err
	}
	if len(sig) != 65 {
		return "", fmt.Errorf("SIGNATURE_DECLINED")
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	hash := accountsHash([]byte(message))
	pub, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", fmt.Errorf("SIGNATURE_DECLINED")
	}
	addr := crypto.PubkeyToAddress(*pub)
	return identity.NormalizeAddress(addr.Hex())
}

func accountsHash(msg []byte) []byte {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))
	return crypto.Keccak256([]byte(prefix), msg)
}

func decodeHex(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("SIGNATURE_DECLINED")
	}
	out := make([]byte, len(s)/2)
	for i := 0; i < len(out); i++ {
		a := unhex(s[i*2])
		b := unhex(s[i*2+1])
		if a < 0 || b < 0 {
			return nil, fmt.Errorf("SIGNATURE_DECLINED")
		}
		out[i] = byte(a<<4 | b)
	}
	return out, nil
}

func unhex(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	default:
		return -1
	}
}

func AssertChain(got, want int64) error {
	if got != want {
		return fmt.Errorf("WRONG_NETWORK")
	}
	return nil
}

func MustAddress(a common.Address) identity.Address {
	out, _ := identity.NormalizeAddress(a.Hex())
	return out
}
