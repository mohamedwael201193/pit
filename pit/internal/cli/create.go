package cli

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mohamedwael201193/pit/internal/keyring"
	"github.com/mohamedwael201193/pit/internal/session"
)

func CreateLocalSession(dir, workspace, network, policyVer string) (SessionFile, error) {
	if err := session.CapTTLHours(1); err != nil {
		return SessionFile{}, err
	}
	key, exp, err := session.GenerateAgent(1)
	if err != nil {
		return SessionFile{}, err
	}
	ring, err := keyring.OpenProduct(filepath.Join(dir, "keyring"))
	if err != nil {
		return SessionFile{}, err
	}
	id := session.NewID()
	if err := key.Store(ring, workspace, id); err != nil {
		return SessionFile{}, err
	}
	sf := SessionFile{
		ID:        id,
		AgentAddr: key.Address,
		Workspace: workspace,
		Network:   network,
		PolicyVer: policyVer,
		Expires:   exp,
	}
	if err := SaveSession(dir, sf); err != nil {
		return SessionFile{}, err
	}
	if err := key.ExportJSON(); err == nil {
		return SessionFile{}, fmt.Errorf("session_export_denied")
	}
	return sf, nil
}

func LiveFromDisk(dir string, kill bool, nowMs int64) (session.Session, error) {
	sf, err := LoadSession(dir)
	if err != nil {
		return session.Session{}, err
	}
	s := sf.Meta().Session()
	s.Kill = kill
	if !session.Alive(s, nowMs) {
		return session.Session{}, fmt.Errorf("session_expired")
	}
	return s, nil
}

func KeyringDir(dir string) string {
	return filepath.Join(dir, "keyring")
}

func NowMs() int64 {
	return time.Now().UTC().UnixMilli()
}
