package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/session"
)

func TestGateExpired(t *testing.T) {
	s := session.Session{ID: "s", Workspace: "w", AgentAddr: "0xa", Expires: 1, PolicyVer: "v1"}
	if err := Gate(s, 2, "v1", "w"); err == nil {
		t.Fatal("expired")
	}
	s.Expires = 9e18
	s.Revoked = true
	if err := Gate(s, 2, "v1", "w"); err == nil {
		t.Fatal("revoked")
	}
	s.Revoked = false
	if err := Gate(s, 2, "v1", "w"); err != nil {
		t.Fatal(err)
	}
}
