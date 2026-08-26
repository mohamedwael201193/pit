package session

import "testing"

func TestAlive(t *testing.T) {
	s := Session{ID: "s", Workspace: "w", AgentAddr: "0xa", Expires: 20, PolicyVer: "v1"}
	if Alive(s, 20) {
		t.Fatal("expired")
	}
	if !Alive(s, 19) {
		t.Fatal("live")
	}
	s.Revoked = true
	if Alive(s, 19) {
		t.Fatal("revoked")
	}
}
