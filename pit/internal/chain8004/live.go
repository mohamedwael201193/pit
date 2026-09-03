package chain8004

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/chainrpc"
	"github.com/mohamedwael201193/pit/internal/config"
)

const eip1967ImplSlot = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"

type IdentityLive struct {
	Owner string
	URI   string
}

type FeedbackRead struct {
	Value         *big.Int
	ValueDecimals uint8
	Tag1          string
	Tag2          string
	Revoked       bool
}

type Summary struct {
	Count         uint64
	Value         *big.Int
	ValueDecimals uint8
}

func OwnerOf(ch config.Chain, agentID uint64) (string, error) {
	data, err := EncodeOwnerOf(agentID)
	if err != nil {
		return "", err
	}
	out, err := chainrpc.Call(ch.RPC, ch.Identity8004, CalldataHex(data))
	if err != nil {
		return "", err
	}
	if len(out) < 32 {
		return "", fmt.Errorf("owner_unread")
	}
	return common.BytesToAddress(out[12:]).Hex(), nil
}

func AgentURI(ch config.Chain, agentID uint64) (string, error) {
	data, err := EncodeTokenURI(agentID)
	if err != nil {
		return "", err
	}
	out, err := chainrpc.Call(ch.RPC, ch.Identity8004, CalldataHex(data))
	if err != nil {
		return "", err
	}
	vals, err := identityABI.Unpack("tokenURI", out)
	if err != nil || len(vals) == 0 {
		return "", fmt.Errorf("uri_unread")
	}
	s, _ := vals[0].(string)
	return s, nil
}

func LastIndex(ch config.Chain, agentID uint64, client string) (uint64, error) {
	data, err := EncodeGetLastIndex(agentID, client)
	if err != nil {
		return 0, err
	}
	out, err := chainrpc.Call(ch.RPC, ch.Reputation8004, CalldataHex(data))
	if err != nil {
		return 0, err
	}
	if len(out) == 0 {
		return 0, nil
	}
	return new(big.Int).SetBytes(out).Uint64(), nil
}

func ReadFeedback(ch config.Chain, agentID uint64, client string, index uint64) (FeedbackRead, error) {
	var got FeedbackRead
	data, err := EncodeReadFeedback(agentID, client, index)
	if err != nil {
		return got, err
	}
	out, err := chainrpc.Call(ch.RPC, ch.Reputation8004, CalldataHex(data))
	if err != nil {
		return got, err
	}
	vals, err := reputationABI.Unpack("readFeedback", out)
	if err != nil || len(vals) < 5 {
		return got, fmt.Errorf("feedback_unread")
	}
	got.Value, _ = vals[0].(*big.Int)
	if d, ok := vals[1].(uint8); ok {
		got.ValueDecimals = d
	}
	got.Tag1, _ = vals[2].(string)
	got.Tag2, _ = vals[3].(string)
	got.Revoked, _ = vals[4].(bool)
	return got, nil
}

func GetSummary(ch config.Chain, agentID uint64, clients []string, tag1, tag2 string) (Summary, error) {
	var s Summary
	data, err := EncodeGetSummary(agentID, clients, tag1, tag2)
	if err != nil {
		return s, err
	}
	out, err := chainrpc.Call(ch.RPC, ch.Reputation8004, CalldataHex(data))
	if err != nil {
		return s, err
	}
	vals, err := reputationABI.Unpack("getSummary", out)
	if err != nil || len(vals) < 3 {
		return s, fmt.Errorf("summary_unread")
	}
	switch v := vals[0].(type) {
	case uint64:
		s.Count = v
	case *big.Int:
		s.Count = v.Uint64()
	}
	s.Value, _ = vals[1].(*big.Int)
	if d, ok := vals[2].(uint8); ok {
		s.ValueDecimals = d
	}
	return s, nil
}

func IdentityRegistry(ch config.Chain) (string, error) {
	data, err := EncodeGetIdentityRegistry()
	if err != nil {
		return "", err
	}
	out, err := chainrpc.Call(ch.RPC, ch.Reputation8004, CalldataHex(data))
	if err != nil {
		return "", err
	}
	if len(out) < 32 {
		return "", fmt.Errorf("registry_unread")
	}
	return common.BytesToAddress(out[12:]).Hex(), nil
}

func Implementation(rpc, proxy string) (string, error) {
	raw, err := chainrpc.StorageAt(rpc, proxy, eip1967ImplSlot)
	if err != nil {
		return "", err
	}
	if len(raw) < 20 {
		return "", fmt.Errorf("impl_unread")
	}
	return common.BytesToAddress(raw[len(raw)-20:]).Hex(), nil
}

func IdentitySnapshot(ch config.Chain, agentID uint64) (IdentityLive, error) {
	var live IdentityLive
	owner, err := OwnerOf(ch, agentID)
	if err != nil {
		return live, err
	}
	live.Owner = owner
	live.URI, _ = AgentURI(ch, agentID)
	return live, nil
}
