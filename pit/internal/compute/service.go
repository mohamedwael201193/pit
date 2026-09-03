package compute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/config"
)

// LiveService is InferenceServing.getService. PIT never auto-swaps the frozen SKU.
type LiveService struct {
	Provider      string
	ServiceType   string
	URL           string
	Model         string
	Verifiability string
	TeeSigner     string
	TeeAck        bool
	Present       bool
	Err           string
}

var (
	svcMu  sync.Mutex
	svcAt  time.Time
	svcKey string
	svcVal LiveService
)

func GetService(ch config.Chain, provider string) (LiveService, error) {
	key := strings.ToLower(ch.RPC + "|" + ch.Serving + "|" + provider)
	svcMu.Lock()
	if svcKey == key && time.Since(svcAt) < 12*time.Second {
		v := svcVal
		svcMu.Unlock()
		return v, nil
	}
	svcMu.Unlock()
	got, err := getServiceUncached(ch, provider)
	if err != nil {
		return got, err
	}
	svcMu.Lock()
	svcKey = key
	svcAt = time.Now()
	svcVal = got
	svcMu.Unlock()
	return got, nil
}

func getServiceUncached(ch config.Chain, provider string) (LiveService, error) {
	out := LiveService{}
	parsed, err := abi.JSON(strings.NewReader(getServiceABI))
	if err != nil {
		return out, err
	}
	prov, err := ChecksumAddress(provider)
	if err != nil {
		out.Err = "bad_provider"
		return out, nil
	}
	data, err := parsed.Pack("getService", common.HexToAddress(prov))
	if err != nil {
		return out, err
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
		return out, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var rpc struct {
		Result string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(b, &rpc); err != nil {
		return out, err
	}
	if rpc.Error != nil {
		out.Err = "rpc_error"
		return out, nil
	}
	if rpc.Result == "" || rpc.Result == "0x" {
		out.Err = "sku_unread"
		return out, nil
	}
	decoded, err := parsed.Unpack("getService", common.FromHex(rpc.Result))
	if err != nil || len(decoded) == 0 {
		out.Err = "sku_decode"
		return out, nil
	}
	applyDecodedService(decoded[0], &out)
	return out, nil
}

func applyDecodedService(v any, out *LiveService) {
	vals := asAnySlice(v)
	if len(vals) < 11 {
		out.Err = "sku_shape"
		return
	}
	out.Present = true
	out.Provider = asAddr(vals[0])
	out.ServiceType = asString(vals[1])
	out.URL = asString(vals[2])
	out.Model = asString(vals[6])
	out.Verifiability = asString(vals[7])
	out.TeeSigner = asAddr(vals[9])
	if ack, ok := vals[10].(bool); ok {
		out.TeeAck = ack
	}
}

func asAddr(v any) string {
	switch a := v.(type) {
	case common.Address:
		return a.Hex()
	case [20]byte:
		return common.BytesToAddress(a[:]).Hex()
	case []byte:
		if len(a) == 20 {
			return common.BytesToAddress(a).Hex()
		}
	case string:
		if strings.TrimSpace(a) == "" {
			return ""
		}
		return common.HexToAddress(a).Hex()
	}
	return ""
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// MatchDirectPin fails closed on provider identity. Model is chosen by PickDirectModel.
func MatchDirectPin(live LiveService, pin SKU) error {
	if !live.Present {
		if live.Err != "" {
			return fmt.Errorf("%s", live.Err)
		}
		return fmt.Errorf("sku_unread")
	}
	if !live.TeeAck {
		return fmt.Errorf("unacked_provider")
	}
	if err := RequireTeeML(live.Verifiability); err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(live.URL), pin.URL) {
		return fmt.Errorf("sku_drift_url")
	}
	if !addrEq(live.Provider, pin.Provider) {
		return fmt.Errorf("sku_drift_provider")
	}
	if !addrEq(live.TeeSigner, pin.TeeSigner) {
		return fmt.Errorf("sku_drift_teesigner")
	}
	return nil
}

// MatchFrozenSKU fails closed. It never returns a replacement SKU.
func MatchFrozenSKU(live LiveService, frozen SKU) error {
	if err := MatchDirectPin(live, frozen); err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(live.Model), frozen.Model) {
		return fmt.Errorf("sku_drift_model")
	}
	return nil
}

func addrEq(a, b string) bool {
	if strings.TrimSpace(a) == "" || strings.TrimSpace(b) == "" {
		return false
	}
	return strings.EqualFold(common.HexToAddress(a).Hex(), common.HexToAddress(b).Hex())
}

// FreezeLiveSKU is the sealed-path gate. Transport failure keeps the pinned SKU.
// A successful getService that does not match the Direct pin refuses the ask.
// Model is glm-5.3 only when THIS provider lists it as TeeML; TeeTLS/Router are never used.
func FreezeLiveSKU(net config.Network) (SKU, error) {
	sku := ForNetwork(net)
	if net != config.Mainnet {
		return sku, nil
	}
	got, err := GetService(config.MainnetChain(), sku.Provider)
	if err != nil {
		return sku, nil
	}
	if err := MatchDirectPin(got, sku); err != nil {
		return SKU{}, err
	}
	models, _ := FetchProviderModels(sku.URL)
	model, err := PickDirectModel(got, models)
	if err != nil {
		return SKU{}, err
	}
	sku.Model = model
	sku.TeeSigner = got.TeeSigner
	return sku, nil
}

func fetchPubkeySigner(ctx context.Context, rawURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("pubkey_http_%d", resp.StatusCode)
	}
	var body struct {
		SignerAddress string `json:"signer_address"`
		Signer        string `json:"signer"`
		KemID         string `json:"kem_id"`
	}
	if err := json.Unmarshal(b, &body); err != nil {
		return "", err
	}
	s := strings.TrimSpace(body.SignerAddress)
	if s == "" {
		s = strings.TrimSpace(body.Signer)
	}
	if s == "" {
		return "", fmt.Errorf("pubkey_signer_missing")
	}
	if body.KemID != "" && body.KemID != "0x0020" && !strings.EqualFold(body.KemID, "32") {
		return "", fmt.Errorf("pubkey_kem")
	}
	return s, nil
}

const getServiceABI = `[{"name":"getService","type":"function","stateMutability":"view","inputs":[{"name":"provider","type":"address"}],"outputs":[{"name":"service","type":"tuple","components":[{"name":"provider","type":"address"},{"name":"serviceType","type":"string"},{"name":"url","type":"string"},{"name":"inputPrice","type":"uint256"},{"name":"outputPrice","type":"uint256"},{"name":"updatedAt","type":"uint256"},{"name":"model","type":"string"},{"name":"verifiability","type":"string"},{"name":"additionalInfo","type":"string"},{"name":"teeSignerAddress","type":"address"},{"name":"teeSignerAcknowledged","type":"bool"}]}]}]`
