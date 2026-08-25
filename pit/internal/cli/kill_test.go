package cli

import "testing"

func TestRefuseKilledSession(t *testing.T) {
	if err := RefuseKilledSession(true); err == nil {
		t.Fatal("kill")
	}
	if err := RefuseKilledSession(false); err != nil {
		t.Fatal(err)
	}
}
