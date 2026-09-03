package deskid

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Canonical ERC-165 interface IDs used by PitDeskID on Aristotle.
const (
	ID7857          = "0x2afbede9"
	ID7857Authorize = "0xdf597d99"
	ID7857Cloneable = "0x74f8628b"
	ID721           = "0x80ac58cd"
	MainnetTokenID  = 1
)

var deskABI = mustABI(`[
  {"type":"function","name":"mint","inputs":[{"name":"to","type":"address"},{"name":"uri","type":"string"},{"name":"dataHash","type":"bytes32"},{"name":"dataDescription","type":"string"}],"outputs":[{"type":"uint256"}]},
  {"type":"function","name":"authorizeUsage","inputs":[{"name":"tokenId","type":"uint256"},{"name":"user","type":"address"}]},
  {"type":"function","name":"revokeAuthorization","inputs":[{"name":"tokenId","type":"uint256"},{"name":"user","type":"address"}]},
  {"type":"function","name":"isAuthorized","inputs":[{"name":"tokenId","type":"uint256"},{"name":"user","type":"address"}],"outputs":[{"type":"bool"}],"stateMutability":"view"},
  {"type":"function","name":"ownerOf","inputs":[{"name":"tokenId","type":"uint256"}],"outputs":[{"type":"address"}],"stateMutability":"view"},
  {"type":"function","name":"authorizedUsersOf","inputs":[{"name":"tokenId","type":"uint256"}],"outputs":[{"type":"address[]"}],"stateMutability":"view"},
  {"type":"function","name":"supportsInterface","inputs":[{"name":"interfaceId","type":"bytes4"}],"outputs":[{"type":"bool"}],"stateMutability":"view"},
  {"type":"function","name":"tokenURI","inputs":[{"name":"tokenId","type":"uint256"}],"outputs":[{"type":"string"}],"stateMutability":"view"}
]`)

func mustABI(raw string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		panic(err)
	}
	return parsed
}

func Selector(sig string) string {
	h := crypto.Keccak256([]byte(sig))
	return "0x" + common.Bytes2Hex(h[:4])
}

func EncodeMint(to, uri string, dataHash common.Hash, desc string) ([]byte, error) {
	if to == "" || uri == "" || dataHash == (common.Hash{}) {
		return nil, fmt.Errorf("mint_args")
	}
	return deskABI.Pack("mint", common.HexToAddress(to), uri, dataHash, desc)
}

func EncodeAuthorizeUsage(tokenID uint64, user string) ([]byte, error) {
	if user == "" || common.HexToAddress(user) == (common.Address{}) {
		return nil, fmt.Errorf("user_required")
	}
	return deskABI.Pack("authorizeUsage", new(big.Int).SetUint64(tokenID), common.HexToAddress(user))
}

func EncodeRevokeAuthorization(tokenID uint64, user string) ([]byte, error) {
	if user == "" || common.HexToAddress(user) == (common.Address{}) {
		return nil, fmt.Errorf("user_required")
	}
	return deskABI.Pack("revokeAuthorization", new(big.Int).SetUint64(tokenID), common.HexToAddress(user))
}

func EncodeOwnerOf(tokenID uint64) ([]byte, error) {
	return deskABI.Pack("ownerOf", new(big.Int).SetUint64(tokenID))
}

func EncodeIsAuthorized(tokenID uint64, user string) ([]byte, error) {
	return deskABI.Pack("isAuthorized", new(big.Int).SetUint64(tokenID), common.HexToAddress(user))
}

func EncodeSupportsInterface(id string) ([]byte, error) {
	b := common.FromHex(id)
	if len(b) != 4 {
		return nil, fmt.Errorf("interface_id")
	}
	var id4 [4]byte
	copy(id4[:], b)
	return deskABI.Pack("supportsInterface", id4)
}

func EncodeAuthorizedUsersOf(tokenID uint64) ([]byte, error) {
	return deskABI.Pack("authorizedUsersOf", new(big.Int).SetUint64(tokenID))
}

func EncodeTokenURI(tokenID uint64) ([]byte, error) {
	return deskABI.Pack("tokenURI", new(big.Int).SetUint64(tokenID))
}

func CalldataHex(data []byte) string {
	return "0x" + common.Bytes2Hex(data)
}
