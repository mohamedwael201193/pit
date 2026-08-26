package ledger

import "testing"

func TestRefuseForeignWorkspace(t *testing.T) {
	dir := t.TempDir()
	a, err := Open(dir, "mainnet", "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	ok, err := a.Apply(Record{Workspace: "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", Cloid: "0x9", Preview: "h", Status: "signed"})
	if ok || err == nil || err.Error() != "wrong_workspace" {
		t.Fatalf("%v %v", ok, err)
	}
}
