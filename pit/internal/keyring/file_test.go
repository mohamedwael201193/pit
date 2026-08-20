package keyring

import "testing"

func TestNamespaceIsolation(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Put("ws-a", "session", []byte("aaa")); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("ws-b", "session"); err == nil {
		t.Fatal("B must not read A")
	}
	got, err := s.Get("ws-a", "session")
	if err != nil || string(got) != "aaa" {
		t.Fatal(err)
	}
	k, err := NewMemoryKey()
	if err != nil || len(k) != 66 {
		t.Fatalf("%s %v", k, err)
	}
}
