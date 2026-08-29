// Package receipt builds the public evidence record PIT files on 0G Storage
// after every research verdict and every venue action. The record is public on
// purpose: anyone holding the storage root can download it and recompute the
// digest without a key. Nothing private is allowed inside it.
package receipt

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Version is the schema tag written into every record.
const Version = "pit.receipt.v1"

// Kinds of evidence PIT files.
const (
	KindResearch = "research"
	KindOrder    = "order"
	KindCancel   = "cancel"
)

// Role carries the per-role TEE attestation summary from the sealed committee.
type Role struct {
	Role       string `json:"role"`
	VerifyE2EE string `json:"verify_e2ee"`
	Signer     string `json:"signer,omitempty"`
	TeeSigner  string `json:"tee_signer,omitempty"`
	Side       string `json:"proposed_side,omitempty"`
	Survives   bool   `json:"survives,omitempty"`
	Kill       bool   `json:"kill,omitempty"`
}

// Receipt is the canonical public record. Field order is the canonical byte
// order: encoding/json emits struct fields in declaration order, so marshalling
// the same Receipt twice always produces the same bytes and the same digest.
type Receipt struct {
	V           string  `json:"v"`
	Kind        string  `json:"kind"`
	Network     string  `json:"network"`
	ChainID     int64   `json:"chain_id"`
	Workspace   string  `json:"workspace"`
	Venue       string  `json:"venue,omitempty"`
	Market      string  `json:"market,omitempty"`
	Side        string  `json:"side,omitempty"`
	Size        float64 `json:"size,omitempty"`
	NotionalUSD float64 `json:"notional_usd,omitempty"`
	Verdict     string  `json:"verdict,omitempty"`
	Deny        string  `json:"deny,omitempty"`
	PreviewHash string  `json:"preview_hash,omitempty"`
	PolicyHash  string  `json:"policy_hash,omitempty"`
	OID         string  `json:"oid,omitempty"`
	Cloid       string  `json:"cloid,omitempty"`
	OrderStatus string  `json:"order_status,omitempty"`
	JobID       string  `json:"job_id,omitempty"`
	Model       string  `json:"compute_model,omitempty"`
	Provider    string  `json:"compute_provider,omitempty"`
	TeeSigner   string  `json:"tee_signer,omitempty"`
	Roles       []Role  `json:"roles,omitempty"`
	CreatedAt   string  `json:"created_at"`
	Sign        bool    `json:"sign"`
	Trade       bool    `json:"trade"`
}

// New stamps schema, timestamp, and the two honesty flags that must never be
// true in evidence.
func New(kind, network string, chainID int64, workspace string) Receipt {
	return Receipt{
		V:         Version,
		Kind:      kind,
		Network:   network,
		ChainID:   chainID,
		Workspace: workspace,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Sign:      false,
		Trade:     false,
	}
}

// Canonical returns the exact bytes that get uploaded and digested.
func (r Receipt) Canonical() ([]byte, error) {
	r.V = Version
	r.Sign = false
	r.Trade = false
	if strings.TrimSpace(r.Kind) == "" {
		return nil, fmt.Errorf("receipt_kind_required")
	}
	if strings.TrimSpace(r.Network) == "" || r.ChainID == 0 {
		return nil, fmt.Errorf("receipt_network_required")
	}
	if strings.TrimSpace(r.Workspace) == "" {
		return nil, fmt.Errorf("receipt_workspace_required")
	}
	if strings.TrimSpace(r.CreatedAt) == "" {
		return nil, fmt.Errorf("receipt_time_required")
	}
	raw, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	if err := PublicSafe(raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// Digest is the sha256 of the canonical bytes. A verifier downloads the record
// from 0G Storage and recomputes this to prove the bytes were not swapped.
func (r Receipt) Digest() (string, error) {
	raw, err := r.Canonical()
	if err != nil {
		return "", err
	}
	return DigestBytes(raw), nil
}

// DigestBytes hashes arbitrary downloaded bytes the same way.
func DigestBytes(raw []byte) string {
	sum := sha256.Sum256(raw)
	return "0x" + hex.EncodeToString(sum[:])
}

var forbidden = []string{
	"app-sk-",
	"private_key",
	"privatekey",
	"encryption_key",
	"encryptionkey",
	"session_key",
	"sessionkey",
	"agent_key",
	"agentkey",
	"payer_key",
	"payerkey",
	"mnemonic",
	"passphrase",
	"secret",
	"begin ec private",
	"begin private",
}

// PublicSafe refuses to publish anything that looks like key material. A
// receipt is world-readable once its root exists, so this gate is the last
// thing standing between the user and an irreversible leak.
func PublicSafe(raw []byte) error {
	low := strings.ToLower(string(raw))
	for _, bad := range forbidden {
		if strings.Contains(low, bad) {
			return fmt.Errorf("receipt_not_public_safe")
		}
	}
	if strings.Contains(low, `"sign":true`) || strings.Contains(low, `"trade":true`) {
		return fmt.Errorf("receipt_not_public_safe")
	}
	return nil
}

// Parse reads bytes that came back from 0G Storage.
func Parse(raw []byte) (Receipt, error) {
	var r Receipt
	if err := json.Unmarshal(raw, &r); err != nil {
		return Receipt{}, fmt.Errorf("receipt_unreadable")
	}
	if r.V != Version {
		return Receipt{}, fmt.Errorf("receipt_version_unknown")
	}
	if r.Sign || r.Trade {
		return Receipt{}, fmt.Errorf("receipt_not_public_safe")
	}
	return r, nil
}

// RolesVerified reports whether every listed role carried an OK envelope
// verification and a signer that matched the on-chain TEE signer.
func RolesVerified(roles []Role) bool {
	if len(roles) == 0 {
		return false
	}
	for _, role := range roles {
		if !strings.EqualFold(strings.TrimSpace(role.VerifyE2EE), "OK") {
			return false
		}
		if role.TeeSigner != "" && !strings.EqualFold(role.Signer, role.TeeSigner) {
			return false
		}
	}
	return true
}
