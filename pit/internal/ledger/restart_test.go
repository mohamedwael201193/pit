package ledger

import "testing"

func TestRestartKeepsPreview(t *testing.T) {
	dir := t.TempDir()
	ws := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	a, err := Open(dir, "testnet", ws)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := a.Apply(Record{Workspace: ws, Cloid: "0x3", Preview: "h", Status: StatusPreviewed})
	if err != nil || !ok {
		t.Fatal(err)
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := Open(dir, "testnet", ws)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	rec, err := b.Get(ws, "0x3")
	if err != nil || rec.Status != StatusPreviewed {
		t.Fatalf("%+v %v", rec, err)
	}
	ok, err = b.Apply(Record{Workspace: ws, Cloid: "0x3", Preview: "h", Status: StatusPreviewed})
	if err != nil || ok {
		t.Fatal("duplicate click")
	}
}
