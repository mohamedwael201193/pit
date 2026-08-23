package session

import "testing"

func TestExpiredAndRevoked(t *testing.T) {
	s := Session{ID: "s", Workspace: "w", AgentAddr: "0xa", Expires: 10, PolicyVer: "v1"}
	if err := CheckSession(s, 10, "v1", "w"); err == nil {
		t.Fatal("expired")
	}
	s.Expires = 11
	s.Revoked = true
	if err := CheckSession(s, 1, "v1", "w"); err == nil {
		t.Fatal("revoked")
	}
}
