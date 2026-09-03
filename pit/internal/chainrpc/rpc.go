package chainrpc

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// Call is a public eth_call. It never signs.
func Call(rpc, to, data string) ([]byte, error) {
	raw, err := rpcDo(rpc, "eth_call", []any{
		map[string]string{"to": to, "data": data},
		"latest",
	})
	if err != nil {
		return nil, err
	}
	if raw == "" || raw == "0x" {
		return nil, fmt.Errorf("empty_call")
	}
	return common.FromHex(raw), nil
}

func Code(rpc, addr string) ([]byte, error) {
	raw, err := rpcDo(rpc, "eth_getCode", []any{addr, "latest"})
	if err != nil {
		return nil, err
	}
	return common.FromHex(raw), nil
}

func StorageAt(rpc, addr, slot string) ([]byte, error) {
	raw, err := rpcDo(rpc, "eth_getStorageAt", []any{addr, slot, "latest"})
	if err != nil {
		return nil, err
	}
	return common.FromHex(raw), nil
}

func parseKey(hexKey string) (*ecdsa.PrivateKey, error) {
	k := strings.TrimSpace(hexKey)
	k = strings.TrimPrefix(strings.TrimPrefix(k, "0x"), "0X")
	if len(k) != 64 {
		return nil, fmt.Errorf("key_required")
	}
	return crypto.HexToECDSA(k)
}

// AddressOfKey returns the EOA for a hex private key. The key is never logged.
func AddressOfKey(hexKey string) (common.Address, error) {
	key, err := parseKey(hexKey)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(key.PublicKey), nil
}

// Send signs a legacy EIP-155 tx from hexKey and broadcasts it. hexKey is never logged.
func Send(rpc string, chainID int64, hexKey, to string, data []byte) (string, error) {
	key, err := parseKey(hexKey)
	if err != nil {
		return "", err
	}
	from := crypto.PubkeyToAddress(key.PublicKey)
	toAddr := common.HexToAddress(to)
	nonceHex, err := rpcDo(rpc, "eth_getTransactionCount", []any{from.Hex(), "pending"})
	if err != nil {
		return "", err
	}
	nonce, err := hexutil.DecodeUint64(nonceHex)
	if err != nil {
		return "", err
	}
	gasPriceHex, err := rpcDo(rpc, "eth_gasPrice", []any{})
	if err != nil {
		return "", err
	}
	gasPrice, ok := new(big.Int).SetString(strings.TrimPrefix(gasPriceHex, "0x"), 16)
	if !ok {
		return "", fmt.Errorf("gas_price")
	}
	estHex, err := rpcDo(rpc, "eth_estimateGas", []any{map[string]string{
		"from": from.Hex(), "to": toAddr.Hex(), "data": "0x" + common.Bytes2Hex(data),
	}})
	gas := uint64(250000)
	if err == nil {
		if g, e := hexutil.DecodeUint64(estHex); e == nil && g > 21000 {
			gas = g + g/5
		}
	}
	raw, err := signLegacyEIP155(key, chainID, nonce, gasPrice, gas, toAddr, data)
	if err != nil {
		return "", err
	}
	hash, err := rpcDo(rpc, "eth_sendRawTransaction", []any{hexutil.Encode(raw)})
	if err != nil {
		return "", err
	}
	if len(hash) != 66 || !strings.HasPrefix(hash, "0x") {
		return "", fmt.Errorf("bad_tx_hash")
	}
	return hash, nil
}

func signLegacyEIP155(key *ecdsa.PrivateKey, chainID int64, nonce uint64, gasPrice *big.Int, gas uint64, to common.Address, data []byte) ([]byte, error) {
	cid := big.NewInt(chainID)
	toBytes := to.Bytes()
	value := big.NewInt(0)
	sighash, err := rlp.EncodeToBytes([]any{nonce, gasPrice, gas, toBytes, value, data, cid, uint(0), uint(0)})
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(crypto.Keccak256(sighash), key)
	if err != nil {
		return nil, err
	}
	v := new(big.Int).Mul(cid, big.NewInt(2))
	v.Add(v, big.NewInt(int64(sig[64])+35))
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	return rlp.EncodeToBytes([]any{nonce, gasPrice, gas, toBytes, value, data, v, r, s})
}

func rpcDo(rpc, method string, params []any) (string, error) {
	body, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Post(rpc, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var out struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return "", err
	}
	if out.Error != nil {
		return "", fmt.Errorf("%s", out.Error.Message)
	}
	var s string
	if json.Unmarshal(out.Result, &s) == nil {
		return s, nil
	}
	return strings.TrimSpace(string(out.Result)), nil
}
