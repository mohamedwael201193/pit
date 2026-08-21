package ledger

import "testing"

func TestConcurrentDuplicateCloid(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir, "testnet", "ws-concurrent")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	rec := Record{Workspace: "ws-concurrent", Cloid: "0xabc", Preview: "p", Status: StatusPreviewed}
	ok1, err := s.Apply(rec)
	if err != nil || !ok1 {
		t.Fatalf("%v %v", ok1, err)
	}
	ok2, err := s.Apply(rec)
	if err != nil {
		t.Fatal(err)
	}
	if ok2 {
		t.Fatal("duplicate applied")
	}
}
