package session

import "testing"

func TestMetaSession(t *testing.T) {
	m := Meta{
		ID:        NewID(),
		AgentAddr: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Workspace: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		Network:   "mainnet",
		PolicyVer: "v1",
		Expires:   99,
	}
	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
	s := m.Session()
	if s.AgentAddr != m.AgentAddr || s.Workspace != m.Workspace {
		t.Fatalf("%+v", s)
	}
	if Alive(s, 99) {
		t.Fatal("expired at boundary")
	}
}
