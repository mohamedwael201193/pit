package chain8004

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGiveFeedbackSelector(t *testing.T) {
	got := Selector("giveFeedback(uint256,int128,uint8,string,string,string,string,bytes32)")
	if got != GiveFeedbackSelector {
		t.Fatal(got)
	}
	if Selector("setAgentURI(uint256,string)") != SetAgentURISelector {
		t.Fatal(Selector("setAgentURI(uint256,string)"))
	}
	if Selector("register(string)") != RegisterStringSel {
		t.Fatal(Selector("register(string)"))
	}
}

func TestEncodeGiveFeedbackPrefix(t *testing.T) {
	data, err := EncodeGiveFeedback(FeedbackArgs{
		AgentID: MainnetAgentID, Value: 1, Tag1: "hype_fill", Tag2: "successful_job",
		Endpoint: WebEndpoint, FeedbackURI: ProofURL, FeedbackHash: common.HexToHash("0x01"),
	})
	if err != nil {
		t.Fatal(err)
	}
	hex := CalldataHex(data)
	if !strings.HasPrefix(hex, GiveFeedbackSelector) {
		t.Fatal(hex[:10])
	}
}

func TestEncodeGiveFeedbackRefusesPrivateTag(t *testing.T) {
	_, err := EncodeGiveFeedback(FeedbackArgs{
		AgentID: MainnetAgentID, Value: 1, Tag1: "book", Tag2: "ok",
		Endpoint: WebEndpoint, FeedbackURI: ProofURL,
	})
	if err == nil {
		t.Fatal("private tag")
	}
}

func TestEncodeSetAgentURIPrefix(t *testing.T) {
	data, err := EncodeSetAgentURI(MainnetAgentID, AgentCardURL)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(CalldataHex(data), SetAgentURISelector) {
		t.Fatal(CalldataHex(data)[:10])
	}
}
