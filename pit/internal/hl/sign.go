package hl

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
	V int    `json:"v"`
}

type Envelope struct {
	Action    json.RawMessage `json:"action"`
	Signature Signature       `json:"signature"`
	Nonce     int64           `json:"nonce"`
}

func (e Envelope) Signed() bool {
	return e.Signature.V == 27 || e.Signature.V == 28
}

func l1Digest(source string, conn [32]byte) ([]byte, error) {
	bytes32, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}
	uint256, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}
	addr, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}
	domainArgs := abi.Arguments{
		{Type: bytes32}, {Type: bytes32}, {Type: bytes32}, {Type: uint256}, {Type: addr},
	}
	domainPacked, err := domainArgs.Pack(
		crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)")),
		crypto.Keccak256Hash([]byte("Exchange")),
		crypto.Keccak256Hash([]byte("1")),
		big.NewInt(1337),
		common.Address{},
	)
	if err != nil {
		return nil, err
	}
	agentArgs := abi.Arguments{{Type: bytes32}, {Type: bytes32}, {Type: bytes32}}
	agentPacked, err := agentArgs.Pack(
		crypto.Keccak256Hash([]byte("Agent(string source,bytes32 connectionId)")),
		crypto.Keccak256Hash([]byte(source)),
		common.Hash(conn),
	)
	if err != nil {
		return nil, err
	}
	domain := crypto.Keccak256(domainPacked)
	agent := crypto.Keccak256(agentPacked)
	return crypto.Keccak256([]byte{0x19, 0x01}, domain, agent), nil
}

func SignL1(key *ecdsa.PrivateKey, raw json.RawMessage, nonce int64, mainnet bool) (Envelope, error) {
	if key == nil {
		return Envelope{}, fmt.Errorf("empty_session_key")
	}
	if err := AssertActionType(raw); err != nil {
		return Envelope{}, err
	}
	var head struct {
		Type string `json:"type"`
	}
	_ = json.Unmarshal(raw, &head)
	if err := SessionMustNotSign(head.Type); err != nil {
		return Envelope{}, err
	}
	conn, err := ActionHash(raw, uint64(nonce))
	if err != nil {
		return Envelope{}, err
	}
	source := "b"
	if mainnet {
		source = "a"
	}
	digest, err := l1Digest(source, conn)
	if err != nil {
		return Envelope{}, err
	}
	sig, err := crypto.Sign(digest, key)
	if err != nil {
		return Envelope{}, err
	}
	if len(sig) != 65 {
		return Envelope{}, fmt.Errorf("bad_signature")
	}
	env := Envelope{
		Action: raw,
		Nonce:  nonce,
		Signature: Signature{
			R: "0x" + hex.EncodeToString(sig[:32]),
			S: "0x" + hex.EncodeToString(sig[32:64]),
			V: int(sig[64]) + 27,
		},
	}
	got, err := RecoverL1(env, mainnet)
	if err != nil {
		return Envelope{}, err
	}
	want := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
	if got != want {
		return Envelope{}, fmt.Errorf("signer_mismatch")
	}
	return env, nil
}

func RecoverL1(env Envelope, mainnet bool) (string, error) {
	conn, err := ActionHash(env.Action, uint64(env.Nonce))
	if err != nil {
		return "", err
	}
	source := "b"
	if mainnet {
		source = "a"
	}
	digest, err := l1Digest(source, conn)
	if err != nil {
		return "", err
	}
	r, err := hex.DecodeString(strings.TrimPrefix(env.Signature.R, "0x"))
	if err != nil || len(r) != 32 {
		return "", fmt.Errorf("bad_r")
	}
	s, err := hex.DecodeString(strings.TrimPrefix(env.Signature.S, "0x"))
	if err != nil || len(s) != 32 {
		return "", fmt.Errorf("bad_s")
	}
	v := byte(env.Signature.V)
	if v >= 27 {
		v -= 27
	}
	sig := append(append(append([]byte{}, r...), s...), v)
	pub, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return "", err
	}
	return strings.ToLower(crypto.PubkeyToAddress(*pub).Hex()), nil
}
