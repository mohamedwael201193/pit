package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

// DiskState is the local CLI bind. It never stores a session secret or mnemonic.
type DiskState struct {
	WorkspaceID string `json:"workspaceId"`
	Network     string `json:"network"`
	Wallet      string `json:"wallet"`
	Kill        bool   `json:"kill"`
}

func DefaultDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "pit"), nil
}

func Path(dir string) string {
	return filepath.Join(dir, "state.json")
}

func Save(dir string, st DiskState) error {
	if _, err := identity.ParseWorkspaceID(st.WorkspaceID); err != nil {
		return err
	}
	if _, err := config.ParseNetwork(st.Network); err != nil {
		return err
	}
	if _, err := identity.NormalizeAddress(st.Wallet); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(string(b)), "private") || strings.Contains(string(b), "0x"+strings.Repeat("a", 64)) {
		return fmt.Errorf("refusing_secret_material")
	}
	return os.WriteFile(Path(dir), b, 0o600)
}

func Load(dir string) (DiskState, error) {
	b, err := os.ReadFile(Path(dir))
	if err != nil {
		return DiskState{}, err
	}
	var st DiskState
	if err := json.Unmarshal(b, &st); err != nil {
		return DiskState{}, err
	}
	if st.WorkspaceID == "" {
		return DiskState{}, fmt.Errorf("unbound")
	}
	return st, nil
}

func SetKill(dir string, on bool) error {
	st, err := Load(dir)
	if err != nil {
		return err
	}
	st.Kill = on
	return Save(dir, st)
}
