package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedwael201193/pit/internal/session"
)

type SessionFile struct {
	ID        string `json:"id"`
	AgentAddr string `json:"agentAddr"`
	Workspace string `json:"workspace"`
	Network   string `json:"network"`
	PolicyVer string `json:"policyVer"`
	Expires   int64  `json:"expires"`
}

func SessionPath(dir string) string {
	return filepath.Join(dir, "session.json")
}

func RefuseSessionSecret(b []byte) error {
	low := strings.ToLower(string(b))
	for _, w := range []string{"private_key", "mnemonic", "session_key", "hl_secret", "app-sk-"} {
		if strings.Contains(low, w) {
			return fmt.Errorf("secret_print")
		}
	}
	return nil
}

func SaveSession(dir string, s SessionFile) error {
	if s.ID == "" || s.AgentAddr == "" || s.Workspace == "" {
		return fmt.Errorf("incomplete_session")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := RefuseSessionSecret(b); err != nil {
		return err
	}
	return os.WriteFile(SessionPath(dir), b, 0o600)
}

func LoadSession(dir string) (SessionFile, error) {
	b, err := os.ReadFile(SessionPath(dir))
	if err != nil {
		return SessionFile{}, fmt.Errorf("session_expired")
	}
	if err := RefuseSessionSecret(b); err != nil {
		return SessionFile{}, err
	}
	var s SessionFile
	if err := json.Unmarshal(b, &s); err != nil {
		return SessionFile{}, err
	}
	if s.ID == "" || s.AgentAddr == "" {
		return SessionFile{}, fmt.Errorf("session_expired")
	}
	return s, nil
}

func (s SessionFile) Meta() session.Meta {
	return session.Meta{
		ID:        s.ID,
		AgentAddr: s.AgentAddr,
		Workspace: s.Workspace,
		Network:   s.Network,
		PolicyVer: s.PolicyVer,
		Expires:   s.Expires,
	}
}
