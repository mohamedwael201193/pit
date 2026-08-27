package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/config"
)

// AccountProbe is a public eth_call against InferenceServing.getAccount.
// A missing account is generation 0 / not acknowledged — the official SDK catch path.
type AccountProbe struct {
	Generation   uint64 `json:"generation"`
	Acknowledged bool   `json:"acknowledged"`
	Present      bool   `json:"present"`
}

func ProbeDirectAccount(ch config.Chain, user, provider string) AccountProbe {
	out := AccountProbe{}
	parsed, err := abi.JSON(strings.NewReader(getAccountABI))
	if err != nil {
		return out
	}
	userAddr, err := ChecksumAddress(user)
	if err != nil {
		return out
	}
	provAddr, err := ChecksumAddress(provider)
	if err != nil {
		return out
	}
	data, err := parsed.Pack("getAccount", common.HexToAddress(userAddr), common.HexToAddress(provAddr))
	if err != nil {
		return out
	}
	body := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_call",
		"params": []any{
			map[string]string{"to": ch.Serving, "data": "0x" + common.Bytes2Hex(data)},
			"latest",
		},
	}
	raw, _ := json.Marshal(body)
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Post(ch.RPC, "application/json", bytes.NewReader(raw))
	if err != nil {
		return out
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var rpc struct {
		Result string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(b, &rpc); err != nil || rpc.Result == "" || rpc.Result == "0x" {
		return out
	}
	decoded, err := parsed.Unpack("getAccount", common.FromHex(rpc.Result))
	if err != nil || len(decoded) == 0 {
		return out
	}
	vals, ok := decoded[0].([]any)
	if !ok {
		return unpackAccountMap(decoded[0], &out)
	}
	if len(vals) >= 11 {
		out.Present = true
		if ack, ok := vals[7].(bool); ok {
			out.Acknowledged = ack
		}
		out.Generation = asUint64(vals[9])
	}
	return out
}

func unpackAccountMap(v any, out *AccountProbe) AccountProbe {
	_ = fmt.Sprintf("%v", v)
	return *out
}

func asUint64(v any) uint64 {
	switch n := v.(type) {
	case uint8:
		return uint64(n)
	case uint64:
		return n
	case int64:
		if n < 0 {
			return 0
		}
		return uint64(n)
	default:
		if s, ok := n.(interface{ Uint64() uint64 }); ok {
			return s.Uint64()
		}
	}
	return 0
}

const getAccountABI = `[{"name":"getAccount","type":"function","stateMutability":"view","inputs":[{"name":"user","type":"address"},{"name":"provider","type":"address"}],"outputs":[{"name":"account","type":"tuple","components":[{"name":"user","type":"address"},{"name":"provider","type":"address"},{"name":"nonce","type":"uint256"},{"name":"balance","type":"uint256"},{"name":"pendingRefund","type":"uint256"},{"name":"refunds","type":"tuple[]","components":[{"name":"index","type":"uint256"},{"name":"amount","type":"uint256"},{"name":"createdAt","type":"uint256"},{"name":"processed","type":"bool"}]},{"name":"additionalInfo","type":"string"},{"name":"acknowledged","type":"bool"},{"name":"validRefundsLength","type":"uint256"},{"name":"generation","type":"uint256"},{"name":"revokedBitmap","type":"uint256"}]}]}]`
