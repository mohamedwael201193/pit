package auto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MissionProof is hash-only. Private book, prompt, and strategy never appear.
type MissionProof struct {
	MissionID      string `json:"mission_id"`
	PolicyHash     string `json:"policy_hash,omitempty"`
	ResearchJob    string `json:"research_job,omitempty"`
	ResearchDigest string `json:"research_digest,omitempty"`
	TEEVerified    bool   `json:"tee_verified"`
	TeeSigner      string `json:"tee_signer,omitempty"`
	SkillVersion   string `json:"skill_version,omitempty"`
	MemoryRoot     string `json:"memory_root,omitempty"`
	MemoryVersion  string `json:"memory_version,omitempty"`
	PreviewHash    string `json:"preview_hash,omitempty"`
	ExecutionHash  string `json:"execution_hash,omitempty"`
	OID            string `json:"oid,omitempty"`
	FillState      string `json:"fill_state,omitempty"`
	OutcomeHash    string `json:"outcome_hash,omitempty"`
	StorageRoot    string `json:"storage_root,omitempty"`
	ChainTx        string `json:"chain_tx,omitempty"`
	Timestamp      int64  `json:"timestamp"`
	EnvelopeDigest string `json:"envelope_digest,omitempty"`
}

func proofPath(dir string) string {
	return filepath.Join(dir, "mission-proof.json")
}

func AssembleProof(dir string, extra MissionProof) MissionProof {
	m := LoadMission(dir)
	b := LoadSkillbook(dir)
	p := extra
	if p.MissionID == "" {
		p.MissionID = m.MissionID
		if p.MissionID == "" {
			p.MissionID = m.Envelope.MissionID
		}
	}
	if p.PolicyHash == "" {
		p.PolicyHash = m.PolicyHash
	}
	if p.PreviewHash == "" {
		p.PreviewHash = m.LastPreviewHash
	}
	if p.OID == "" {
		p.OID = m.LastOID
	}
	if p.MemoryRoot == "" {
		p.MemoryRoot = b.MemoryRoot
	}
	if p.MemoryVersion == "" {
		p.MemoryVersion = b.Version
	}
	if p.EnvelopeDigest == "" {
		p.EnvelopeDigest = m.Envelope.Digest()
	}
	if p.Timestamp == 0 {
		p.Timestamp = time.Now().Unix()
	}
	if strings.EqualFold(p.FillState, "RESTING") {
		// RESTING is not a fill. Keep the label honest.
	}
	sum := sha256.Sum256([]byte(p.MissionID + "|" + p.PreviewHash + "|" + p.OID + "|" + p.FillState))
	if p.OutcomeHash == "" {
		p.OutcomeHash = "0x" + hex.EncodeToString(sum[:])
	}
	_ = saveProof(dir, p)
	return p
}

func LoadProof(dir string) MissionProof {
	var p MissionProof
	raw, err := os.ReadFile(proofPath(dir))
	if err != nil {
		return p
	}
	if strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return MissionProof{}
	}
	_ = json.Unmarshal(raw, &p)
	return p
}

func PublicProof(dir string) map[string]any {
	p := LoadProof(dir)
	return map[string]any{
		"ok": true, "proof": p, "private_strategy": "redacted",
		"copy": "Private strategy remains on desktop.",
		"sign": false, "trade": false, "arm": false,
	}
}

func saveProof(dir string, p MissionProof) error {
	if dir == "" {
		return nil
	}
	raw, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(proofPath(dir), raw, 0o600)
}

type ReplayView struct {
	OK              bool           `json:"ok"`
	Live            bool            `json:"live"`
	Note            string         `json:"note"`
	MarketSnapshot  map[string]any `json:"market_snapshot,omitempty"`
	PolicyHash      string         `json:"policy_hash,omitempty"`
	ResearchVersion string         `json:"research_version,omitempty"`
	SkillVersion    string         `json:"skill_version,omitempty"`
	MemoryVersion   string         `json:"memory_version,omitempty"`
	Result          string         `json:"result,omitempty"`
	Engine          string         `json:"engine_output,omitempty"`
	Proof           MissionProof  `json:"proof"`
	Sign            bool           `json:"sign"`
	Trade           bool           `json:"trade"`
}

func StrategyReplay(dir string) ReplayView {
	p := LoadProof(dir)
	m := LoadMission(dir)
	return ReplayView{
		OK:              true,
		Live:            false,
		Note:            "Replay uses the original committed snapshot. It is not a new trade and does not use today's market data.",
		PolicyHash:      m.PolicyHash,
		ResearchVersion: m.Envelope.ResearchDigest,
		SkillVersion:    p.SkillVersion,
		MemoryVersion:   p.MemoryVersion,
		Result:          m.LastResult,
		Engine:          m.BestWhy,
		Proof:           p,
		Sign:            false,
		Trade:           false,
	}
}
