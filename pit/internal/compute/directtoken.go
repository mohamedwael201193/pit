package compute

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/identity"
)

// Official 0G Direct session token (serving-broker SessionToken / SDK generateSessionToken).
// Field order matches JSON.stringify insertion order. Do not reorder.
const (
	EphemeralTokenID            = 255
	EphemeralTokenMaxDurationMS = int64(24 * 60 * 60 * 1000)
)

type SessionToken struct {
	Address    string `json:"address"`
	Provider   string `json:"provider"`
	Timestamp  int64  `json:"timestamp"`
	ExpiresAt  int64  `json:"expiresAt"`
	Nonce      string `json:"nonce"`
	Generation uint64 `json:"generation"`
	TokenId    uint8  `json:"tokenId"`
}

type Challenge struct {
	Message   string `json:"message"`
	Digest    string `json:"digest"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Network   string `json:"network"`
	Wallet    string `json:"wallet"`
	ExpiresAt int64  `json:"expiresAt"`
	TokenId   uint8  `json:"tokenId"`
	Explain   string `json:"explain"`
}

type DirectMeta struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Wallet    string `json:"wallet"`
	ExpiresAt int64  `json:"expiresAt"`
	TokenId   uint8  `json:"tokenId"`
	Source    string `json:"source"`
}

func ChecksumAddress(raw string) (string, error) {
	addr, err := identity.NormalizeAddress(raw)
	if err != nil {
		return "", err
	}
	return common.HexToAddress(string(addr)).Hex(), nil
}

func CanonicalJSON(t SessionToken) string {
	return fmt.Sprintf(
		`{"address":%s,"provider":%s,"timestamp":%d,"expiresAt":%d,"nonce":%s,"generation":%d,"tokenId":%d}`,
		jsonString(t.Address),
		jsonString(t.Provider),
		t.Timestamp,
		t.ExpiresAt,
		jsonString(t.Nonce),
		t.Generation,
		t.TokenId,
	)
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func DigestHex(message string) string {
	sum := crypto.Keccak256([]byte(message))
	return "0x" + hex.EncodeToString(sum)
}

func NewChallenge(wallet, provider string, generation uint64, now time.Time) (SessionToken, Challenge, error) {
	addr, err := ChecksumAddress(wallet)
	if err != nil {
		return SessionToken{}, Challenge{}, err
	}
	prov, err := ChecksumAddress(provider)
	if err != nil {
		return SessionToken{}, Challenge{}, err
	}
	nonce, err := randomNonce()
	if err != nil {
		return SessionToken{}, Challenge{}, err
	}
	ts := now.UTC().UnixMilli()
	tok := SessionToken{
		Address:    addr,
		Provider:   prov,
		Timestamp:  ts,
		ExpiresAt:  ts + EphemeralTokenMaxDurationMS,
		Nonce:      nonce,
		Generation: generation,
		TokenId:    EphemeralTokenID,
	}
	msg := CanonicalJSON(tok)
	return tok, Challenge{
		Message:   msg,
		Digest:    DigestHex(msg),
		Provider:  prov,
		Wallet:    addr,
		ExpiresAt: tok.ExpiresAt,
		TokenId:   tok.TokenId,
		Explain:   "This signature authorizes sealed 0G Direct research for 24 hours. It cannot withdraw funds. It cannot place a Hyperliquid order. The website never receives the assembled token.",
	}, nil
}

func randomNonce() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func AssembleBearer(message, signature string) (string, error) {
	sig := normalizeSigHex(signature)
	if sig == "" {
		return "", fmt.Errorf("SIGNATURE_DECLINED")
	}
	raw := message + "|" + sig
	auth := "Bearer app-sk-" + base64.StdEncoding.EncodeToString([]byte(raw))
	if err := RefuseRouterKey(auth); err != nil {
		return "", err
	}
	return auth, nil
}

func ParseBearer(auth string) (SessionToken, string, error) {
	a := strings.TrimSpace(auth)
	if !strings.HasPrefix(a, "Bearer app-sk-") {
		return SessionToken{}, "", fmt.Errorf("direct_token_required")
	}
	enc := strings.TrimPrefix(a, "Bearer app-sk-")
	decoded, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return SessionToken{}, "", fmt.Errorf("direct_token_required")
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return SessionToken{}, "", fmt.Errorf("direct_token_required")
	}
	var tok SessionToken
	if err := json.Unmarshal([]byte(parts[0]), &tok); err != nil {
		return SessionToken{}, "", fmt.Errorf("direct_token_required")
	}
	return tok, parts[1], nil
}

func RecoverDirectSigner(message, signature string) (identity.Address, error) {
	sig, err := decodeSig(signature)
	if err != nil {
		return "", err
	}
	messageHash := crypto.Keccak256Hash([]byte(message))
	prefixed := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), messageHash.Bytes())
	pub, err := crypto.SigToPub(prefixed.Bytes(), sig)
	if err != nil {
		return "", fmt.Errorf("SIGNATURE_DECLINED")
	}
	return identity.NormalizeAddress(crypto.PubkeyToAddress(*pub).Hex())
}

func AcceptDirectSignature(message, signature, boundWallet string, now time.Time) (string, SessionToken, error) {
	if strings.TrimSpace(message) == "" {
		return "", SessionToken{}, fmt.Errorf("direct_challenge_required")
	}
	got, err := RecoverDirectSigner(message, signature)
	if err != nil {
		return "", SessionToken{}, err
	}
	want, err := identity.NormalizeAddress(boundWallet)
	if err != nil {
		return "", SessionToken{}, fmt.Errorf("wallet_required")
	}
	if got != want {
		return "", SessionToken{}, fmt.Errorf("signature_mismatch")
	}
	var tok SessionToken
	if err := json.Unmarshal([]byte(message), &tok); err != nil {
		return "", SessionToken{}, fmt.Errorf("direct_token_required")
	}
	if CanonicalJSON(tok) != message {
		return "", SessionToken{}, fmt.Errorf("direct_token_required")
	}
	claimed, err := identity.NormalizeAddress(tok.Address)
	if err != nil || claimed != want {
		return "", SessionToken{}, fmt.Errorf("signature_mismatch")
	}
	if tok.TokenId != EphemeralTokenID {
		return "", SessionToken{}, fmt.Errorf("direct_token_required")
	}
	if tok.ExpiresAt <= 0 || tok.ExpiresAt-tok.Timestamp > EphemeralTokenMaxDurationMS {
		return "", SessionToken{}, fmt.Errorf("direct_token_required")
	}
	if now.Unix() > tok.ExpiresAt/1000 {
		return "", SessionToken{}, fmt.Errorf("direct_token_expired")
	}
	auth, err := AssembleBearer(message, signature)
	if err != nil {
		return "", SessionToken{}, err
	}
	return auth, tok, nil
}

func TokenExpired(tok SessionToken, now time.Time) bool {
	if tok.ExpiresAt <= 0 {
		return true
	}
	return now.Unix() > tok.ExpiresAt/1000
}

func PublicMeta(sku SKU, tok SessionToken, source string) DirectMeta {
	return DirectMeta{
		Provider:  sku.Provider,
		Model:     sku.Model,
		Wallet:    tok.Address,
		ExpiresAt: tok.ExpiresAt,
		TokenId:   tok.TokenId,
		Source:    source,
	}
}

func normalizeSigHex(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		s = "0x" + s
	}
	return s
}

func decodeSig(sigHex string) ([]byte, error) {
	s := strings.TrimPrefix(strings.TrimSpace(sigHex), "0x")
	s = strings.TrimPrefix(s, "0X")
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 65 {
		return nil, fmt.Errorf("SIGNATURE_DECLINED")
	}
	if b[64] >= 27 {
		b[64] -= 27
	}
	return b, nil
}
