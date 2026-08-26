package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/session"
)

func TestRequireLive(t *testing.T) {
	s := session.Session{ID: "s", Workspace: "w", AgentAddr: "0xa", Expires: 5, PolicyVer: "v1"}
	if err := RequireLive(s, 5); err == nil {
		t.Fatal("expired")
	}
	s.Expires = 6
	if err := RequireLive(s, 5); err != nil {
		t.Fatal(err)
	}
}
