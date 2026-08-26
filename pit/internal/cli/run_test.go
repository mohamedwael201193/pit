package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/session"
)

func TestRunAuthorizeNeedsLiveSession(t *testing.T) {
	s := session.Session{ID: "s", Workspace: "w", AgentAddr: "0xa", Expires: 10, PolicyVer: "v1"}
	if err := RunAuthorize(true, "AUTHORIZE", true, s, 10); err == nil {
		t.Fatal("expired")
	}
	s.Expires = 11
	if err := RunAuthorize(true, "yes", true, s, 1); err == nil {
		t.Fatal("token")
	}
	if err := RunAuthorize(true, "AUTHORIZE", true, s, 1); err != nil {
		t.Fatal(err)
	}
}
