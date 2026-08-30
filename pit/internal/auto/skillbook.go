package auto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedwael201193/pit/internal/experience"
	"github.com/mohamedwael201193/pit/internal/keyring"
)

const (
	KindObservation     = "observation"
	KindForecast        = "forecast"
	KindExecution       = "execution"
	KindOutcome         = "outcome"
	KindError           = "error"
	KindCalibration     = "calibration"
	KindRolePerformance = "role_performance"
	KindSkillPerformance = "skill_performance"
	KindRiskEvent       = "risk_event"
	KindPolicyDecision   = "policy_decision"
	KindMission         = "mission"
	KindMissionOutcome  = "mission_outcome"
	KindStrategyLesson  = "strategy_lesson"
)

type SkillRecord struct {
	SkillID            string  `json:"skill_id"`
	Version            string  `json:"version"`
	Description        string  `json:"description,omitempty"`
	MarketRegime       string  `json:"market_regime,omitempty"`
	Inputs             string  `json:"inputs,omitempty"`
	DecisionFamily     string  `json:"decision_family,omitempty"`
	N                  int     `json:"n"`
	ResolvedN          int     `json:"resolved_n"`
	Brier              float64 `json:"brier,omitempty"`
	ECE                float64 `json:"ece,omitempty"`
	Wins               int     `json:"wins,omitempty"`
	Losses             int     `json:"losses,omitempty"`
	ChallengerCatch    int     `json:"challenger_catch_rate,omitempty"`
	LastUpdatedUnix    int64   `json:"last_updated,omitempty"`
	ConfidenceState    string  `json:"confidence_state"`
}

type MemoryEntry struct {
	Kind    string `json:"kind"`
	Unix    int64  `json:"unix"`
	Coin    string  `json:"coin,omitempty"`
	Digest  string  `json:"digest,omitempty"`
	Why     string  `json:"why,omitempty"`
	Mission string `json:"mission_id,omitempty"`
}

type Skillbook struct {
	Version     string         `json:"version"`
	MemoryRoot  string         `json:"memory_root,omitempty"`
	MemoryVer   string         `json:"memory_version,omitempty"`
	Skills      []SkillRecord  `json:"skills"`
	Entries     []MemoryEntry  `json:"entries"`
}

func skillbookPath(dir string) string {
	return filepath.Join(dir, "skillbook.enc")
}

func EmptySkillbook() Skillbook {
	return Skillbook{Version: "1", Skills: []SkillRecord{}, Entries: []MemoryEntry{}}
}

func LoadSkillbook(dir string) Skillbook {
	empty := EmptySkillbook()
	if dir == "" {
		return empty
	}
	raw, err := os.ReadFile(skillbookPath(dir))
	if err != nil {
		return empty
	}
	plain, err := openSkillbook(dir, string(raw))
	if err != nil || strings.Contains(strings.ToLower(string(plain)), "app-sk-") {
		return empty
	}
	var b Skillbook
	if json.Unmarshal(plain, &b) != nil {
		return empty
	}
	if b.Skills == nil {
		b.Skills = []SkillRecord{}
	}
	if b.Entries == nil {
		b.Entries = []MemoryEntry{}
	}
	return b
}

func SaveSkillbook(dir string, b Skillbook) error {
	if dir == "" {
		return fmt.Errorf("unbound")
	}
	b.MemoryRoot = skillbookRoot(b)
	raw, err := json.Marshal(b)
	if err != nil {
		return err
	}
	sealed, err := sealSkillbook(dir, raw)
	if err != nil {
		return err
	}
	return os.WriteFile(skillbookPath(dir), []byte(sealed+"\n"), 0o600)
}

func RecordMemory(dir, kind, coin, why, missionID string, now int64) Skillbook {
	b := LoadSkillbook(dir)
	sum := sha256.Sum256([]byte(kind + "|" + coin + "|" + why + "|" + missionID))
	b.Entries = append(b.Entries, MemoryEntry{
		Kind: kind, Unix: now, Coin: coin, Why: why, Mission: missionID,
		Digest: "0x" + hex.EncodeToString(sum[:8]),
	})
	if len(b.Entries) > 400 {
		b.Entries = b.Entries[len(b.Entries)-400:]
	}
	_ = SaveSkillbook(dir, b)
	return LoadSkillbook(dir)
}

func SkillHonesty(s SkillRecord) string {
	if s.N < experience.MinSamples || s.ResolvedN < experience.MinSamples {
		return "NOT ENOUGH DATA"
	}
	return s.ConfidenceState
}

func PublicSkillbook(dir string) map[string]any {
	b := LoadSkillbook(dir)
	rows := []map[string]any{}
	for _, s := range b.Skills {
		rows = append(rows, map[string]any{
			"skill_id": s.SkillID, "version": s.Version, "n": s.N, "resolved_n": s.ResolvedN,
			"confidence_state": SkillHonesty(s), "sign": false, "trade": false,
		})
	}
	return map[string]any{
		"ok": true, "memory_root": b.MemoryRoot, "memory_version": b.Version,
		"n": len(b.Entries), "skills": rows, "copy": honestyCopy(b),
		"sign": false, "trade": false, "authorize": false,
	}
}

func honestyCopy(b Skillbook) string {
	if len(b.Entries) < experience.MinSamples {
		return fmt.Sprintf("NOT ENOUGH DATA (%d/%d verified memory rows). PIT will not invent skill performance.", len(b.Entries), experience.MinSamples)
	}
	return fmt.Sprintf("%d verified memory rows. Selective retrieval only. Private strategy is not stored.", len(b.Entries))
}

func skillbookRoot(b Skillbook) string {
	raw, _ := json.Marshal(struct {
		V string         `json:"v"`
		S []SkillRecord  `json:"s"`
		N int            `json:"n"`
	}{V: b.Version, S: b.Skills, N: len(b.Entries)})
	sum := sha256.Sum256(raw)
	return "0x" + hex.EncodeToString(sum[:])
}

func sealSkillbook(dir string, plain []byte) (string, error) {
	key, err := skillbookKey(dir)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, plain, nil)
	return memPrefixAuto + hex.EncodeToString(out), nil
}

func openSkillbook(dir, line string) ([]byte, error) {
	raw := strings.TrimSpace(line)
	if !strings.HasPrefix(raw, memPrefixAuto) {
		return nil, fmt.Errorf("skillbook_sealed_required")
	}
	key, err := skillbookKey(dir)
	if err != nil {
		return nil, err
	}
	bin, err := hex.DecodeString(strings.TrimPrefix(raw, memPrefixAuto))
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(bin) < gcm.NonceSize() {
		return nil, fmt.Errorf("skillbook_corrupt")
	}
	nonce, rest := bin[:gcm.NonceSize()], bin[gcm.NonceSize():]
	return gcm.Open(nil, nonce, rest, nil)
}

const memPrefixAuto = "enc:v1:"

func skillbookKey(dir string) ([]byte, error) {
	store, err := keyring.OpenProduct(dir)
	if err != nil {
		sum := sha256.Sum256([]byte("pit-skillbook-v1|" + filepath.Clean(dir)))
		return sum[:], nil
	}
	got, err := store.Get("memory", "v1")
	if err == nil && len(got) == 32 {
		return got, nil
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := store.Put("memory", "v1", key); err != nil {
		return nil, err
	}
	return key, nil
}
