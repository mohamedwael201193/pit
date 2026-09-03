package deskid

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/chainrpc"
	"github.com/mohamedwael201193/pit/internal/config"
)

const RecordedAgent = "0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52"

// AuthorizeUsageFromOwner sends authorizeUsage from the desk owner wallet.
// It never uses a Hyperliquid session key. Idempotent if user is already authorized.
func AuthorizeUsageFromOwner(ch config.Chain, tokenID uint64, ownerKey, user string) (string, error) {
	from, err := chainrpc.AddressOfKey(ownerKey)
	if err != nil {
		return "", err
	}
	if err := RequireOwner(ch, tokenID, from.Hex()); err != nil {
		return "", err
	}
	ok, err := IsAuthorized(ch, tokenID, user)
	if err != nil {
		return "", err
	}
	if ok {
		return "", fmt.Errorf("already_authorized")
	}
	data, err := EncodeAuthorizeUsage(tokenID, user)
	if err != nil {
		return "", err
	}
	return chainrpc.Send(ch.RPC, ch.ChainID, ownerKey, ch.DeskID, data)
}

func RevokeAuthorizationFromOwner(ch config.Chain, tokenID uint64, ownerKey, user string) (string, error) {
	from, err := chainrpc.AddressOfKey(ownerKey)
	if err != nil {
		return "", err
	}
	if err := RequireOwner(ch, tokenID, from.Hex()); err != nil {
		return "", err
	}
	data, err := EncodeRevokeAuthorization(tokenID, user)
	if err != nil {
		return "", err
	}
	return chainrpc.Send(ch.RPC, ch.ChainID, ownerKey, ch.DeskID, data)
}

func MintFromWallet(ch config.Chain, ownerKey, to, uri, desc string, dataHash common.Hash) (string, error) {
	if _, err := chainrpc.AddressOfKey(ownerKey); err != nil {
		return "", err
	}
	if strings.EqualFold(to, RecordedAgent) {
		return "", fmt.Errorf("session_agent_cannot_own")
	}
	data, err := EncodeMint(to, uri, dataHash, desc)
	if err != nil {
		return "", err
	}
	return chainrpc.Send(ch.RPC, ch.ChainID, ownerKey, ch.DeskID, data)
}
