package session

import "testing"

func TestAgentName(t *testing.T) {
	n, err := AgentName("abcd1234-ffff-4000-8000-000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if n != "PIT-abcd1234" {
		t.Fatalf("%s", n)
	}
	if _, err := AgentName("short"); err == nil {
		t.Fatal("short")
	}
	if len(n) >= 17 {
		t.Fatal("hyperliquid name limit")
	}
}
