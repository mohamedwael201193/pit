package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/config"
)

func ListServices(ch config.Chain, offset, limit int64) ([]LiveService, int64, error) {
	parsed, err := abi.JSON(strings.NewReader(getAllServicesABI))
	if err != nil {
		return nil, 0, err
	}
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	data, err := parsed.Pack("getAllServices", big.NewInt(offset), big.NewInt(limit))
	if err != nil {
		return nil, 0, err
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
	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Post(ch.RPC, "application/json", bytes.NewReader(raw))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	var rpc struct {
		Result string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(b, &rpc); err != nil {
		return nil, 0, err
	}
	if rpc.Error != nil {
		return nil, 0, fmt.Errorf("rpc_error")
	}
	if rpc.Result == "" || rpc.Result == "0x" {
		return nil, 0, fmt.Errorf("sku_unread")
	}
	decoded, err := parsed.Unpack("getAllServices", common.FromHex(rpc.Result))
	if err != nil || len(decoded) < 2 {
		return nil, 0, fmt.Errorf("sku_decode")
	}
	rows := asAnySlice(decoded[0])
	out := make([]LiveService, 0, len(rows))
	for _, row := range rows {
		var got LiveService
		applyDecodedService(row, &got)
		if got.Present {
			out = append(out, got)
		}
	}
	total := int64(0)
	switch v := decoded[1].(type) {
	case *big.Int:
		if v != nil {
			total = v.Int64()
		}
	}
	return out, total, nil
}

const getAllServicesABI = `[{"name":"getAllServices","type":"function","stateMutability":"view","inputs":[{"name":"offset","type":"uint256"},{"name":"limit","type":"uint256"}],"outputs":[{"name":"services","type":"tuple[]","components":[{"name":"provider","type":"address"},{"name":"serviceType","type":"string"},{"name":"url","type":"string"},{"name":"inputPrice","type":"uint256"},{"name":"outputPrice","type":"uint256"},{"name":"updatedAt","type":"uint256"},{"name":"model","type":"string"},{"name":"verifiability","type":"string"},{"name":"additionalInfo","type":"string"},{"name":"teeSignerAddress","type":"address"},{"name":"teeSignerAcknowledged","type":"bool"}]},{"name":"total","type":"uint256"}]}]`
