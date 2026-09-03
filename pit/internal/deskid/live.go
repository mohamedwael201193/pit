package deskid

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/chainrpc"
	"github.com/mohamedwael201193/pit/internal/config"
)

type Live struct {
	Owner             string
	URI               string
	OwnerAuthorized   bool
	UserAuthorized    bool
	AuthorizedUsers   []string
	Supports7857      bool
	SupportsAuthorize bool
	SupportsCloneable bool
	Supports721       bool
}

func OwnerOf(ch config.Chain, tokenID uint64) (string, error) {
	data, err := EncodeOwnerOf(tokenID)
	if err != nil {
		return "", err
	}
	out, err := chainrpc.Call(ch.RPC, ch.DeskID, CalldataHex(data))
	if err != nil {
		return "", err
	}
	if len(out) < 32 {
		return "", fmt.Errorf("owner_unread")
	}
	return common.BytesToAddress(out[12:]).Hex(), nil
}

func IsAuthorized(ch config.Chain, tokenID uint64, user string) (bool, error) {
	data, err := EncodeIsAuthorized(tokenID, user)
	if err != nil {
		return false, err
	}
	out, err := chainrpc.Call(ch.RPC, ch.DeskID, CalldataHex(data))
	if err != nil {
		return false, err
	}
	return len(out) > 0 && out[len(out)-1] == 1, nil
}

func Supports(ch config.Chain, id string) (bool, error) {
	data, err := EncodeSupportsInterface(id)
	if err != nil {
		return false, err
	}
	out, err := chainrpc.Call(ch.RPC, ch.DeskID, CalldataHex(data))
	if err != nil {
		return false, err
	}
	return len(out) > 0 && out[len(out)-1] == 1, nil
}

func TokenURI(ch config.Chain, tokenID uint64) (string, error) {
	data, err := EncodeTokenURI(tokenID)
	if err != nil {
		return "", err
	}
	out, err := chainrpc.Call(ch.RPC, ch.DeskID, CalldataHex(data))
	if err != nil {
		return "", err
	}
	vals, err := deskABI.Unpack("tokenURI", out)
	if err != nil || len(vals) == 0 {
		return "", fmt.Errorf("uri_unread")
	}
	s, _ := vals[0].(string)
	return s, nil
}

func AuthorizedUsers(ch config.Chain, tokenID uint64) ([]string, error) {
	data, err := EncodeAuthorizedUsersOf(tokenID)
	if err != nil {
		return nil, err
	}
	out, err := chainrpc.Call(ch.RPC, ch.DeskID, CalldataHex(data))
	if err != nil {
		return nil, err
	}
	vals, err := deskABI.Unpack("authorizedUsersOf", out)
	if err != nil || len(vals) == 0 {
		return nil, fmt.Errorf("authz_unread")
	}
	addrs, _ := vals[0].([]common.Address)
	var outAddr []string
	for _, a := range addrs {
		outAddr = append(outAddr, a.Hex())
	}
	return outAddr, nil
}

func Snapshot(ch config.Chain, tokenID uint64, user string) (Live, error) {
	var live Live
	owner, err := OwnerOf(ch, tokenID)
	if err != nil {
		return live, err
	}
	live.Owner = owner
	live.URI, _ = TokenURI(ch, tokenID)
	live.OwnerAuthorized, _ = IsAuthorized(ch, tokenID, owner)
	if strings.TrimSpace(user) != "" {
		live.UserAuthorized, _ = IsAuthorized(ch, tokenID, user)
	}
	live.AuthorizedUsers, _ = AuthorizedUsers(ch, tokenID)
	live.Supports7857, _ = Supports(ch, ID7857)
	live.SupportsAuthorize, _ = Supports(ch, ID7857Authorize)
	live.SupportsCloneable, _ = Supports(ch, ID7857Cloneable)
	live.Supports721, _ = Supports(ch, ID721)
	return live, nil
}

func RequireOwner(ch config.Chain, tokenID uint64, caller string) error {
	owner, err := OwnerOf(ch, tokenID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(owner, caller) {
		return fmt.Errorf("not_desk_owner")
	}
	return nil
}

func DecodeUint(out []byte) *big.Int {
	return new(big.Int).SetBytes(out)
}
