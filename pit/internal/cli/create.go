package cli

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mohamedwael201193/pit/internal/keyring"
	"github.com/mohamedwael201193/pit/internal/session"
)

func SessionPublic(sf SessionFile) map[string]any {
	return map[string]any{
		"ok":        true,
		"agent":     sf.AgentAddr,
		"expires":   sf.Expires,
		"workspace": sf.Workspace,
		"network":   sf.Network,
		"order":     true,
		"cancel":    true,
		"withdraw":  false,
		"transfer":  false,
		"leverage":  false,
		"sign":      false,
		"trade":     false,
	}
}

func EnsureLocalSession(dir string) (SessionFile, bool, error) {
	st, err := Load(dir)
	if err != nil {
		return SessionFile{}, false, fmt.Errorf("unbound")
	}
	now := time.Now().UnixMilli()
	sf, loadErr := LoadSession(dir)
	if loadErr == nil {
		linked, until, qerr := LookupAgent(st.Network, st.Wallet, st.WorkspaceID, sf.AgentAddr, now)
		if qerr == nil && linked {
			return extendApprovedSession(dir, sf, now, until)
		}
		if _, err := LiveFromDisk(dir, st.Kill, now); err == nil {
			return sf, false, nil
		}
		if qerr != nil {
			return SessionFile{}, false, qerr
		}
	}
	sf, err = CreateLocalSession(dir, st.WorkspaceID, st.Network, "v1")
	return sf, true, err
}

func extendApprovedSession(dir string, sf SessionFile, nowMs, until int64) (SessionFile, bool, error) {
	want := until
	if want <= nowMs {
		want = nowMs + int64(session.DefaultTTLHours)*3600*1000
	}
	if sf.Expires != want {
		sf.Expires = want
		if err := SaveSession(dir, sf); err != nil {
			return SessionFile{}, false, err
		}
	}
	return sf, false, nil
}

func CreateLocalSession(dir, workspace, network, policyVer string) (SessionFile, error) {
	if err := session.CapTTLHours(session.DefaultTTLHours); err != nil {
		return SessionFile{}, err
	}
	key, exp, err := session.GenerateAgent(session.DefaultTTLHours)
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
