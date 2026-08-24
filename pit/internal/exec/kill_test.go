package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/session"
)

func TestRefuseKill(t *testing.T) {
	if err := RefuseKill(session.Session{}); err != nil {
		t.Fatal(err)
	}
	if err := RefuseKill(session.Session{Kill: true}); err == nil {
		t.Fatal("kill")
	}
}
