package session

import (
	"testing"
	"time"
)

func TestAllowlist(t *testing.T) {
	for _, a := range []string{"order", "cancel"} {
		if err := CheckAction(a); err != nil {
			t.Fatal(err)
		}
	}
	for a := range DeniedActions {
		if err := CheckAction(a); err == nil {
			t.Fatalf("%s must deny", a)
		}
	}
	if err := CheckAction("withdraw3"); err == nil {
		t.Fatal("withdraw")
	}
	if err := CheckAction("approveAgent"); err == nil {
		t.Fatal("approveAgent")
	}
	if err := CheckAction("updateLeverage"); err == nil {
		t.Fatal("leverage")
	}
}

func TestSessionWorkspace(t *testing.T) {
	now := time.Now().UnixMilli()
	s := Session{ID: "s", Workspace: "ws-a", AgentAddr: "0xabc", Expires: now + 1000, PolicyVer: "v1"}
	if err := CheckSession(s, now, "v1", "ws-a"); err != nil {
		t.Fatal(err)
	}
	if err := CheckSession(s, now, "v1", "ws-b"); err == nil {
		t.Fatal("cross workspace")
	}
	s.Kill = true
	if err := CheckSession(s, now, "v1", "ws-a"); err == nil {
		t.Fatal("kill")
	}
}
