package session

import "testing"

func TestRequireAgentName(t *testing.T) {
	ws := "abcdefgh-1234-5678-1234-567812345678"
	want := AgentKey{}.Name(ws)
	if err := RequireAgentName(want, ws); err != nil {
		t.Fatal(err)
	}
	if err := RequireAgentName("bot", ws); err == nil {
		t.Fatal("name")
	}
}
