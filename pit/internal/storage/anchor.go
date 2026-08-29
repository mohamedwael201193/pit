package storage

import (
	"context"
	"errors"
	"strings"
)

// When the official client submits a log entry, the 0G Storage flow contract
// emits an event whose indexed topics carry the merkle root. That means the link
// between "these bytes" and "this chain transaction" is not something PIT
// asserts: it is readable from any public 0G RPC with no PIT code involved.

// Anchor is the on-chain half of an evidence record.
type Anchor struct {
	Tx          string `json:"tx"`
	Root        string `json:"root"`
	Success     bool   `json:"success"`
	RootInLogs  bool   `json:"root_in_logs"`
	Flow        string `json:"flow,omitempty"`
	FlowMatch   bool   `json:"flow_match"`
	BlockNumber string `json:"block_number,omitempty"`
	From        string `json:"from,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Bound reports whether the transaction really commits this root.
func (a Anchor) Bound() bool {
	return a.Success && a.RootInLogs && a.FlowMatch && a.Error == ""
}

// CheckAnchor reads the transaction receipt from a 0G RPC endpoint and looks for
// the root among the emitted topics.
func CheckAnchor(ctx context.Context, rpcURL, tx, root, flow string) Anchor {
	got := Anchor{Tx: tx, Root: root, Flow: flow}
	if strings.TrimSpace(rpcURL) == "" {
		got.Error = "rpc_required"
		return got
	}
	if err := RejectBadRoot(root); err != nil {
		got.Error = err.Error()
		return got
	}
	if !looksLikeHash32(tx) {
		got.Error = "tx_required"
		return got
	}
	var body struct {
		Status      string `json:"status"`
		BlockNumber string `json:"blockNumber"`
		From        string `json:"from"`
		To          string `json:"to"`
		Logs        []struct {
			Address string   `json:"address"`
			Topics  []string `json:"topics"`
		} `json:"logs"`
	}
	if err := rpc(ctx, rpcURL, "eth_getTransactionReceipt", []any{tx}, &body); err != nil {
		got.Error = err.Error()
		return got
	}
	if strings.TrimSpace(body.BlockNumber) == "" {
		got.Error = "tx_not_found"
		return got
	}
	got.BlockNumber = body.BlockNumber
	got.From = body.From
	got.Success = strings.EqualFold(strings.TrimSpace(body.Status), "0x1")
	want := strings.ToLower(strings.TrimSpace(root))
	for _, log := range body.Logs {
		for _, topic := range log.Topics {
			if strings.ToLower(strings.TrimSpace(topic)) == want {
				got.RootInLogs = true
				if strings.TrimSpace(flow) == "" || strings.EqualFold(strings.TrimSpace(log.Address), strings.TrimSpace(flow)) {
					got.FlowMatch = true
				}
			}
		}
	}
	if strings.TrimSpace(flow) == "" {
		got.FlowMatch = got.RootInLogs
	}
	if !got.RootInLogs {
		got.Error = "root_not_in_transaction_logs"
	}
	return got
}

func looksLikeHash32(v string) bool {
	v = strings.TrimSpace(v)
	if len(v) != 66 || !strings.HasPrefix(v, "0x") {
		return false
	}
	return hex32.MatchString(v)
}

// ErrNoAnchor is returned when a caller demands an anchor that does not exist.
var ErrNoAnchor = errors.New("no_chain_anchor")
