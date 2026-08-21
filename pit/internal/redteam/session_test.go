package redteam

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/session"
)

func TestSessionExpiryAndRevoke(t *testing.T) {
	s := session.Session{
		ID: "s1", Workspace: "ws-a", AgentAddr: "0xagent", Expires: 100, PolicyVer: "v1",
	}
	if err := session.CheckSession(s, 50, "v1", "ws-a"); err != nil {
		t.Fatal(err)
	}
	if err := session.CheckSession(s, 100, "v1", "ws-a"); err == nil {
		t.Fatal("expired")
	}
	s.Revoked = true
	s.Expires = 999
	if err := session.CheckSession(s, 50, "v1", "ws-a"); err == nil {
		t.Fatal("revoked")
	}
}

func TestWrongWorkspaceAndPolicy(t *testing.T) {
	s := session.Session{ID: "s1", Workspace: "ws-a", AgentAddr: "0xagent", Expires: 999, PolicyVer: "v1"}
	if err := session.CheckSession(s, 1, "v1", "ws-b"); err == nil {
		t.Fatal("ws")
	}
	if err := session.CheckSession(s, 1, "v2", "ws-a"); err == nil {
		t.Fatal("policy")
	}
}

func TestWithdrawNeverAllowlisted(t *testing.T) {
	for _, a := range []string{"withdraw3", "usdSend", "approveAgent", "updateLeverage"} {
		if err := session.CheckAction(a); err == nil {
			t.Fatal(a)
		}
	}
}
