package companion

import "testing"

func TestSnapshotNeverAttachesOid(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	snap := h.snapshotResearch()
	if _, ok := snap["oid"]; ok {
		t.Fatal("research snapshot must not carry a venue oid")
	}
	if snap["sign"] != false || snap["trade"] != false {
		t.Fatal(snap)
	}
}
