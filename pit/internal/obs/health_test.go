package obs

import "testing"

func TestRefuseHealthSecrets(t *testing.T) {
	if err := RefuseHealthSecrets(map[string]any{"ok": true, "sign": false}); err != nil {
		t.Fatal(err)
	}
	if err := RefuseHealthSecrets(map[string]any{"session": "x"}); err == nil {
		t.Fatal("leak")
	}
}
