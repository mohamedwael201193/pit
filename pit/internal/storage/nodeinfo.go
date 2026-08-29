package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// The 0G storage network answers questions about a root without any PIT code in
// the loop: the indexer lists trusted nodes, and each node reports whether it
// holds a finalized copy of that root. That makes an evidence root checkable by
// anyone with curl, not just by someone running this desk.

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

// NodeFile is the subset of zgs_getFileInfo that proves availability.
type NodeFile struct {
	Node            string `json:"node"`
	Finalized       bool   `json:"finalized"`
	Size            int64  `json:"size"`
	Seq             int64  `json:"seq"`
	StartEntryIndex int64  `json:"start_entry_index"`
	Root            string `json:"root"`
	Error           string `json:"error,omitempty"`
}

func rpc(ctx context.Context, url, method string, params []any, out any) error {
	raw, err := json.Marshal(rpcRequest{JSONRPC: "2.0", Method: method, Params: params, ID: 1})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := (&http.Client{Timeout: 12 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("rpc_http_%d", res.StatusCode)
	}
	var env struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return fmt.Errorf("rpc_unreadable")
	}
	if env.Error != nil {
		return fmt.Errorf("rpc_error")
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(env.Result, out)
}

// TrustedNodes asks the indexer which storage nodes it trusts. No node URL is
// hardcoded so the check keeps working as the network rotates.
func TrustedNodes(ctx context.Context, indexer string, limit int) ([]string, error) {
	if strings.TrimSpace(indexer) == "" {
		return nil, fmt.Errorf("indexer_required")
	}
	var body struct {
		Trusted []struct {
			URL     string `json:"url"`
			Latency int64  `json:"latency"`
		} `json:"trusted"`
	}
	if err := rpc(ctx, indexer, "indexer_getShardedNodes", []any{}, &body); err != nil {
		return nil, err
	}
	nodes := body.Trusted
	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Latency < nodes[j].Latency })
	out := make([]string, 0, len(nodes))
	seen := map[string]bool{}
	for _, n := range nodes {
		u := strings.TrimSpace(n.URL)
		if u == "" || seen[u] {
			continue
		}
		seen[u] = true
		out = append(out, u)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no_trusted_nodes")
	}
	return out, nil
}

// FileInfo asks one storage node whether it holds the root.
func FileInfo(ctx context.Context, node, root string) NodeFile {
	got := NodeFile{Node: node, Root: root}
	if err := RejectBadRoot(root); err != nil {
		got.Error = err.Error()
		return got
	}
	var body struct {
		Finalized bool `json:"finalized"`
		Tx        struct {
			Size            int64  `json:"size"`
			Seq             int64  `json:"seq"`
			StartEntryIndex int64  `json:"startEntryIndex"`
			DataMerkleRoot  string `json:"dataMerkleRoot"`
		} `json:"tx"`
	}
	if err := rpc(ctx, node, "zgs_getFileInfo", []any{root, true}, &body); err != nil {
		got.Error = err.Error()
		return got
	}
	got.Finalized = body.Finalized
	got.Size = body.Tx.Size
	got.Seq = body.Tx.Seq
	got.StartEntryIndex = body.Tx.StartEntryIndex
	if !strings.EqualFold(strings.TrimSpace(body.Tx.DataMerkleRoot), root) {
		got.Error = "root_mismatch"
		got.Finalized = false
	}
	return got
}

// Availability polls the trusted nodes and reports each answer. A caller that
// sees zero finalized copies must not claim the evidence is stored.
func Availability(ctx context.Context, indexer, root string, maxNodes int) ([]NodeFile, error) {
	nodes, err := TrustedNodes(ctx, indexer, maxNodes)
	if err != nil {
		return nil, err
	}
	out := make([]NodeFile, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, FileInfo(ctx, n, root))
	}
	return out, nil
}

// FinalizedCount counts nodes that reported a finalized copy.
func FinalizedCount(list []NodeFile) int {
	n := 0
	for _, item := range list {
		if item.Finalized && item.Error == "" {
			n++
		}
	}
	return n
}
