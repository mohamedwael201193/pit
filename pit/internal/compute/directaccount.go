package compute

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/config"
)

// AccountProbe is a public eth_call against InferenceServing.getAccount.
// A missing account is generation 0 / not acknowledged — the official SDK catch path.
type AccountProbe struct {
	Generation   uint64   `json:"generation"`
	Acknowledged bool     `json:"acknowledged"`
	Present      bool     `json:"present"`
	BalanceWei   *big.Int `json:"-"`
	Err          string   `json:"error,omitempty"`
}

// CommitteeFloorWei is 3 0G. Official Direct rejects a call when locked balance
// is below the provider minimum; a three-role committee needs headroom after the first role.
var CommitteeFloorWei = new(big.Int).Mul(big.NewInt(3), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

func (p AccountProbe) EnoughForCommittee() bool {
	if p.BalanceWei == nil {
		return false
	}
	return p.Acknowledged && p.BalanceWei.Cmp(CommitteeFloorWei) >= 0
}

func (p AccountProbe) BalanceOG() string {
	if p.BalanceWei == nil {
		return "0"
	}
	r := new(big.Rat).SetInt(p.BalanceWei)
	r.Quo(r, new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
	return strings.TrimRight(strings.TrimRight(r.FloatString(4), "0"), ".")
}

var (
	probeMu  sync.Mutex
	probeAt  time.Time
	probeKey string
	probeVal AccountProbe
)

func cloneProbe(p AccountProbe) AccountProbe {
	out := p
	if p.BalanceWei != nil {
		out.BalanceWei = new(big.Int).Set(p.BalanceWei)
	}
	return out
}

func ProbeDirectAccount(ch config.Chain, user, provider string) AccountProbe {
	key := strings.ToLower(ch.RPC + "|" + user + "|" + provider)
	probeMu.Lock()
	if probeKey == key && time.Since(probeAt) < 12*time.Second {
		v := cloneProbe(probeVal)
		probeMu.Unlock()
		return v
	}
	probeMu.Unlock()
	out := probeDirectAccountUncached(ch, user, provider)
	probeMu.Lock()
	probeKey = key
	probeAt = time.Now()
	probeVal = cloneProbe(out)
	probeMu.Unlock()
	return out
}

func probeDirectAccountUncached(ch config.Chain, user, provider string) AccountProbe {
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
		out.Err = "ledger_decode"
		return out
	}
	applyDecodedAccount(decoded[0], &out)
	if !out.Present {
		out.Err = "ledger_shape"
	}
	return out
}

func applyDecodedAccount(v any, out *AccountProbe) {
	vals := asAnySlice(v)
	if len(vals) < 8 {
		return
	}
	out.Present = true
	out.BalanceWei = asBigInt(vals[3])
	if ack, ok := vals[7].(bool); ok {
		out.Acknowledged = ack
	}
	if len(vals) >= 10 {
		out.Generation = asUint64(vals[9])
	}
}

func asAnySlice(v any) []any {
	if vals, ok := v.([]any); ok {
		return vals
	}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		out := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out[i] = rv.Index(i).Interface()
		}
		return out
	case reflect.Struct:
		out := make([]any, rv.NumField())
		for i := 0; i < rv.NumField(); i++ {
			out[i] = rv.Field(i).Interface()
		}
		return out
	default:
		return nil
	}
}

func asBigInt(v any) *big.Int {
	switch n := v.(type) {
	case *big.Int:
		if n == nil {
			return big.NewInt(0)
		}
		return new(big.Int).Set(n)
	case big.Int:
		return new(big.Int).Set(&n)
	case uint64:
		return new(big.Int).SetUint64(n)
	case int64:
		if n < 0 {
			return big.NewInt(0)
		}
		return big.NewInt(n)
	default:
		if s, ok := n.(interface{ String() string }); ok {
			txt := strings.TrimSpace(s.String())
			if parsed, ok := new(big.Int).SetString(strings.TrimPrefix(txt, "0x"), 0); ok {
				return parsed
			}
		}
	}
	return big.NewInt(0)
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
