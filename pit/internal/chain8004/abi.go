package chain8004

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	GiveFeedbackSelector = "0x3c036a7e"
	SetAgentURISelector  = "0x0af28bd3"
	RegisterStringSel    = "0xf2c298be"
	MainnetAgentID       = 3489333
	RecordedHLAgent      = "0xfc64e36babe7dfe9eb779ee3a9f2362d16881d52"
	RecordedOwner        = "0xbdfcee82bd42fefa58ee850b3709636a8b6b0034"
	AgentCardURL         = "https://pit0g.vercel.app/.well-known/agent-card.json"
	ProofURL             = "https://pit0g.vercel.app/proof"
	WebEndpoint          = "https://pit0g.vercel.app/"
)

var identityABI = mustABI(`[
  {"type":"function","name":"ownerOf","inputs":[{"name":"tokenId","type":"uint256"}],"outputs":[{"type":"address"}],"stateMutability":"view"},
  {"type":"function","name":"tokenURI","inputs":[{"name":"tokenId","type":"uint256"}],"outputs":[{"type":"string"}],"stateMutability":"view"},
  {"type":"function","name":"setAgentURI","inputs":[{"name":"agentId","type":"uint256"},{"name":"newURI","type":"string"}]},
  {"type":"function","name":"register","inputs":[{"name":"agentURI","type":"string"}],"outputs":[{"type":"uint256"}]}
]`)

var reputationABI = mustABI(`[
  {"type":"function","name":"giveFeedback","inputs":[
    {"name":"agentId","type":"uint256"},
    {"name":"value","type":"int128"},
    {"name":"valueDecimals","type":"uint8"},
    {"name":"tag1","type":"string"},
    {"name":"tag2","type":"string"},
    {"name":"endpoint","type":"string"},
    {"name":"feedbackURI","type":"string"},
    {"name":"feedbackHash","type":"bytes32"}
  ]},
  {"type":"function","name":"getLastIndex","inputs":[{"name":"agentId","type":"uint256"},{"name":"clientAddress","type":"address"}],"outputs":[{"type":"uint64"}],"stateMutability":"view"},
  {"type":"function","name":"readFeedback","inputs":[
    {"name":"agentId","type":"uint256"},
    {"name":"clientAddress","type":"address"},
    {"name":"feedbackIndex","type":"uint64"}
  ],"outputs":[
    {"name":"value","type":"int128"},
    {"name":"valueDecimals","type":"uint8"},
    {"name":"tag1","type":"string"},
    {"name":"tag2","type":"string"},
    {"name":"isRevoked","type":"bool"}
  ],"stateMutability":"view"},
  {"type":"function","name":"getSummary","inputs":[
    {"name":"agentId","type":"uint256"},
    {"name":"clientAddresses","type":"address[]"},
    {"name":"tag1","type":"string"},
    {"name":"tag2","type":"string"}
  ],"outputs":[
    {"name":"count","type":"uint64"},
    {"name":"summaryValue","type":"int128"},
    {"name":"summaryValueDecimals","type":"uint8"}
  ],"stateMutability":"view"},
  {"type":"function","name":"getIdentityRegistry","outputs":[{"type":"address"}],"stateMutability":"view"}
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

type FeedbackArgs struct {
	AgentID        uint64
	Value          int64
	ValueDecimals  uint8
	Tag1           string
	Tag2           string
	Endpoint       string
	FeedbackURI    string
	FeedbackHash   common.Hash
}

func EncodeGiveFeedback(a FeedbackArgs) ([]byte, error) {
	if a.AgentID == 0 {
		return nil, fmt.Errorf("agent_required")
	}
	if err := RefusePrivateTag(a.Tag1); err != nil {
		return nil, err
	}
	if err := RefusePrivateTag(a.Tag2); err != nil {
		return nil, err
	}
	return reputationABI.Pack(
		"giveFeedback",
		new(big.Int).SetUint64(a.AgentID),
		big.NewInt(a.Value),
		a.ValueDecimals,
		a.Tag1,
		a.Tag2,
		a.Endpoint,
		a.FeedbackURI,
		a.FeedbackHash,
	)
}

func EncodeSetAgentURI(agentID uint64, uri string) ([]byte, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, fmt.Errorf("uri_required")
	}
	return identityABI.Pack("setAgentURI", new(big.Int).SetUint64(agentID), uri)
}

func EncodeOwnerOf(agentID uint64) ([]byte, error) {
	return identityABI.Pack("ownerOf", new(big.Int).SetUint64(agentID))
}

func EncodeTokenURI(agentID uint64) ([]byte, error) {
	return identityABI.Pack("tokenURI", new(big.Int).SetUint64(agentID))
}

func EncodeGetLastIndex(agentID uint64, client string) ([]byte, error) {
	return reputationABI.Pack("getLastIndex", new(big.Int).SetUint64(agentID), common.HexToAddress(client))
}

func EncodeReadFeedback(agentID uint64, client string, index uint64) ([]byte, error) {
	return reputationABI.Pack("readFeedback", new(big.Int).SetUint64(agentID), common.HexToAddress(client), index)
}

func EncodeGetSummary(agentID uint64, clients []string, tag1, tag2 string) ([]byte, error) {
	addrs := make([]common.Address, 0, len(clients))
	for _, c := range clients {
		addrs = append(addrs, common.HexToAddress(c))
	}
	return reputationABI.Pack("getSummary", new(big.Int).SetUint64(agentID), addrs, tag1, tag2)
}

func EncodeGetIdentityRegistry() ([]byte, error) {
	return reputationABI.Pack("getIdentityRegistry")
}

func CalldataHex(data []byte) string {
	return "0x" + common.Bytes2Hex(data)
}
