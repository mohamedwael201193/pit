package policy

import "testing"

func TestRequireSession(t *testing.T) {
	if err := RequireSession(Context{}); err == nil {
		t.Fatal("dead")
	}
	if err := RequireSession(Context{SessionAlive: true}); err != nil {
		t.Fatal(err)
	}
}
